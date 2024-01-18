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
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"perun.network/go-perun/wire"
	"perun.network/perun-collateralized-channels/app"
	"perun.network/perun-collateralized-channels/client"
	"perun.network/perun-collateralized-channels/eth"

	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
)

type Environment struct {
	clientNames map[common.Address]string
}

func (e *Environment) ClientName(clientAddr common.Address) string {
	name, ok := e.clientNames[clientAddr]
	if !ok {
		return fmt.Sprintf("%v", clientAddr)
	}
	return name
}

func (e *Environment) logAccountBalance(clients ...*client.AppClient) {
	for _, c := range clients {
		globalBalance, err := c.OnChainBalance()
		if err != nil {
			log.Fatal(err)
		}
		collateralBalance, err := c.PeerCollateral()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v: Account Balance: Global=%v, Locked as Collateral=%v", e.ClientName(c.WalletAddress()), toEth(globalBalance), toEth(collateralBalance))
	}
}

func (l *Environment) logChannelBalances(clients ...*client.AppClient) {
	for _, c := range clients {
		channelBalances, err := c.ChannelBalances()
		if err != nil {
			log.Panicf("getting channel balances: %v", err)
		}

		var b strings.Builder
		fmt.Fprintf(&b, "%v: ", l.ClientName(c.WalletAddress()))
		for peer, bal := range channelBalances {
			balCollateral, err := c.ChannelFunding(peer)
			if err != nil {
				log.Panicf("getting channel collateral balance: %v", err)
			}
			total := new(big.Int).Add(balCollateral, bal)
			fmt.Fprintf(&b, "Channel with %v: Balance: %v (OnChain=%v, OffChain=%v)\n", l.ClientName(peer), toEth(total), toEth(balCollateral), toEth(bal))
		}
		log.Print(b.String())
	}
}

func toEth(weiAmount *big.Int) string {
	return fmt.Sprintf("%vETH", eth.WeiToEth(weiAmount))
}

func deployContracts(nodeURL string, deploymentKey *ecdsa.PrivateKey, contextTimeout time.Duration, collateralWithdrawalDelay time.Duration) (contracts ContractAddresses, err error) {
	ethClient, err := eth.NewEthClient(nodeURL, deploymentKey, contextTimeout)
	if err != nil {
		err = errors.WithMessage(err, "creating ethereum client")
		return
	}

	// Deploy adjudicator
	adjudicatorAddr, txAdj, err := ethClient.DeployAdjudicator()
	if err != nil {
		err = errors.WithMessage(err, "deploying adjudicator")
		return
	}

	// Deploy collateralized channels app
	appAddr, txApp, err := ethClient.DeployCollateralApp(adjudicatorAddr)
	if err != nil {
		err = errors.WithMessage(err, "deploying CollateralApp")
		return
	}

	// Deploy asset holder
	assetHolderAddr, txAss, err := ethClient.DeployCollateralAssetHolderETH(adjudicatorAddr, appAddr, big.NewInt(int64(collateralWithdrawalDelay.Seconds())))
	if err != nil {
		err = errors.WithMessage(err, "deploying CollateralAssetHolderETH")
		return
	}

	err = ethClient.WaitDeployment(txAdj, txApp, txAss)
	if err != nil {
		err = errors.WithMessage(err, "waiting for contract deployment")
		return
	}

	return ContractAddresses{
		AdjudicatorAddr: adjudicatorAddr,
		AssetHolderAddr: assetHolderAddr,
		AppAddr:         appAddr,
	}, nil
}

type ContractAddresses struct {
	AdjudicatorAddr, AssetHolderAddr, AppAddr common.Address
}

func setupClient(
	bus wire.Bus,
	nodeURL string,
	chainID uint64,
	assetHolderAddr common.Address,
	privateKey *ecdsa.PrivateKey,
	app *app.CollateralApp,
	updatePolicy client.PaymentAcceptancePolicy,
	challengeDuration time.Duration,
	contextTimeout time.Duration,

) (*client.AppClient, error) {
	w := swallet.NewWallet(privateKey)
	acc := crypto.PubkeyToAddress(privateKey.PublicKey)
	eaddr := ethwallet.AsWalletAddr(acc)

	// Create and start client.
	return client.SetupAppClient(
		bus,
		w,
		acc,
		assetHolderAddr,
		eaddr,
		nodeURL,
		chainID,
		app,
		updatePolicy,
		challengeDuration,
		contextTimeout,
	)
}
