// Copyright 2021 PolyCrypt GmbH, Germany
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

// Package client implements the client side of the collateralized channel.
package client

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	ethwire "github.com/perun-network/perun-eth-backend/wire"
	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/log"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"
	"perun.network/perun-collateralized-channels/app"
	"perun.network/perun-collateralized-channels/contracts/generated/collateralAssetHolderETH"
)

const (
	chainID         = 1337 // default chainID of Ganache
	txFinalityDepth = 1    // Number of blocks required to confirm a transaction.
)

// PaymentAcceptancePolicy is a policy that decides whether a payment should be accepted.
type PaymentAcceptancePolicy = func(
	amount *big.Int,
	collateral *big.Int,
	funding *big.Int,
	balance *big.Int,
	hasOverdrawn bool,
) (ok bool)

// Channel represents an app channel.
type Channel struct {
	*client.Channel
	state *channel.State
}

// AppClient is an app channel client.
type AppClient struct {
	perunClient       *client.Client
	account           wallet.Address
	walletAccount     wallet.Account
	transactor        *swallet.Transactor
	ethAcc            accounts.Account
	waddress          wire.Address
	ethClient         *ethclient.Client
	assetHolderAddr   common.Address
	asset             channel.Asset
	assetHolder       *collateralAssetHolderETH.CollateralAssetHolderETH
	contractBackend   ethchannel.ContractInterface
	channels          map[common.Address]*Channel
	challengeDuration time.Duration
	collateralApp     *app.CollateralApp
	contextTimeout    time.Duration
	updatePolicy      PaymentAcceptancePolicy
}

// SetupAppClient creates a new app channel client.
func SetupAppClient(
	bus wire.Bus,
	w *swallet.Wallet,
	acc common.Address,
	assetHolderAddr common.Address,
	eaddress *ethwallet.Address, // eaddress is the address of the Ethereum account to be used for signing transactions.
	nodeURL string, // nodeURL is the URL of the blockchain node.
	chainID uint64, // chainID is the identifier of the blockchain.
	app *app.CollateralApp,
	updatePolicy PaymentAcceptancePolicy,
	challengeDuration time.Duration,
	contextTimeout time.Duration,

) (*AppClient, error) {
	// Create Ethereum client and contract backend.
	ethClient, cb, tr, err := CreateContractBackend(nodeURL, chainID, w)
	if err != nil {
		return nil, errors.WithMessage(err, "creating contract backend")
	}

	walletAcc, err := w.Unlock(eaddress)
	if err != nil {
		return nil, errors.WithMessage(err, "unlocking wallet")
	}

	// Setup funder.
	funder := ethchannel.NewFunder(cb)
	dep := ethchannel.NewETHDepositor()
	ethAcc := accounts.Account{Address: acc}
	asset := ethchannel.NewAsset(big.NewInt(int64(chainID)), assetHolderAddr)
	funder.RegisterAsset(*asset, dep, ethAcc)

	// Setup adjudicator.
	adjudicator := ethchannel.NewAdjudicator(cb, assetHolderAddr, ethAcc.Address, ethAcc)

	// Setup dispute watcher.
	watcher, err := local.NewWatcher(adjudicator)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}

	// Setup client.
	waddr := &ethwire.Address{Address: eaddress}
	perunClient, err := client.New(waddr, bus, funder, adjudicator, w, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	ah, err := collateralAssetHolderETH.NewCollateralAssetHolderETH(assetHolderAddr, cb)
	if err != nil {
		return nil, errors.WithMessage(err, "loading asset holder")
	}

	// Create client and start request handler.
	c := &AppClient{
		perunClient:       perunClient,
		account:           eaddress,
		walletAccount:     walletAcc,
		transactor:        tr,
		ethAcc:            ethAcc,
		waddress:          waddr,
		ethClient:         ethClient,
		assetHolderAddr:   assetHolderAddr,
		asset:             asset,
		assetHolder:       ah,
		contractBackend:   cb,
		channels:          make(map[common.Address]*Channel),
		challengeDuration: challengeDuration,
		collateralApp:     app,
		contextTimeout:    contextTimeout,
		updatePolicy:      updatePolicy,
	}

	channel.RegisterApp(app)

	go perunClient.Handle(c, c)
	return c, nil
}

