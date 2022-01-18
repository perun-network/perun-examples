// Copyright 2022 PolyCrypt GmbH
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
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
	"perun.network/go-perun/wire"
	"perun.network/perun-examples/payment-channel/client"
)

// deployContracts deploys the Perun smart contracts on the specified ledger.
func deployContracts(nodeURL string, chainID uint64, privateKey string) (adj, ah common.Address) {
	k, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}
	w := swallet.NewWallet(k)
	cb, err := client.CreateContractBackend(nodeURL, chainID, w)
	if err != nil {
		panic(err)
	}
	acc := accounts.Account{Address: crypto.PubkeyToAddress(k.PublicKey)}

	// Deploy adjudicator.
	adj, err = ethchannel.DeployAdjudicator(context.TODO(), cb, acc) //TODO:go-perun accept ethwallet Account instead?
	if err != nil {
		panic(err)
	}

	// Deploy asset holder.
	ah, err = ethchannel.DeployETHAssetholder(context.TODO(), cb, adj, acc) //TODO:go-perun accept ethwallet Account instead?
	if err != nil {
		panic(err)
	}

	return adj, ah
}

// setupPaymentClient sets up a new client with the given parameters.
func setupPaymentClient(
	bus wire.Bus,
	nodeURL string,
	adjudicator, assetHolder common.Address,
	privateKey string,
) *client.PaymentClient {
	// Create wallet and account.
	k, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}
	w := swallet.NewWallet(k)
	acc := crypto.PubkeyToAddress(k.PublicKey)

	// Create and start client.
	c, err := client.SetupPaymentClient(
		bus,
		w,
		acc,
		nodeURL,
		chainID,
		adjudicator,
		assetHolder,
	)
	if err != nil {
		panic(err)
	}

	return c
}

// balanceLogger is a utility for logging client balances.
type balanceLogger struct {
	ethClient *ethclient.Client
}

// newBalanceLogger creates a new balance logger for the specified ledger.
func newBalanceLogger(chainURL string) balanceLogger {
	c, err := ethclient.Dial(chainURL)
	if err != nil {
		panic(err)
	}
	return balanceLogger{ethClient: c}
}

// LogBalances prints the balances of the specified clients.
func (l balanceLogger) LogBalances(clients ...*client.PaymentClient) {
	bals := make([]*big.Int, len(clients))
	for i, c := range clients {
		bal, err := l.ethClient.BalanceAt(context.TODO(), c.AccountAddress(), nil)
		if err != nil {
			log.Fatal(err)
		}
		bals[i] = bal
	}
	log.Println("Client balances:", bals)
}
