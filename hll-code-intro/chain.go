// Copyright (c) 2021, PolyCrypt GmbH, Germany. All rights reserved.
// This file is part of perun-tutorial. Use of this source code is
// governed by the Apache 2.0 license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire/net"
	"perun.network/go-perun/wire/net/simple"
)

// connectToChain connects to the Ethereum blockchain.
func connectToChain(transactor ethchannel.Transactor) ethchannel.ContractBackend {
	client, err := ethclient.Dial(cfg.chainURL)
	noError(err)

	return ethchannel.NewContractBackend(client, transactor)
}

// setupContracts validates the Adjudicator and AssetHolder contracts.
func setupContracts(chain ethchannel.ContractBackend, account accounts.Account) (*ethchannel.Adjudicator, common.Address) {
	err := ethchannel.ValidateAssetHolderETH(context.Background(), chain, cfg.assetAddr, cfg.adjAddr)
	noError(err)
	fmt.Printf("Adjudicator at %v\n AssetHolder at %v\n", cfg.adjAddr, cfg.assetAddr)
	adjudicator := ethchannel.NewAdjudicator(chain, cfg.adjAddr, account.Address, account)

	return adjudicator, cfg.assetAddr
}

// setupNetworking sets up the TCP network connection to Alice.
// go-perun will use this connection to send off-chain payments.
func setupNetworking(account wallet.Account) *net.Bus {
	// Create a Dialer to initiate a connection with Alice.
	dialer := simple.NewTCPDialer(10 * time.Second)
	dialer.Register(cfg.aliceAddr, cfg.aliceHost)

	// Create a Bus which will forward network messages.
	return net.NewBus(account, dialer)
}

// setupFunder creates a Funder to fund channels.
// The passed account will be used to deposit funds with.
func setupFunder(chain ethchannel.ContractBackend, account accounts.Account, assetHolder common.Address) channel.Funder {
	funder := ethchannel.NewFunder(chain)
	// depositor will be used to create Ethereum transactions.
	depositor := new(ethchannel.ETHDepositor)
	// Register the depositor for our asset.
	funder.RegisterAsset(ethwallet.Address(assetHolder), depositor, account)

	return funder
}
