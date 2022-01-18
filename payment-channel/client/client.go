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
	"math/big"
	"sync"

	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
	"perun.network/go-perun/channel"
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
	proposerIdx     = 0 // Participant index of the proposer.  //TODO:go-perun expose channel.ProposerIdx and ReceiverIdx.
	receiverIdx     = 1 // Participant index of the receiver.
)

type Client struct { //TODO:code add coments to variables?
	Name        string
	perunClient *client.Client
	account     wallet.Address
	channels    map[channel.ID]*Channel
	channelsMtx sync.RWMutex
	currency    channel.Asset // The currency we expect to get paid in.
}

func StartClient(
	name string,
	bus wire.Bus,
	w *swallet.Wallet,
	acc common.Address,
	nodeURL string,
	chainID uint64,
	adjudicator common.Address,
	assetHolder common.Address,
) (*Client, error) {
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
	err = ethchannel.ValidateAssetHolderETH(context.TODO(), cb, assetHolder, adjudicator)
	if err != nil {
		return nil, fmt.Errorf("validating adjudicator: %w", err)
	}

	// Setup funder.
	funder := ethchannel.NewFunder(cb)
	asset := *NewAsset(assetHolder)
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
	c := &Client{
		Name:        name,
		perunClient: perunClient,
		account:     waddr,
		channels:    map[channel.ID]*Channel{},
		currency:    &asset,
	}
	go perunClient.Handle(c, c)

	return c, nil
}

func (c *Client) OpenChannel(peer *Client, asset channel.Asset, amount uint64) Channel {
	// We define the channel participants. The proposer always has index 0. Here
	// we use the on-chain addresses as off-chain addresses, but we could also
	// use different ones.
	participants := []wire.Address{c.account, peer.account}

	// We create an initial allocation which defines the starting balances.
	initAlloc := channel.NewAllocation(2, asset) //TODO:go-perun balances should be initialized to zero
	initAlloc.SetAssetBalances(asset, []channel.Bal{
		new(big.Int).SetUint64(amount), // Our initial balance.
		big.NewInt(0),                  // Peer's initial balance.
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

	return *newChannel(ch)
}

func (c *Client) Shutdown() {
	c.perunClient.Close()
}
