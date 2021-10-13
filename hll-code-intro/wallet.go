// Copyright (c) 2021, PolyCrypt GmbH, Germany. All rights reserved.
// This file is part of perun-tutorial. Use of this source code is
// governed by the Apache 2.0 license that can be found in the LICENSE file.

package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	perunhd "perun.network/go-perun/backend/ethereum/wallet/hd"
)

// setupWallet create a wallet and an account.
// Both will be used by go-perun to sign on- and off-chain data.
func setupWallet() (*perunhd.Account, *perunhd.Wallet) {
	rootWallet, err := hdwallet.NewFromMnemonic(cfg.mnemonic)
	noError(err)

	wallet, err := perunhd.NewWallet(rootWallet, accounts.DefaultBaseDerivationPath.String(), uint(2))
	noError(err)

	acc, err := wallet.NewAccount()
	noError(err)
	return acc, wallet
}

// setupTransactor creates a transactor for signing Ethereum transations.
// This is Ethereum specific and will be different for other chains.
func setupTransactor(wallet *perunhd.Wallet) *perunhd.Transactor {
	// 1337 is the default chain id for ganache-cli.
	signer := types.NewEIP155Signer(big.NewInt(1337))
	return perunhd.NewTransactor(wallet.Wallet(), signer)
}
