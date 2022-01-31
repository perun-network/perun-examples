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
	"log"

	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/wire"
)

const (
	chainURL = "ws://127.0.0.1:8545"
	chainID  = 1337

	// Private keys.
	keyDeployer = "79ea8f62d97bc0591a4224c1725fca6b00de5b2cea286fe2e0bb35c5e76be46e"
	keyAlice    = "1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f"
	keyBob      = "f63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e"
)

//todo:tutorial Mention that we use context.TODO and panic(err) to keep the code in simple, but in production code one should always use proper context and handle error appropriately.

// main runs a demo of the payment client. It assumes that a blockchain node is
// available at `chainURL` and that the accounts corresponding to the specified
// secret keys are provided with sufficient funds.
func main() {
	// Deploy contracts.
	log.Println("Deploying contracts.")
	adjudicator, assetHolder := deployContracts(chainURL, chainID, keyDeployer)
	asset := ethwallet.AsWalletAddr(assetHolder)

	// Setup clients.
	log.Println("Setting up clients.")
	bus := wire.NewLocalBus() // Message bus used for off-chain communication.	//TODO:tutorial Extension that explains tcp/ip bus.
	alice := setupPaymentClient(bus, chainURL, adjudicator, *asset, keyAlice)
	bob := setupPaymentClient(bus, chainURL, adjudicator, *asset, keyBob)

	// Print balances before transactions.
	l := newBalanceLogger(chainURL)
	l.LogBalances(alice, bob)

	// Open channel, transact, close.
	log.Println("Opening channel.")
	chAlice := alice.OpenChannel(bob, 10)
	chBob := bob.AcceptedChannel()

	log.Println("Sending payments.")
	chAlice.SendPayment(1)
	chBob.SendPayment(2)
	chAlice.SendPayment(3)

	log.Println("Settling channel.")
	chAlice.Settle() // Conclude and withdraw.
	chBob.Settle()   // Withdraw.

	// Print balances after transactions.
	l.LogBalances(alice, bob)

	// Cleanup.
	alice.Shutdown()
	bob.Shutdown()
}
