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
	"fmt"
	"math/big"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"perun.network/go-perun/backend/ethereum/bindings/assetholdereth"
	"perun.network/perun-examples/payment-channel/eth"
)

type ClientConfig struct {
	PerunClientConfig
	ContextTimeout time.Duration
}

type Client struct {
	role            Role
	PerunClient     *PerunClient
	AssetHolderAddr common.Address
	AssetHolder     *assetholdereth.AssetHolderETH
	ContextTimeout  time.Duration
	Channel         *client.Channel
}

func StartClient(cfg ClientConfig) (*Client, error) {
	perunClient, err := setupPerunClient(cfg.PerunClientConfig)
	if err != nil {
		return nil, errors.WithMessage(err, "creating perun client")
	}

	ah, err := assetholdereth.NewAssetHolderETH(cfg.AssetHolderAddr, perunClient.ContractBackend)
	if err != nil {
		return nil, errors.WithMessage(err, "loading asset holder")
	}

	c := &Client{
		cfg.Role,
		perunClient,
		cfg.AssetHolderAddr,
		ah,
		cfg.ContextTimeout,
		nil,
	}

	go c.PerunClient.StateChClient.Handle(c, c)
	go c.PerunClient.Bus.Listen(c.PerunClient.Listener)

	return c, nil
}

func (c *Client) OpenChannel(opponent wallet.Address) error {
	fmt.Printf("%s: Opening channel from %s to %s\n", c.RoleAsString(), c.RoleAsString(), c.OpponentRoleAsString())
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
	peers := []wire.Address{c.PerunClient.Account.Address(), opponent}

	// Prepare the proposal by defining the channel parameters.
	proposal, err := client.NewLedgerChannelProposal(10, c.PerunAddress(), initBals, peers)
	if err != nil {
		return fmt.Errorf("creating channel proposal: %w", err)
	}
	ctx, cancel := c.defaultContextWithTimeout()
	defer cancel()

	// Send the proposal.
	ch, err := c.PerunClient.StateChClient.ProposeChannel(ctx, proposal)
	c.Channel = ch

	if err != nil {
		return fmt.Errorf("proposing channel: %w", err)
	}

	fmt.Printf("\n 🎉 Opened channel with id 0x%x \n\n", ch.ID())
	return nil
}

func (c *Client) UpdateChannel() error {
	fmt.Printf("%s: Update channel by sending 5 ETH to %s \n", c.RoleAsString(), c.OpponentRoleAsString())

	ctx, cancel := c.defaultContextWithTimeout()
	defer cancel()
	// Use UpdateBy to conveniently update the channels state.
	return c.Channel.UpdateBy(ctx, func(state *channel.State) error {
		// Shift 5 ETH from caller to opponent.
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
