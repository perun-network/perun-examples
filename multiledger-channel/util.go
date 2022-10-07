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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/perun-network/perun-eth-backend/bindings/peruntoken"

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
func deployContracts(chains []client.ChainConfig) {
	k, err := crypto.HexToECDSA(keyDeployer)
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

		// Deploy the ERC20 PerunToken.
		initAccs := []common.Address{privKeyToAddress(keyAlice), privKeyToAddress(keyBob)}
		initAmount := big.NewInt(initialTokenAmount)
		token, err := ethchannel.DeployPerunToken(context.TODO(), cb, acc, initAccs, initAmount)
		chains[i].Token = token
		if err != nil {
			panic(err)
		}

		// Deploy adjudicator.
		chains[i].Adjudicator, err = ethchannel.DeployAdjudicator(context.TODO(), cb, acc)
		if err != nil {
			panic(err)
		}

		// Deploy ERC20 asset holder.
		chains[i].AssetHolder, err = ethchannel.DeployERC20Assetholder(context.TODO(), cb, chains[i].Adjudicator, token, acc)
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

func privKeyToAddress(privateKey string) common.Address {
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}

	return crypto.PubkeyToAddress(pk.PublicKey)
}

// balanceLogger is a utility for logging client balances on two chains (A & B).
type balanceLogger struct {
	ethClientA *ethclient.Client
	ethClientB *ethclient.Client
	tokenA     common.Address
	tokenB     common.Address
}

// newBalanceLogger creates a new balance logger for the specified ledger.
func newBalanceLogger(chainA, chainB client.ChainConfig) balanceLogger {
	clientA, err := ethclient.Dial(chainA.ChainURL)
	if err != nil {
		panic(err)
	}

	clientB, err := ethclient.Dial(chainB.ChainURL)
	if err != nil {
		panic(err)
	}
	return balanceLogger{
		ethClientA: clientA,
		ethClientB: clientB,
		tokenA:     chainA.Token,
		tokenB:     chainA.Token,
	}
}

// LogBalances prints the token balances of the specified clients on the two chains.
func (l balanceLogger) LogBalances(addresses ...common.Address) {
	getBals := func(cb bind.ContractBackend, token common.Address) []*big.Int {
		bals := make([]*big.Int, len(addresses))

		t, err := peruntoken.NewPeruntoken(token, cb)
		if err != nil {
			panic(err)
		}

		for i, a := range addresses {
			bal, err := t.BalanceOf(&bind.CallOpts{}, a)
			if err != nil {
				log.Fatal(err)
			}
			bals[i] = bal
		}
		return bals
	}

	log.Println("Client balances Chain A (PRN):", getBals(l.ethClientA, l.tokenA))
	log.Println("Client balances Chain B (PRN):", getBals(l.ethClientB, l.tokenB))
}
