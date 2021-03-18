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

package perun

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"perun.network/go-perun/backend/ethereum/wallet"
	pwallet "perun.network/go-perun/wallet"
)

func createWallet(privateKeys ...*ecdsa.PrivateKey) *Wallet {
	accounts := make(map[common.Address]*Account)
	for _, k := range privateKeys {
		acc := createAccount(k)
		accounts[acc.EthAddress()] = acc
	}
	return &Wallet{accounts}
}

func createAccount(k *ecdsa.PrivateKey) *Account {
	return &Account{
		Account: accounts.Account{Address: crypto.PubkeyToAddress(k.PublicKey)},
		key:     k,
	}
}

type Wallet struct {
	accounts map[common.Address]*Account
}

type Account struct {
	accounts.Account
	key *ecdsa.PrivateKey
}

type Address wallet.Address

func (w Wallet) Unlock(addr pwallet.Address) (pwallet.Account, error) {
	pa, ok := addr.(*wallet.Address)
	if !ok {
		return nil, errors.New("casting address")
	}

	a, ok := w.accounts[Address(*pa).EthAddress()]
	if !ok {
		return nil, errors.Errorf("getting account: %v", addr)
	}

	return a, nil
}

// LockAll is called by the framework when a Client shuts down.
func (w Wallet) LockAll() {}

// IncrementUsage is called whenever a new channel is created or restored.
func (w Wallet) IncrementUsage(addr pwallet.Address) {}

// DecrementUsage is called whenever a channel is settled.
func (w Wallet) DecrementUsage(addr pwallet.Address) {}

// Contains returns whether the wallet contains the given account.
func (w Wallet) Contains(acc accounts.Account) bool {
	_, ok := w.accounts[acc.Address]
	return ok
}

func (w *Wallet) NewTransactor(account accounts.Account) (*bind.TransactOpts, error) {
	// 1337 is the default chain id for ganache-cli.
	signer := types.NewEIP155Signer(big.NewInt(1337))

	acc, ok := w.accounts[account.Address]
	if !ok {
		return nil, errors.New("account not found")
	}
	return &bind.TransactOpts{
		From: account.Address,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != account.Address {
				return nil, errors.New("not authorized to sign this account")
			}

			signature, err := acc.SignHash(signer.Hash(tx).Bytes())
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}, nil
}

// Address returns the address of this account.
func (a Account) Address() pwallet.Address {
	return wallet.AsWalletAddr(a.Account.Address)
}

// SignData requests a signature from this account.
// It returns the signature or an error.
func (a Account) SignData(data []byte) ([]byte, error) {
	hash := wallet.PrefixedHash(data)
	sig, err := a.SignHash(hash)
	if err != nil {
		return nil, errors.Wrap(err, "SignHash")
	}
	sig[64] += 27
	return sig, nil
}

func (a Account) SignHash(h []byte) ([]byte, error) {
	return crypto.Sign(h, a.key)
}

// EthAccount returns the corresponding Ethereum account.
func (a *Account) EthAccount() accounts.Account {
	return a.Account
}

// EthAddress returns the Ethereum address of this account.
func (a *Account) EthAddress() common.Address {
	return a.Account.Address
}

func (a Address) EthAddress() common.Address {
	return common.Address(a)
}
