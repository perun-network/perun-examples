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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"perun.network/go-perun/backend/ethereum/bindings/assetholdereth"
	ewallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wire"
	"perun.network/perun-examples/app-channel/app"
	"perun.network/perun-examples/app-channel/perun"
)

type ClientConfig struct {
	perun.ClientConfig
	ChallengeDuration time.Duration
	AppAddress        common.Address
	ContextTimeout    time.Duration
}

type Client struct {
	sync.Mutex
	perunClient       *perun.Client
	assetHolderAddr   common.Address
	assetHolder       *assetholdereth.AssetHolderETH
	games             map[channel.ID]*Game
	challengeDuration time.Duration
	appAddress        common.Address
	contextTimeout    time.Duration
	gameProposals     chan *GameProposal
}

func StartClient(cfg ClientConfig) (*Client, error) {
	perunClient, err := perun.SetupClient(cfg.ClientConfig)
	if err != nil {
		return nil, errors.WithMessage(err, "creating perun client")
	}

	ah, err := assetholdereth.NewAssetHolderETH(cfg.AssetHolderAddr, perunClient.ContractBackend)
	if err != nil {
		return nil, errors.WithMessage(err, "loading asset holder")
	}

	c := &Client{
		sync.Mutex{},
		perunClient,
		cfg.AssetHolderAddr,
		ah,
		make(map[channel.ID]*Game),
		cfg.ChallengeDuration,
		cfg.AppAddress,
		cfg.ContextTimeout,
		make(chan *GameProposal, 1),
	}

	_app := app.NewTicTacToeApp(ewallet.AsWalletAddr(cfg.AppAddress))
	channel.RegisterApp(_app)

	go c.perunClient.PerunClient.Handle(c, c)
	go c.perunClient.Bus.Listen(c.perunClient.Listener)

	return c, nil
}

func (c *Client) ProposeGame(opponent common.Address, stake *big.Int) (*Game, error) {
	c.Lock()
	defer c.Unlock()

	ctx, cancel := c.defaultContextWithTimeout()
	defer cancel()
	_app := app.NewTicTacToeApp(ewallet.AsWalletAddr(c.appAddress))
	peers := []wire.Address{c.perunClient.Account.Address(), ewallet.AsWalletAddr(opponent)}
	withApp := client.WithApp(_app, _app.InitData(0))

	prop, err := client.NewLedgerChannelProposal(
		c.challengeDurationInSeconds(),
		c.PerunAddress(),
		makeStakeAllocation(c.assetHolderAddr, stake),
		peers,
		withApp,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "creating channel proposal")
	}

	perunChannel, err := c.perunClient.PerunClient.ProposeChannel(ctx, prop)
	if err != nil {
		return nil, errors.WithMessage(err, "proposing channel")
	}
	g := c.newGame(perunChannel)
	return g, nil
}

func (c *Client) NextGameProposal() (*GameProposal, error) {
	p, ok := <-c.gameProposals
	if !ok {
		return nil, fmt.Errorf("channel closed")
	}
	return p, nil
}

func (c *Client) newGame(perunChannel *client.Channel) *Game {
	g := &Game{
		sync.Mutex{},
		c,
		perunChannel,
		perunChannel.State().Clone(),
		make(chan error, 1),
	}
	c.games[perunChannel.ID()] = g
	// Start the on-chain watcher.
	go func() {
		g.errs <- g.ch.Watch(c)
	}()
	return g
}
