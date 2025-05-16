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

package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	perunhd "perun.network/go-perun/backend/ethereum/wallet/hd"
)

func setupWallet(role Role) (*perunhd.Account, *perunhd.Wallet, error) {
	rootWallet, err := hdwallet.NewFromMnemonic(cfg.mnemonic)
	if err != nil {
		return nil, nil, fmt.Errorf("creating hd wallet: %w", err)
	}
	// Alice has account index 0 and Bob 1.
	wallet, err := perunhd.NewWallet(rootWallet, "m/44'/60'/0'/0/0", uint(role))
	if err != nil {
		return nil, nil, fmt.Errorf("deriving path: %w", err)
	}
	// Derive the first account from the wallet.
	acc, err := wallet.NewAccount()
	if err != nil {
		return nil, nil, fmt.Errorf("deriving hd account: %w", err)
	}

	return acc, wallet, nil
}

func createTransactor(wallet *perunhd.Wallet) *perunhd.Transactor {
	// 1337 is the default chain id for ganache-cli.
	signer := types.NewEIP155Signer(big.NewInt(1337))
	return perunhd.NewTransactor(wallet.Wallet(), signer)
}