// IncreaseCollateral increases the collateral of the client.
func (c *AppClient) IncreaseCollateral(amount *big.Int) (err error) {
	return c.transactAndConfirm(func(tr *bind.TransactOpts) (*types.Transaction, error) {
		tr.Value = amount
		fundingID := calcPeerCollateralID(c.WalletAddress())
		return c.assetHolder.Deposit(tr, fundingID, amount)
	})
}

// SendPayment sends a payment to the peer.
func (c *AppClient) SendPayment(peer common.Address, amount *big.Int) (err error) {
	ch, ok := c.channels[peer]
	if !ok {
		ctx := c.defaultContext()
		challengeDuration := uint64(c.challengeDuration.Seconds())

		participants := []wire.Address{c.waddress, &ethwire.Address{Address: ethwallet.AsWalletAddr(peer)}}
		fmt.Println("Participants: ", participants)
		if err != nil {
			return errors.WithMessage(err, "creating peers")
		}
		withApp := client.WithApp(c.collateralApp, c.collateralApp.ZeroBalance())

		prop, err := client.NewLedgerChannelProposal(challengeDuration, c.account, zeroBalance(c.asset), participants, withApp)
		if err != nil {
			return errors.WithMessage(err, "creating channel proposal")
		}
		perunChannel, err := c.perunClient.ProposeChannel(ctx, prop)
		if err != nil {
			return errors.WithMessage(err, "proposing channel")
		}
		ch = c.onNewChannel(perunChannel)
	}

	time.Sleep(1000 * time.Millisecond) // have to wait until other side has set up the channel, can be removed once updated to new go-perun version

	return ch.Update(c.defaultContext(), func(s *channel.State) {
		app.Transfer(ch.Params().Parts, s, c.WalletAddress(), peer, amount)
	})
}

func (c *AppClient) onNewChannel(perunChannel *client.Channel) *Channel {
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

// IncreaseChannelCollateral increases the collateral of the channel with the peer.
func (c *AppClient) IncreaseChannelCollateral(peer common.Address, amount *big.Int) (err error) {
	ch, ok := c.channels[peer]
	if !ok {
		return errors.WithMessage(err, "getting channel")
	}

	return c.transactAndConfirm(func(tr *bind.TransactOpts) (*types.Transaction, error) {
		tr.Value = amount
		fundingID := calcChannelCollateralID(ch.ID(), c.WalletAddress())
		return c.assetHolder.Deposit(tr, fundingID, amount)
		// return c.assetHolder.DepositIntoChannelCollateral(tr, ch.ID(), c.Address(), amount)
	})
}

// Settle settles the channel with the peer.
func (c *AppClient) Settle(peer common.Address) (err error) {
	ch, ok := c.channels[peer]
	if !ok {
		return errors.WithMessage(err, "getting channel")
	}

	err = c.settle(ch)
	if err != nil {
		return errors.WithMessage(err, "settling channel")
	}

	// Create withdrawal authentication
	amount, err := app.ChannelBalance(ch.Params().Parts, ch.State().Data, c.WalletAddress())
	if err != nil {
		return errors.WithMessage(err, "getting balance")
	}
	auth, sig, err := NewWithdrawalAuth(ch.ID(), c.walletAccount, c.walletAccount.Address(), amount)
	if err != nil {
		return errors.WithMessage(err, "creating withdrawal authentication")
	}

	// Withdraw channel collateral
	return c.transactAndConfirm(func(tr *bind.TransactOpts) (*types.Transaction, error) {
		return c.assetHolder.Withdraw(tr, auth, sig)
	})
}
