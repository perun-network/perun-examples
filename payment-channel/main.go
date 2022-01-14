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
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
	"perun.network/perun-examples/payment-channel/client"
)

const (
	chainURL              = "ws://127.0.0.1:8545"
	chainID               = 1337
	defaultContextTimeout = 15 * time.Second
	deploymentKey         = "" // Key used for contract deployment. //TODO insert key

	// Alice
	hostAlice = "0.0.0.0:8401"
	keyAlice  = "1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f"

	// Bob
	hostBob = "0.0.0.0:8402"
	keyBob  = "f63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e"
)

func main() {
	contracts := deployContracts(chainURL, chainID, deploymentKey)

	alice := startClient(chainURL, contracts, hostAlice, keyAlice)
	bob := startClient(chainURL, contracts, hostBob, keyBob)

	logAccountBalance(alice, bob)

	//TODO validate asset holder when proposing a new channel or receiving a channel proposal.

	alice.OpenChannel()
	alice.UpdateChannel()
	bob.UpdateChannel()
	bob.CloseChannel()

	logAccountBalance(alice, bob)
}

// startClient sets up a new client with the given parameters.
func startClient(
	nodeURL string,
	contracts ContractAddresses,
	host string,
	privateKey string,
) *client.Client {
	// Create wallet and account.
	k, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}
	w := swallet.NewWallet(k)
	acc := crypto.PubkeyToAddress(k.PublicKey)

	// Create and start client.
	c, err := client.StartClient(
		host,
		w,
		acc,
		nodeURL,
		chainID,
		contracts.Adjudicator,
	)
	if err != nil {
		panic(err)
	}

	return c
}
