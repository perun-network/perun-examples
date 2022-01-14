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
	"perun.network/go-perun/wire/net"
	"perun.network/go-perun/wire/net/simple"
	snet "perun.network/go-perun/wire/net/simple"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	dialerTimeout   = 10 * time.Second
	txFinalityDepth = 1
)

type Network struct {
	Dialer net.Dialer
	Bus    *net.Bus
}

type Client struct {
	PerunClient     *client.Client
	Dialer          net.Dialer
	ContractBackend ethchannel.ContractInterface

	// Variables set at runtime.
	Channel *client.Channel
}

func StartClient(
	host string,
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

	// Setup network environment.
	waddr := ethwallet.AsWalletAddr(acc)
	wireAcc := dummyAccount{waddr}
	dialer := simple.NewTCPDialer(dialerTimeout)
	listener, err := snet.NewTCPListener(host)
	if err != nil {
		return nil, fmt.Errorf("creating listener: %w", err)
	}
	bus := net.NewBus(wireAcc, dialer)

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

	c := &Client{perunClient, dialer, cb, nil}

	go c.PerunClient.Handle(c, c)
	go bus.Listen(listener)

	return c, nil
}

func (c *Client) OpenChannel(peer wallet.Address) error {
	fmt.Printf("%s: Opening channel from %s to %s\n", c.RoleAsString(), c.RoleAsString(), c.PeerRoleAsString())
	// Alice and Bob will both start with 10 ETH.
	initBal := eth.EthToWei(big.NewFloat(10))
	// Perun needs an initial allocation which defines the balances of all
	// participants. The same structure is used for multi-asset channels.
	initBals := &channel.Allocation{
		Assets:   []channel.Asset{ethwallet.AsWalletAddr(c.AssetHolderAddr)},
		Balances: [][]*big.Int{{initBal, initBal}},
	}
	// All perun identities that we want to open a channel with. In this case
	// we use the same on- and off-chain accounts but you could use different.
	peers := []wire.Address{c.Addr, peer}

	// Prepare the proposal by defining the channel parameters.
	proposal, err := client.NewLedgerChannelProposal(10, c.PerunAddress(), initBals, peers)
	if err != nil {
		return fmt.Errorf("creating channel proposal: %w", err)
	}

	// Send the proposal.
	ch, err := c.PerunClient.ProposeChannel(context.TODO(), proposal)
	c.Channel = ch

	if err != nil {
		return fmt.Errorf("proposing channel: %w", err)
	}

	fmt.Printf("\n 🎉 Opened channel with id 0x%x \n\n", ch.ID())
	return nil
}

func (c *Client) UpdateChannel() error {
	fmt.Printf("%s: Update channel by sending 5 ETH to %s \n", c.RoleAsString(), c.PeerRoleAsString())

	// Use UpdateBy to conveniently update the channels state.
	return c.Channel.UpdateBy(context.TODO(), func(state *channel.State) error {
		// Shift 5 ETH from caller to peer.
		amount := eth.EthToWei(big.NewFloat(5))
		state.Balances[0][1-c.role].Sub(state.Balances[0][1-c.role], amount)
		state.Balances[0][c.role].Add(state.Balances[0][c.role], amount)
		// Finalize the channel, this will be important in the next step.
		state.IsFinal = true
		return nil
	})
}

func (c *Client) CloseChannel() error {
	fmt.Printf("%s: Close Channel \n", c.RoleAsString())

	ctx, cancel := c.defaultContextWithTimeout()
	defer cancel()

	// .Settle() "closes" the channel (= concludes the channel and withdraws the funds)
	if err := c.Channel.Settle(ctx, false); err != nil {
		return fmt.Errorf("settling channel: %w", err)
	}

	// .Close() terminates the local channel object (frees resources) and has nothing to do with the go-perun channel protocol.
	if err := c.Channel.Close(); err != nil {
		return fmt.Errorf("closing channel object: %w", err)
	}
	return nil
}
