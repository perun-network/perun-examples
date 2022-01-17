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
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"perun.network/go-perun/backend/ethereum/channel"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
	"perun.network/perun-examples/payment-channel/client"
)

func printAccountBalance(clients ...*client.Client) {
	for _, c := range clients {
		globalBalance, err := c.OnChainBalance()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s with address %v - Account Balance: %v", c.Name(), c.Address(), toEth(globalBalance))
	}
}

func toEth(weiAmount *big.Int) string {
	return fmt.Sprintf("%vETH", eth.WeiToEth(weiAmount))
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
	adj, err = channel.DeployAdjudicator(context.TODO(), cb, acc)
	if err != nil {
		panic(err)
	}

	// Deploy asset holder.
	ah, err = channel.DeployETHAssetholder(context.TODO(), cb, adj, acc)
	if err != nil {
		panic(err)
	}

	return adj, ah
}
