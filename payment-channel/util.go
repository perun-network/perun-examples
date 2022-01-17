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
	"context"
	"log"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	"perun.network/go-perun/backend/ethereum/wallet"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
	"perun.network/perun-examples/payment-channel/client"
)

type balanceLogger struct {
	ethClient *ethclient.Client
}

func newBalanceLogger(chainURL string) balanceLogger {
	c, err := ethclient.Dial(chainURL)
	if err != nil {
		panic(err)
	}
	return balanceLogger{ethClient: c}
}

func (l balanceLogger) LogBalances(clients ...*client.Client) {
	for _, c := range clients {
		addr := common.Address(*c.AccountAddress.(*wallet.Address))
		bal, err := l.ethClient.BalanceAt(context.TODO(), addr, nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Balance of %v: %v", c.Name, bal)
	}
}

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
	adj, err = ethchannel.DeployAdjudicator(context.TODO(), cb, acc) //TODO accept ethwallet Account instead?
	if err != nil {
		panic(err)
	}

	// Deploy asset holder.
	ah, err = ethchannel.DeployETHAssetholder(context.TODO(), cb, adj, acc) //TODO accept ethwallet Account instead?
	if err != nil {
		panic(err)
	}

	return adj, ah
}
