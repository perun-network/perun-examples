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
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"perun.network/go-perun/backend/ethereum/channel"
	"perun.network/go-perun/backend/ethereum/wallet"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wire"
	"perun.network/go-perun/wire/net"
	"perun.network/go-perun/wire/net/simple"
)

type PeerWithAddress struct {
	Peer    wire.Address
	Address string
}

type PerunClientConfig struct {
	PrivateKey      *ecdsa.PrivateKey
	Host            string
	ETHNodeURL      string
	AdjudicatorAddr common.Address
	AssetHolderAddr common.Address
	DialerTimeout   time.Duration
	PeerAddresses   []PeerWithAddress
}

type PerunClient struct {
	EthClient       *ethclient.Client
	StateChClient   *client.Client
	Bus             *net.Bus
	Listener        net.Listener
	ContractBackend channel.ContractInterface
	Wallet          *swallet.Wallet
	Account         *swallet.Account
}

func setupPerunClient(cfg PerunClientConfig) (*PerunClient, error) {
	// Create wallet and account
	clientWallet := swallet.NewWallet(cfg.PrivateKey)
	addr := wallet.AsWalletAddr(crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey))
	pAccount, err := clientWallet.Unlock(addr)
	if err != nil {
		panic("failed to create account")
	}
	account := pAccount.(*swallet.Account)

	// Create Ethereum client and contract backend
	signer := types.NewEIP155Signer(big.NewInt(1337))
	transactor := swallet.NewTransactor(clientWallet, signer)

	ethClient, cb, err := createContractBackend(cfg.ETHNodeURL, transactor)
	if err != nil {
		return nil, errors.WithMessage(err, "creating contract backend")
	}

	//REMARK
	// A client should validate the smart contracts it is using them, if the
	// given contract addresses come from an untrusted source.

	adjudicator := channel.NewAdjudicator(cb, cfg.AdjudicatorAddr, account.Account.Address, account.Account)

	listener, bus, err := setupNetwork(account, cfg.Host, cfg.PeerAddresses, cfg.DialerTimeout)
	if err != nil {
		return nil, errors.WithMessage(err, "setting up network")
	}

	funder := createFunder(cb, account.Account, cfg.AssetHolderAddr)

	stateChClient, err := client.New(account.Address(), bus, funder, adjudicator, clientWallet)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	return &PerunClient{ethClient, stateChClient, bus, listener, cb, clientWallet, account}, nil
}

func createContractBackend(nodeURL string, transactor channel.Transactor) (*ethclient.Client, channel.ContractBackend, error) {
	ethClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, channel.ContractBackend{}, nil
	}

	return ethClient, channel.NewContractBackend(ethClient, transactor), nil
}

func setupNetwork(account wire.Account, host string, peerAddresses []PeerWithAddress, dialerTimeout time.Duration) (listener net.Listener, bus *net.Bus, err error) {
	dialer := simple.NewTCPDialer(dialerTimeout)

	for _, pa := range peerAddresses {
		dialer.Register(pa.Peer, pa.Address)
	}

	listener, err = simple.NewTCPListener(host)
	if err != nil {
		err = fmt.Errorf("creating listener: %w", err)
		return
	}

	bus = net.NewBus(account, dialer)
	return listener, bus, nil
}

func createFunder(cb channel.ContractBackend, account accounts.Account, assetHolder common.Address) *channel.Funder {
	f := channel.NewFunder(cb)
	asset := wallet.Address(assetHolder)
	depositor := new(channel.ETHDepositor)
	f.RegisterAsset(asset, depositor, account)
	return f
}
