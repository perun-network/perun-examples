// Copyright 2021 PolyCrypt GmbH, Germany
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
	"time"

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
	dialerTimeout   = 10 * time.Second
	txFinalityDepth = 1
)

type Client struct {
	PerunClient     *client.Client
	ContractBackend ethchannel.ContractInterface
	Adjudicator     common.Address
	AccountAddress  wallet.Address
}

func StartClient(
	w *swallet.Wallet,
	acc common.Address,
	nodeURL string,
	chainID uint64,
	adjudicator common.Address,
) (*Client, error) {
	// Create Ethereum client and contract backend.
	cb, err := CreateContractBackend(nodeURL, chainID, w)
	if err != nil {
		return nil, fmt.Errorf("creating contract backend: %w", err)
	}

	// Validate adjudicator.
	err = ethchannel.ValidateAdjudicator(context.TODO(), cb, adjudicator)
	if err != nil {
		return nil, fmt.Errorf("validating adjudicator: %w", err)
	}

	// Setup message bus.
	waddr := ethwallet.AsWalletAddr(acc)
	bus := wire.NewLocalBus()
	// TODO add tutorial that explains tcp/ip bus.

	// Setup funder and adjudicator.
	funder := ethchannel.NewFunder(cb)
	ethAcc := accounts.Account{Address: acc}
	adj := ethchannel.NewAdjudicator(cb, adjudicator, acc, ethAcc)

	// Setup dispute watcher.
	watcher, err := local.NewWatcher(adj)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}

	// Setup Perun client.
	perunClient, err := client.New(waddr, bus, funder, adj, w, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	// Create client and start request handler.
	c := &Client{
		perunClient,
		cb,
		adjudicator,
		waddr,
	}
	go perunClient.Handle(c, c)

	return c, nil
}

func (c *Client) OpenChannel(peer *Client, asset common.Address, amount uint64) Channel {
	// We define the channel participants. The proposer always has index 0. Here
	// we use the on-chain addresses as off-chain addresses, but we could also
	// use different ones.
	participants := []wire.Address{c.AccountAddress, peer.AccountAddress}

	// We verify the asset holder contract.
	err := ethchannel.ValidateAssetHolderETH(context.TODO(), c.ContractBackend, asset, c.Adjudicator)
	if err != nil {
		panic(err)
	}

	// We create an initial allocation which defines the starting balances.
	ethAsset := ethwallet.AsWalletAddr(asset)       // Convert to wallet.Address, which implements Asset. //TODO create ethchannel.AsAsset
	initAlloc := channel.NewAllocation(2, ethAsset) //TODO Create issue: init the balances to zero.
	initAlloc.SetAssetBalances(ethAsset, []channel.Bal{
		new(big.Int).SetUint64(amount), // Our initial balance.
		big.NewInt(0),                  // Peer's initial balance.
	})

	// Prepare the channel proposal by defining the channel parameters.
	challengeDuration := uint64(10) // On-chain challenge duration in seconds.
	proposal, err := client.NewLedgerChannelProposal(
		challengeDuration,
		c.AccountAddress,
		initAlloc,
		participants,
	)
	if err != nil {
		panic(err)
	}

	// Send the proposal.
	ch, err := c.PerunClient.ProposeChannel(context.TODO(), proposal)
	if err != nil {
		panic(err)
	}

	return Channel{ch: ch}
}

type Channel struct {
	ch    *client.Channel
	asset channel.Asset
}

//TODO document all exported functions.
func (c Channel) SendPayment(amount uint64) {
	// Transfer the given amount from us to peer.
	// Use UpdateBy to update the channel state.
	err := c.ch.UpdateBy(context.TODO(), func(state *channel.State) error { //TODO we always use context.TODO for simplicity.
		ethAmount := new(big.Int).SetUint64(amount)
		state.Allocation.TransferBalance(0, 1, c.asset, ethAmount)
		return nil
	})
	if err != nil {
		panic(err) //TODO mention that we always panic on error for simplicity.
	}
}

func (c Channel) Close() {
	// Finalize the channel to enable fast settlement.
	err := c.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		state.IsFinal = true
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Settle concludes the channel and withdraws the funds.
	err = c.ch.Settle(context.TODO(), false)
	if err != nil {
		panic(err)
	}

	// Close frees up channel resources.
	c.ch.Close()
}
