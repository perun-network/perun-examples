// Copyright 2025 PolyCrypt GmbH
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
	simplewallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	ethwire "github.com/perun-network/perun-eth-backend/wire"
	"log"
	"math/big"
	"perun.network/go-perun/channel/multi"
	"perun.network/go-perun/wire/net/simple"
	"perun.network/perun-ckb-backend/channel/adjudicator"
	"perun.network/perun-ckb-backend/channel/asset"
	"perun.network/perun-ckb-backend/channel/funder"
	ckbaddress "perun.network/perun-ckb-backend/wallet/address"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	ckbwallet "perun.network/perun-ckb-backend/wallet"
)

// PaymentClient is a payment channel client.
type PaymentClient struct {
	perunClient *client.Client                      // The core Perun client.
	account     map[wallet.BackendID]wallet.Address // The account we use for on-chain and off-chain transactions.
	waddress    map[wallet.BackendID]wire.Address
	currency    []channel.Asset      // The currency we expect to get paid in.
	channels    chan *PaymentChannel // Accepted payment channels.
}

// SetupPaymentClient creates a new payment client.
func SetupPaymentClient(
	bus wire.Bus, // bus is used of off-chain communication.
	ethWallet *simplewallet.Wallet, // w is the wallet used for signing ethereum transactions.
	acc common.Address, // acc is the address of the account to be used for signing transactions.
	ethAddress *ethwallet.Address, // ethAddress is the address of the Ethereum account to be used for signing transactions.
	nodeURL string, // nodeURL is the URL of the blockchain node.
	chainID uint64, // chainID is the identifier of the blockchain.
	adjudicator common.Address, // adjudicator is the address of the adjudicator.
	assetAddr ethwallet.Address, // asset is the address of the asset holder for our payment channels.
	ckbWallet *ckbwallet.EphemeralWallet, // stellarWallet is the wallet used for signing stellar transactions.
	ckbAccount *ckbwallet.Account, // stellarAccount is the account to be used for signing Stellar transactions.
	ckbTokenIDs channel.Asset, // stellarTokenIDs is the list of token IDs to be used for payment channels.
	ckbFunder *funder.Funder, // stellarFunder is the funder to be used for funding Stellar payment channels.
	ckbAdj *adjudicator.Adjudicator, // stellarAdj is the adjudicator to be used for Stellar payment channels.
) (*PaymentClient, error) {
	multiAdjudicator := multi.NewAdjudicator()
	watcher, err := local.NewWatcher(multiAdjudicator)
	multiFunder := multi.NewFunder()
	ccWallet := map[wallet.BackendID]wallet.Wallet{1: ethWallet, 3: ckbWallet}

	ckbPart := ckbAccount.Address().(*ckbaddress.Participant)
	ckbWireAddr := simple.NewAddress(ckbPart.String())
	// Create Ethereum client and contract backend.
	cb, err := CreateContractBackend(nodeURL, chainID, ethWallet)
	if err != nil {
		return nil, fmt.Errorf("creating contract backend: %w", err)
	}

	// Validate contracts.
	err = ethchannel.ValidateAdjudicator(context.TODO(), cb, adjudicator)
	if err != nil {
		return nil, fmt.Errorf("validating adjudicator: %w", err)
	}
	err = ethchannel.ValidateAssetHolderETH(context.TODO(), cb, common.Address(assetAddr), adjudicator)
	if err != nil {
		return nil, fmt.Errorf("validating adjudicator: %w", err)
	}

	// Setup funder.
	ethFunder := ethchannel.NewFunder(cb)
	ethAssetID := ethchannel.MakeLedgerBackendID(big.NewInt(int64(chainID)))
	ckbAssetID := asset.MakeCCID(asset.MakeContractID("03"))
	multiFunder.RegisterFunder(ethAssetID, ethFunder)
	multiFunder.RegisterFunder(ckbAssetID, ckbFunder)
	dep := ethchannel.NewETHDepositor(50000)
	ethAcc := accounts.Account{Address: acc}
	asset := ethchannel.NewAsset(big.NewInt(int64(chainID)), common.Address(assetAddr))
	ethFunder.RegisterAsset(*asset, dep, ethAcc)

	// Setup adjudicator.
	ethAdj := ethchannel.NewAdjudicator(cb, adjudicator, acc, ethAcc, 1000000)
	multiAdjudicator.RegisterAdjudicator(ethAssetID, ethAdj)
	multiAdjudicator.RegisterAdjudicator(ckbAssetID, ckbAdj)

	// Setup Perun client.
	ethWireAddr := &ethwire.Address{Address: ethAddress}
	addresses := map[wallet.BackendID]wire.Address{1: ethWireAddr, 3: ckbWireAddr}
	perunClient, err := client.New(addresses, bus, multiFunder, multiAdjudicator, ccWallet, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	// Setup Accounts
	account := map[wallet.BackendID]wallet.Address{1: ethAddress, 3: ckbAccount.Address()}

	// Create client and start request handler.
	c := &PaymentClient{
		perunClient: perunClient,
		account:     account,
		waddress:    addresses,
		currency:    []channel.Asset{asset, ckbTokenIDs},
		channels:    make(chan *PaymentChannel, 1),
	}
	go perunClient.Handle(c, c)

	return c, nil
}

// OpenChannel opens a new channel with the specified peer and funding.
func (c *PaymentClient) OpenChannel(peer map[wallet.BackendID]wire.Address, ethAmount float64, ckbAmount uint64) *PaymentChannel {
	// We define the channel participants. The proposer has always index 0. Here
	// we use the on-chain addresses as off-chain addresses, but we could also
	// use different ones.
	participants := []map[wallet.BackendID]wire.Address{c.waddress, peer}

	// We create an initial allocation which defines the starting balances.
	initAlloc := channel.NewAllocation(2, []wallet.BackendID{1, 3}, c.currency[0], c.currency[1])
	shannonAmount := CKByteToShannon(big.NewFloat(float64(ckbAmount)))
	log.Println("ETH amount: ", ethAmount, c.currency[0])
	log.Println("CKB amount: ", shannonAmount, c.currency[1])
	initAlloc.SetAssetBalances(c.currency[0], []channel.Bal{
		EthToWei(big.NewFloat(ethAmount)), // Our initial balance.
		big.NewInt(0),                     // Peer's initial balance.
	})
	initAlloc.SetAssetBalances(c.currency[1], []channel.Bal{
		big.NewInt(8000000000),                                 // Our initial balance.
		big.NewInt(int64(8000000000 + shannonAmount.Uint64())), // Peer's initial balance.
	})

	// Prepare the channel proposal by defining the channel parameters.
	challengeDuration := uint64(1000) // On-chain challenge duration in seconds.
	log.Println("Creating channel proposal")
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
	log.Println("Sending channel proposal", proposal)
	// ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	// defer cancel()
	ch, err := c.perunClient.ProposeChannel(context.TODO(), proposal)
	if err != nil {
		panic(err)
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	log.Println("Starting dispute watcher", ch.ID())
	c.startWatching(ch)

	return newPaymentChannel(ch, c.currency)
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
	log.Println("Waiting for accepted channel", c.channels)
	return <-c.channels
}

// Shutdown gracefully shuts down the client.
func (c *PaymentClient) Shutdown() {
	c.perunClient.Close()
}

func (c *PaymentClient) Addresses() map[wallet.BackendID]wallet.Address {
	return c.account
}
