// Copyright 2024 PolyCrypt GmbH
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

package ethereumUtil

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	"log"
	"math/big"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wire"
	"perun.network/perun-examples/payment-channel-cc/client"
	stellarChannel "perun.network/perun-stellar-backend/channel"
	stellarWallet "perun.network/perun-stellar-backend/wallet"
)

// DeployContracts deploys the Perun smart contracts on the specified ledger.
func DeployContracts(nodeURL string, chainID uint64, privateKey string) (adj, ah common.Address) {
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
	adj, err = ethchannel.DeployAdjudicator(context.TODO(), cb, acc)
	if err != nil {
		panic(err)
	}

	// Deploy asset holder.
	ah, err = ethchannel.DeployETHAssetholder(context.TODO(), cb, adj, acc)
	if err != nil {
		panic(err)
	}

	return adj, ah
}

// SetupPaymentClient sets up a new client with the given parameters.
func SetupPaymentClient(
	bus wire.Bus,
	nodeURL string,
	adjudicator common.Address,
	asset ethwallet.Address,
	k *ecdsa.PrivateKey,
	stellarWallet *stellarWallet.EphemeralWallet, // stellarWallet is the wallet used for signing stellar transactions.
	stellarAccount *stellarWallet.Account, // stellarAccount is the account to be used for signing Stellar transactions.
	stellarTokenIDs channel.Asset, // stellarTokenIDs is the list of token IDs to be used for payment channels.
	stellarFunder *stellarChannel.Funder, // stellarFunder is the funder to be used for funding Stellar payment channels.
	stellarAdj *stellarChannel.Adjudicator, // stellarAdj is the adjudicator to be used for Stellar payment channels.
) *client.PaymentClient {
	// Create wallet and account.
	w := swallet.NewWallet(k)
	acc := crypto.PubkeyToAddress(k.PublicKey)
	eaddr := ethwallet.AsWalletAddr(acc)

	// Create and start client.
	c, err := client.SetupPaymentClient(
		bus,
		w,
		acc,
		eaddr,
		nodeURL,
		1337,
		adjudicator,
		asset,
		stellarWallet,
		stellarAccount,
		stellarTokenIDs,
		stellarFunder,
		stellarAdj,
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

// NewBalanceLogger creates a new balance logger for the specified ledger.
func NewBalanceLogger(chainURL string) balanceLogger {
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
