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

	"perun.network/perun-examples/multiledger-channel/client"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"

	"perun.network/go-perun/wire"
)

// deployContracts deploys the Perun smart contracts on the specified chains and
// sets the contract addresses.
func deployContracts(chains []client.ChainConfig, privateKey string) {
	k, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}
	w := swallet.NewWallet(k)

	for i, chain := range chains {
		cb, err := client.CreateContractBackend(chain.ChainURL, chain.ChainID.Int, w)
		if err != nil {
			panic(err)
		}
		acc := accounts.Account{Address: crypto.PubkeyToAddress(k.PublicKey)}

		// Deploy adjudicator.
		chains[i].Adjudicator, err = ethchannel.DeployAdjudicator(context.TODO(), cb, acc)
		if err != nil {
			panic(err)
		}

		// Deploy asset holder.
		chains[i].AssetHolder, err = ethchannel.DeployETHAssetholder(context.TODO(), cb, chains[i].Adjudicator, acc)
		if err != nil {
			panic(err)
		}
	}

}

// setupPaymentClient sets up a new client with the given parameters.
func setupPaymentClient(
	bus wire.Bus,
	privateKey string,
	chains [2]client.ChainConfig,
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
		chains,
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

// LogBalances prints the balances of the specified accounts.
func (l balanceLogger) LogBalances(accounts ...common.Address) {
	bals := make([]*big.Float, len(accounts))
	for i, c := range accounts {
		bal, err := l.ethClient.BalanceAt(context.TODO(), c, nil)
		if err != nil {
			log.Fatal(err)
		}
		bals[i] = client.WeiToEth(bal)
	}
	log.Println("Client balances (ETH):", bals)
}
