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

	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"
	"perun.network/perun-examples/app-channel/app"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	txFinalityDepth = 1 // Number of blocks required to confirm a transaction.
)

// AppClient is an app channel client.
type AppClient struct {
	perunClient *client.Client    // The core Perun client.
	account     wallet.Address    // The account we use for on-chain and off-chain transactions.
	currency    channel.Asset     // The currency we expect to get paid in.
	stake       channel.Bal       // The amount we put at stake.
	app         *app.TicTacToeApp // The app definition.
	channels    chan *TicTacToeChannel
}

// SetupAppClient creates a new app client.
func SetupAppClient(
	bus wire.Bus, // bus is used of off-chain communication.
	w *swallet.Wallet, // w is the wallet used for signing transactions.
	acc common.Address, // acc is the address of the account to be used for signing transactions.
	nodeURL string, // nodeURL is the URL of the blockchain node.
	chainID uint64, // chainID is the identifier of the blockchain.
	adjudicator common.Address, // adjudicator is the address of the adjudicator.
	asset ethwallet.Address, // asset is the address of the asset holder for our app channels.
	app *app.TicTacToeApp, // app is the channel app we want to set up the client with.
	stake channel.Bal, // stake is the balance the client is willing to fund the channel with.
) (*AppClient, error) {
	// Create Ethereum client and contract backend.
	cb, err := CreateContractBackend(nodeURL, chainID, w)
	if err != nil {
		return nil, fmt.Errorf("creating contract backend: %w", err)
	}

	// Validate contracts.
	err = ethchannel.ValidateAdjudicator(context.TODO(), cb, adjudicator)
	if err != nil {
		return nil, fmt.Errorf("validating adjudicator: %w", err)
	}
	err = ethchannel.ValidateAssetHolderETH(context.TODO(), cb, common.Address(asset), adjudicator)
	if err != nil {
		return nil, fmt.Errorf("validating adjudicator: %w", err)
	}

	// Setup funder.
	funder := ethchannel.NewFunder(cb)
	dep := ethchannel.NewETHDepositor()
	ethAcc := accounts.Account{Address: acc}
	funder.RegisterAsset(asset, dep, ethAcc)

	// Setup adjudicator.
	adj := ethchannel.NewAdjudicator(cb, adjudicator, acc, ethAcc)

	// Setup dispute watcher.
	watcher, err := local.NewWatcher(adj)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}

	// Setup Perun client.
	waddr := ethwallet.AsWalletAddr(acc)
	perunClient, err := client.New(waddr, bus, funder, adj, w, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	// Create client and start request handler.
	c := &AppClient{
		perunClient: perunClient,
		account:     waddr,
		currency:    &asset,
		stake:       stake,
		app:         app,
		channels:    make(chan *TicTacToeChannel, 1),
	}

	channel.RegisterApp(app)
	go perunClient.Handle(c, c)

	return c, nil
}

// OpenAppChannel opens a new app channel with the specified peer.
func (c *AppClient) OpenAppChannel(peer wire.Address) *TicTacToeChannel {
	participants := []wire.Address{c.account, peer}

	// We create an initial allocation which defines the starting balances.
	initAlloc := channel.NewAllocation(2, c.currency)
	initAlloc.SetAssetBalances(c.currency, []channel.Bal{
		c.stake, // Our initial balance.
		c.stake, // Peer's initial balance.
	})

	// Prepare the channel proposal by defining the channel parameters.
	challengeDuration := uint64(10) // On-chain challenge duration in seconds.

	firstActorIdx := channel.Index(0)
	withApp := client.WithApp(c.app, c.app.InitData(firstActorIdx))

	proposal, err := client.NewLedgerChannelProposal(
		challengeDuration,
		c.account,
		initAlloc,
		participants,
		withApp,
	)
	if err != nil {
		panic(err)
	}

	// Send the app channel proposal.
	ch, err := c.perunClient.ProposeChannel(context.TODO(), proposal)
	if err != nil {
		panic(err)
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)

	return newTicTacToeChannel(ch)
}

// startWatching starts the dispute watcher for the specified channel.
func (c *AppClient) startWatching(ch *client.Channel) {
	go func() {
		err := ch.Watch(c)
		if err != nil {
			fmt.Printf("Watcher returned with error: %v", err)
		}
	}()
}

// AcceptedChannel returns the next accepted app channel.
func (c *AppClient) AcceptedChannel() *TicTacToeChannel {
	return <-c.channels
}

// Shutdown gracefully shuts down the client.
func (c *AppClient) Shutdown() {
	c.perunClient.Close()
}
