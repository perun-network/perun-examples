// Copyright 2022 PolyCrypt GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"context"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	dotchannel "github.com/perun-network/perun-polkadot-backend/channel"
	"github.com/perun-network/perun-polkadot-backend/channel/pallet"
	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	dotwallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	"math/big"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wallet"

	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"

	"github.com/pkg/errors"
)

// PaymentClient is a payment channel client.
type PaymentClient struct {
	perunClient *client.Client       // The core Perun client.
	account     wallet.Address       // The account we use for on-chain and off-chain transactions.
	currency    channel.Asset        // The currency we expect to get paid in.
	channels    chan *PaymentChannel // Accepted payment channels.
}

// SetupPaymentClient creates a new payment client.
func SetupPaymentClient(
	bus wire.Bus, // bus is used of off-chain communication.
	w *dotwallet.Wallet, // w is the wallet used for signing transactions.
	acc wallet.Account, // acc is the address of the account to be used for signing transactions.
	nodeURL string, // nodeURL is the URL of the blockchain node.
	networkId dot.NetworkID, // networkId is the identifier of the blockchain.
	queryDepth types.BlockNumber, // queryDepth is the number of blocks being evaluated when looking for events.
) (*PaymentClient, error) {
	// Connect to backend.
	api, err := dot.NewAPI(nodeURL, networkId)
	if err != nil {
		panic(err)
	}

	// Create Perun pallet and generate funder + adjudicator from it.
	perun := pallet.NewPallet(pallet.NewPerunPallet(api), api.Metadata())
	funder := pallet.NewFunder(perun, acc, 3)
	adj := pallet.NewAdjudicator(acc, perun, api, queryDepth)

	// Setup dispute watcher.
	watcher, err := local.NewWatcher(adj)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}

	// Setup Perun client.
	waddr := dotwallet.AsAddr(acc.Address())
	perunClient, err := client.New(waddr, bus, funder, adj, w, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	// Create client and start request handler.
	c := &PaymentClient{
		perunClient: perunClient,
		account:     waddr,
		currency:    &dotchannel.Asset,
		channels:    make(chan *PaymentChannel, 1),
	}

	go perunClient.Handle(c, c)
	return c, nil
}

// OpenChannel opens a new channel with the specified peer and funding.
func (c *PaymentClient) OpenChannel(peer wire.Address, amount float64) *PaymentChannel {
	// We define the channel participants. The proposer has always index 0. Here
	// we use the on-chain addresses as off-chain addresses, but we could also
	// use different ones.
	participants := []wire.Address{c.account, peer}

	// We create an initial allocation which defines the starting balances.
	initAlloc := channel.NewAllocation(2, dotchannel.Asset)
	initAlloc.SetAssetBalances(dotchannel.Asset, []channel.Bal{
		DotToPlanck(big.NewFloat(amount)), // Our initial balance.
		DotToPlanck(big.NewFloat(amount)), // Peer's initial balance.
	})

	// Prepare the channel proposal by defining the channel parameters.
	challengeDuration := uint64(10) // On-chain challenge duration in seconds.
	proposal, err := client.NewLedgerChannelProposal(
		challengeDuration,
		c.account,
		initAlloc,
		participants,
	)
	if err != nil {
		panic(err)
	}

	// Send the proposal.
	ch, err := c.perunClient.ProposeChannel(context.TODO(), proposal)
	if err != nil {
		panic(err)
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)

	return newPaymentChannel(ch, dotchannel.Asset)
}

// startWatching starts the dispute watcher for the specified channel.
func (c *PaymentClient) startWatching(ch *client.Channel) {
	go func() {
		err := ch.Watch(c)
		if err != nil {
			fmt.Printf("Watcher returned with error: %v", err)
		}
	}()
}

// AcceptedChannel returns the next accepted channel.
func (c *PaymentClient) AcceptedChannel() *PaymentChannel {
	return <-c.channels
}

// Shutdown gracefully shuts down the client.
func (c *PaymentClient) Shutdown() {
	c.perunClient.Close()
}
