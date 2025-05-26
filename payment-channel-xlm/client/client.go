// Copyright 2025 PolyCrypt GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"log"
	"math/big"

	"context"
	"errors"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"

	stellarwallet "perun.network/perun-stellar-backend/wallet"
	"perun.network/perun-stellar-backend/wallet/types"
)

type PaymentClient struct {
	perunClient *client.Client
	account     map[wallet.BackendID]wallet.Address // The account we use for on-chain and off-chain transactions.
	waddress    map[wallet.BackendID]wire.Address
	currencies  []channel.Asset
	channels    chan *PaymentChannel
	Channel     *PaymentChannel
	balance     *big.Int
}

func SetupPaymentClient(
	w *stellarwallet.EphemeralWallet,
	stellarAccount wallet.Account,
	wireAddr wire.Address,

	stellarTokenIDs []channel.Asset,
	bus wire.Bus,
	funder channel.Funder,
	adj channel.Adjudicator,
) (*PaymentClient, error) {
	watcher, err := local.NewWatcher(adj)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}
	// Setup Perun client.
	addresses := map[wallet.BackendID]wire.Address{types.StellarBackendID: wireAddr}
	stellarWallet := map[wallet.BackendID]wallet.Wallet{types.StellarBackendID: w}

	// Setup Accounts
	account := map[wallet.BackendID]wallet.Address{types.StellarBackendID: stellarAccount.Address()}

	perunClient, err := client.New(addresses, bus, funder, adj, stellarWallet, watcher)
	if err != nil {
		return nil, errors.New("creating client")
	}

	c := &PaymentClient{
		perunClient: perunClient,
		account:     account,
		waddress:    addresses,
		currencies:  stellarTokenIDs,
		channels:    make(chan *PaymentChannel, 1),
		balance:     big.NewInt(0),
	}

	go perunClient.Handle(c, c)
	return c, nil
}

// startWatching starts the dispute watcher for the specified channel.
func (c *PaymentClient) startWatching(ch *client.Channel) {
	go func() {
		err := ch.Watch(c)
		if err != nil {
			log.Printf("watcher returned with error: %v", err)
		}
	}()
}

func (c *PaymentClient) OpenChannel(peer map[wallet.BackendID]wire.Address, balances channel.Balances) {
	// We define the channel participants. The proposer has always index 0. Here
	// we use the on-chain addresses as off-chain addresses, but we could also
	// use different ones.

	participants := []map[wallet.BackendID]wire.Address{c.waddress, peer}

	backends := make([]wallet.BackendID, len(c.currencies))
	for i := range c.currencies {
		backends[i] = types.StellarBackendID
	}
	initAlloc := channel.NewAllocation(2, backends, c.currencies...)
	initAlloc.Balances = balances

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
	c.Channel = newPaymentChannel(ch, c.currencies)
}

func (p *PaymentClient) WireAddress() map[wallet.BackendID]wire.Address {
	return p.waddress
}

func (c *PaymentClient) AcceptedChannel() *PaymentChannel {
	return <-c.channels
}

// Shutdown gracefully shuts down the client.
func (c *PaymentClient) Shutdown() {
	c.perunClient.Close()
}
