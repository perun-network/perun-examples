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
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"perun.network/go-perun/backend/ethereum/channel"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"
	"perun.network/go-perun/wire/net"
	"perun.network/go-perun/wire/net/simple"
)

const (
	dialerTimeout   = 10 * time.Second
	txFinalityDepth = 1
)

type PeerWithAddress struct {
	Peer    wire.Address
	Address string
}

type ContractAddresses struct {
	AdjudicatorAddr, AssetHolderAddr common.Address
}

type PerunClientConfig struct {
	EthNodeURL string
	Wallet     wallet.Wallet
	Account    wallet.Account
	Host       string
	Contracts  ContractAddresses
}

type Network struct {
	Dialer net.Dialer
	Bus    *net.Bus
}

type PerunClient struct {
	EthClient       *ethclient.Client
	StateChClient   *client.Client
	Net             Network
	ContractBackend channel.ContractInterface
}

func setupPerunClient(
	host string,
	w *swallet.Wallet,
	acc common.Address,
	nodeURL string,
	chainID uint64,
	adjudicator common.Address,
) (*PerunClient, error) {
	// Create Ethereum client and contract backend.
	ethClient, cb, err := createContractBackend(nodeURL, chainID, w)
	if err != nil {
		return nil, errors.WithMessage(err, "creating contract backend")
	}

	// Setup network environment.
	waddr := ethwallet.AsWalletAddr(acc)
	wireAcc := dummyAccount{waddr}
	dialer := simple.NewTCPDialer(dialerTimeout)
	bus := net.NewBus(wireAcc, dialer)
	network := Network{
		Dialer: dialer,
		Bus:    bus,
	}

	// Setup funder and adjudicator.
	funder := channel.NewFunder(cb)
	ethAcc := accounts.Account{Address: acc}
	adj := channel.NewAdjudicator(cb, adjudicator, acc, ethAcc)

	// Setup dispute watcher.
	watcher, err := local.NewWatcher(adj)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}

	// Setup Perun client.
	c, err := client.New(waddr, bus, funder, adj, w, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	return &PerunClient{ethClient, c, network, cb}, nil
}

func createContractBackend(
	nodeURL string,
	chainID uint64,
	w *swallet.Wallet,
) (*ethclient.Client, channel.ContractBackend, error) {
	signer := types.NewEIP155Signer(new(big.Int).SetUint64(chainID))
	transactor := swallet.NewTransactor(w, signer) //TODO transactor should be spawnable from Wallet: Add method "NewTransactor"

	ethClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, channel.ContractBackend{}, err
	}

	return ethClient, channel.NewContractBackend(ethClient, transactor, txFinalityDepth), nil
}

type dummyAccount struct {
	addr wire.Address
}

// Address used by this account.
func (a dummyAccount) Address() wallet.Address {
	return a.addr
}

// SignData requests a signature from this account.
// It returns the signature or an error.
func (a dummyAccount) SignData(data []byte) ([]byte, error) {
	panic("unsupported")
}
