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
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	ewallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/log"
	"perun.network/go-perun/wire"
	"perun.network/perun-collateralized-channels/app"
	"perun.network/perun-collateralized-channels/contracts/generated/collateralAssetHolderETH"
	"perun.network/perun-collateralized-channels/perun"
)

type ClientConfig struct {
	perun.ClientConfig
	ChallengeDuration time.Duration
	AppAddress        common.Address
	ContextTimeout    time.Duration
}

type PaymentAcceptancePolicy = func(
	amount *big.Int,
	collateral *big.Int,
	funding *big.Int,
	balance *big.Int,
	hasOverdrawn bool,
) (ok bool)

type Channel struct {
	*client.Channel
	state *channel.State
}

type Client struct {
	mu                sync.Mutex
	perunClient       *perun.Client
	assetHolderAddr   common.Address
	assetHolder       *collateralAssetHolderETH.CollateralAssetHolderETH
	channels          map[common.Address]*Channel
	challengeDuration time.Duration
	appAddress        common.Address
	contextTimeout    time.Duration
	updatePolicy      PaymentAcceptancePolicy
}

func SetupClient(
	cfg ClientConfig,
	updatePolicy PaymentAcceptancePolicy,
) (*Client, error) {
	perunClient, err := perun.SetupClient(cfg.ClientConfig)
	if err != nil {
		return nil, errors.WithMessage(err, "creating perun client")
	}

	ah, err := collateralAssetHolderETH.NewCollateralAssetHolderETH(cfg.AssetHolderAddr, perunClient.ContractBackend)
	if err != nil {
		return nil, errors.WithMessage(err, "loading asset holder")
	}

	c := &Client{
		sync.Mutex{},
		perunClient,
		cfg.AssetHolderAddr,
		ah,
		make(map[common.Address]*Channel),
		cfg.ChallengeDuration,
		cfg.AppAddress,
		cfg.ContextTimeout,
		updatePolicy,
	}

	collateralApp := app.NewCollateralApp(ewallet.AsWalletAddr(cfg.AppAddress))
	channel.RegisterApp(collateralApp)

	go c.perunClient.PerunClient.Handle(c, c)
	go c.perunClient.Bus.Listen(c.perunClient.Listener)

	return c, nil
}

func (c *Client) IncreaseCollateral(amount *big.Int) (err error) {
	return c.transactAndConfirm(func(tr *bind.TransactOpts) (*types.Transaction, error) {
		tr.Value = amount
		fundingID := calcPeerCollateralID(c.Address())
		return c.assetHolder.Deposit(tr, fundingID, amount)
	})
}

func (c *Client) SendPayment(peer common.Address, amount *big.Int) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.channels[peer]
	if !ok {
		ctx := c.defaultContext()
		challengeDuration := uint64(c.challengeDuration.Seconds())
		collateralApp := app.NewCollateralApp(ewallet.AsWalletAddr(c.appAddress))
		peers := []wire.Address{c.perunClient.Account.Address(), ewallet.AsWalletAddr(peer)}
		withApp := client.WithApp(collateralApp, collateralApp.ZeroBalance())
		prop, err := client.NewLedgerChannelProposal(challengeDuration, c.PerunAddress(), zeroBalance(c.assetHolderAddr), peers, withApp)
		if err != nil {
			return errors.WithMessage(err, "creating channel proposal")
		}
		perunChannel, err := c.perunClient.PerunClient.ProposeChannel(ctx, prop)
		if err != nil {
			return errors.WithMessage(err, "proposing channel")
		}
		ch = c.onNewChannel(perunChannel)
	}

	time.Sleep(1000 * time.Millisecond) // have to wait until other side has set up the channel, can be removed once updated to new go-perun version

	return ch.UpdateBy(c.defaultContext(), func(s *channel.State) error {
		return app.Transfer(ch.Params().Parts, s, c.Address(), peer, amount)
	})
}

func (c *Client) onNewChannel(perunChannel *client.Channel) *Channel {
	ch := &Channel{perunChannel, perunChannel.State().Clone()}
	peer := peerAddress(ch.Channel)
	c.channels[peer] = ch
	// Start the on-chain watcher.
	go func() {
		err := ch.Watch(c)
		log.Printf("Watcher for channel %v returned with: %v\n", ch, err)
	}()
	return ch

}

func (c *Client) IncreaseChannelCollateral(peer common.Address, amount *big.Int) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.channels[peer]
	if !ok {
		return errors.New("channel not found")
	}

	return c.transactAndConfirm(func(tr *bind.TransactOpts) (*types.Transaction, error) {
		tr.Value = amount
		fundingID := calcChannelCollateralID(ch.ID(), c.Address())
		return c.assetHolder.Deposit(tr, fundingID, amount)
		// return c.assetHolder.DepositIntoChannelCollateral(tr, ch.ID(), c.Address(), amount)
	})
}

func (c *Client) Settle(peer common.Address) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.channels[peer]
	if !ok {
		return errors.Errorf("channel not found for peer: %s", peer.Hex())
	}

	err = c.settle(ch)
	if err != nil {
		return errors.WithMessage(err, "settling channel")
	}

	// Create withdrawal authentication
	amount, err := app.ChannelBalance(ch.Params().Parts, ch.State().Data, c.Address())
	if err != nil {
		return errors.WithMessage(err, "getting balance")
	}
	auth, sig, err := perun.NewWithdrawalAuth(ch.ID(), c.perunClient.Account, c.PerunAddress(), amount)
	if err != nil {
		return errors.WithMessage(err, "creating withdrawal authentication")
	}

	// Withdraw channel collateral
	return c.transactAndConfirm(func(tr *bind.TransactOpts) (*types.Transaction, error) {
		return c.assetHolder.Withdraw(tr, auth, sig)
	})
}
