// Copyright (c) 2021, PolyCrypt GmbH, Germany. All rights reserved.
// This file is part of perun-tutorial. Use of this source code is
// governed by the Apache 2.0 license that can be found in the LICENSE file.

package main

import (
	"github.com/ethereum/go-ethereum/common"

	"perun.network/go-perun/backend/ethereum/channel"
	"perun.network/go-perun/backend/ethereum/wallet"
)

// config holds the static configuration.
type config struct {
	mnemonic           string          // Ethereum wallet mnemonic.
	chainURL           string          // Url of the Ethereum node.
	bobHost, aliceHost string          // IP:Port of each participant.
	bobAddr, aliceAddr *wallet.Address // On-Chain addresses.
	assetAddr, adjAddr common.Address  // Smart-contract addresses.
}

// cfg provides the configuration globally.
var cfg config

const (
	// The index of each participant in the channel.
	bobIndex, aliceIndex = 0, 1
	// Index of the only Asset that Alice and Bob use.
	assetIdx = 0
)

// init is executed at program start and sets all values of `cfg`.
func init() {
	cfg.mnemonic = "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic"
	cfg.chainURL = "ws://127.0.0.1:8545"
	// Set the IPs,
	cfg.bobHost = "0.0.0.0:5751"
	cfg.aliceHost = "0.0.0.0:5750"
	// Set the on-chain account addresses.
	cfg.bobAddr = wallet.AsWalletAddr(common.HexToAddress("0xF73C1cdA5bD32E6693D3cB313Bd9B9338d96f184"))
	cfg.aliceAddr = wallet.AsWalletAddr(common.HexToAddress("0x2EE1ac154435f542ECEc55C5b0367650d8A5343B"))
	// Set the smart-contract addresses.
	cfg.adjAddr = common.HexToAddress("0x079557d7549d7D44F4b00b51d2C532674129ed51")
	cfg.assetAddr = common.HexToAddress("0x923439be515b6A928cB9650d70000a9044e49E85")
	// Make our lives easier by disabling go-peruns reorg resistance.
	channel.TxFinalityDepth = 1
}
