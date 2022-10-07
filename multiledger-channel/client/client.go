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

	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	ethwire "github.com/perun-network/perun-eth-backend/wire"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/channel/multi"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	txFinalityDepth = 1 // Number of blocks required to confirm a transaction.
)

// ChainConfig is used to hold all information needed about a specific chain.
type ChainConfig struct {
	ChainID     ethchannel.ChainID
	ChainURL    string
	Token       common.Address // The address of the deployed ERC20 token.
	Adjudicator common.Address // The address of the deployed Adjudicator contract.
	AssetHolder common.Address // The address of the deployed AssetHolder contract.
}

// PaymentClient is a payment channel client.
type PaymentClient struct {
	perunClient *client.Client       // The core Perun client.
	account     wallet.Address       // The account we use for on-chain and off-chain transactions.
	currencies  [2]channel.Asset     // The currencies of the different chains we support.
	channels    chan *PaymentChannel // Accepted payment channels.
}

// SetupPaymentClient creates a new payment client.
func SetupPaymentClient(
	bus wire.Bus, // bus is used of off-chain communication.
	w *swallet.Wallet, // w is the wallet used for signing transactions.
	acc common.Address, // acc is the address of the account to be used for signing transactions.
	chains [2]ChainConfig, // chains represent the two chains the client should be able to use.
) (*PaymentClient, error) {
	// The multi-funder and multi-adjudicator will be registered with a funder /
	// adjudicators for each chain.
	multiFunder := multi.NewFunder()
	multiAdjudicator := multi.NewAdjudicator()

	var assets [2]channel.Asset

	for i, chain := range chains {
		// Create Ethereum client and contract backend.
		cb, err := CreateContractBackend(chain.ChainURL, chain.ChainID.Int, w)
		if err != nil {
			return nil, fmt.Errorf("creating contract backend: %w", err)
		}

		// Validate contracts.
		err = ethchannel.ValidateAdjudicator(context.TODO(), cb, chain.Adjudicator)
		if err != nil {
			return nil, fmt.Errorf("validating adjudicator: %w", err)
		}
		err = ethchannel.ValidateAssetHolderERC20(context.TODO(), cb, chain.AssetHolder, chain.Adjudicator, chain.Token)
		if err != nil {
			return nil, fmt.Errorf("validating adjudicator: %w", err)
		}

		// Setup funder.
		funder := ethchannel.NewFunder(cb)
		// Register the asset on the funder.
		dep := ethchannel.NewERC20Depositor(chain.Token)
		ethAcc := accounts.Account{Address: acc}
		assets[i] = ethchannel.NewAsset(chain.ChainID.Int, chain.AssetHolder)
		funder.RegisterAsset(*assets[i].(*ethchannel.Asset), dep, ethAcc)
		// Register the funder on the multi-funder.
		multiFunder.RegisterFunder(chain.ChainID, funder)

		// Setup adjudicator.
		adj := ethchannel.NewAdjudicator(cb, chain.Adjudicator, acc, ethAcc)
		// Register the adjudicator on the multi-adjudicator.
		multiAdjudicator.RegisterAdjudicator(chain.ChainID, adj)
	}

	// Setup dispute watcher.
	watcher, err := local.NewWatcher(multiAdjudicator)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}

	// Setup Perun client.
	walletAddr := ethwallet.AsWalletAddr(acc)
	wireAddr := &ethwire.Address{Address: walletAddr}
	perunClient, err := client.New(wireAddr, bus, multiFunder, multiAdjudicator, w, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	// Create client and start request handler.
	c := &PaymentClient{
		perunClient: perunClient,
		account:     walletAddr,
		currencies:  assets,
		channels:    make(chan *PaymentChannel, 1),
	}
	go perunClient.Handle(c, c)

	return c, nil
}

// OpenChannel opens a new channel with the specified peer and funding.
func (c *PaymentClient) OpenChannel(peer wire.Address, balances channel.Balances) *PaymentChannel {
	// We define the channel participants. The proposer has always index 0. Here
	// we use the on-chain addresses as off-chain addresses, but we could also
	// use different ones.
	wireAddr := &ethwire.Address{Address: c.account.(*ethwallet.Address)}
	participants := []wire.Address{wireAddr, peer}

	// We create an initial allocation which defines the starting balances.
	initAlloc := channel.NewAllocation(2, c.currencies[0], c.currencies[1])
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

	return newPaymentChannel(ch, c.currencies)
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
