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
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"perun.network/go-perun/backend/ethereum/channel"
	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wire"
	"perun.network/go-perun/wire/net"
	"perun.network/go-perun/wire/net/simple"
)

type PeerWithAddress struct {
	Peer    wire.Address
	Address string
}

type ClientConfig struct {
	PrivateKey      *ecdsa.PrivateKey
	Host            string
	ETHNodeURL      string
	AdjudicatorAddr common.Address
	AssetHolderAddr common.Address
	DialerTimeout   time.Duration
	PeerAddresses   []PeerWithAddress
}

type Client struct {
	EthClient       *ethclient.Client
	PerunClient     *client.Client
	Bus             *net.Bus
	Listener        net.Listener
	ContractBackend channel.ContractInterface
	Wallet          *Wallet
	Account         *Account
}

func SetupClient(cfg ClientConfig) (*Client, error) {
	// Create wallet and account
	w := createWallet(cfg.PrivateKey)
	addr := wallet.AsWalletAddr(crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey))
	pAccount, err := w.Unlock(addr)
	if err != nil {
		panic("failed to create account")
	}
	account := pAccount.(*Account)

	// Create Ethereum client and contract backend
	ethClient, cb, err := createContractBackend(cfg.ETHNodeURL, w)
	if err != nil {
		return nil, errors.WithMessage(err, "creating contract backend")
	}

	//REMARK
	// A client should validate the smart contracts it is using them, if the
	// given contract addresses come from an untrusted source.

	adjudicator := channel.NewAdjudicator(cb, cfg.AdjudicatorAddr, account.EthAddress(), account.EthAccount())

	listener, bus, err := setupNetwork(account, cfg.Host, cfg.PeerAddresses, cfg.DialerTimeout)
	if err != nil {
		return nil, errors.WithMessage(err, "setting up network")
	}

	funder := createFunder(cb, account.EthAccount(), cfg.AssetHolderAddr)

	c, err := client.New(account.Address(), bus, funder, adjudicator, w)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	return &Client{ethClient, c, bus, listener, cb, w, account}, nil
}

func createContractBackend(nodeURL string, wallet *Wallet) (*ethclient.Client, channel.ContractBackend, error) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, channel.ContractBackend{}, nil
	}

	return client, channel.NewContractBackend(client, wallet), nil
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
