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
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"perun.network/go-perun/wallet"

	dotwallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	"log"
	"perun.network/go-perun/wire"
	"perun.network/perun-examples/payment-channel/client"
)

// setupPaymentClient sets up a new client with the given parameters.
func setupPaymentClient(
	bus wire.Bus,
	nodeURL string,
	networkId dot.NetworkID,
	queryDepth types.BlockNumber,
	privateKey string,
) *client.PaymentClient {
	// Create wallet and account.
	sk, err := sr25519.NewSKFromHex(privateKey)
	if err != nil {
		panic(err)
	}
	w := dotwallet.NewWallet()
	acc := w.ImportSK(sk)

	// Create and start client.
	c, err := client.SetupPaymentClient(
		bus,
		w,
		acc,
		nodeURL,
		networkId,
		queryDepth,
	)
	if err != nil {
		panic(err)
	}

	return c
}

// balanceLogger is a utility for logging client balances.
type balanceLogger struct {
	api *dot.API
}

// newBalanceLogger creates a new balance logger for the specified ledger.
func newBalanceLogger(nodeURL string, networkId dot.NetworkID) balanceLogger {
	api, err := dot.NewAPI(nodeURL, networkId)
	if err != nil {
		panic(err)
	}
	return balanceLogger{api: api}
}

// LogBalances prints the balances of the specified accounts.
func (l balanceLogger) LogBalances(addrs ...wallet.Address) {
	bals := make([]*dot.Dot, len(addrs))
	for i, addr := range addrs {
		accInfo, err := l.api.AccountInfo(dotwallet.AsAddr(addr).AccountID())
		if err != nil {
			panic(err)
		}
		bals[i] = client.PlanckToDot(accInfo.Free.Int)
	}
	log.Println("Client balances (DOT):", bals)
}
