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

	"perun.network/go-perun/wire"
)

const (
	chainURL        = "ws://127.0.0.1:9944"
	networkID       = 42
	blockQueryDepth = 10

	// Private keys.
	keyAlice = "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a"
	keyBob   = "0x398f0c28f98885e046333d4a41c19cee4c37368a9832c6502f6cfd182e2aef89"
)

// main runs a demo of the payment client. It assumes that a blockchain node is
// available at `chainURL` and that the accounts corresponding to the specified
// secret keys are provided with sufficient funds.
func main() {
	// Setup clients.
	log.Println("Setting up clients.")
	bus := wire.NewLocalBus() // Message bus used for off-chain communication.
	alice := setupPaymentClient(bus, chainURL, networkID, blockQueryDepth, keyAlice)
	bob := setupPaymentClient(bus, chainURL, networkID, blockQueryDepth, keyBob)

	// Print balances before transactions.
	l := newBalanceLogger(chainURL, networkID)
	l.LogBalances(alice.WalletAddress(), bob.WalletAddress())

	// Open channel, transact, close.
	log.Println("Opening channel and depositing funds.")
	chAlice := alice.OpenChannel(bob.WireAddress(), 100000)
	chBob := bob.AcceptedChannel()

	log.Println("Sending payments...")
	chAlice.SendPayment(50000)
	chBob.SendPayment(25000)
	chAlice.SendPayment(25000)

	log.Println("Settling channel.")
	chAlice.Settle(false) // Conclude and withdraw.
	chBob.Settle(false)   // Withdraw.

	// Print balances after transactions.
	l.LogBalances(alice.WalletAddress(), bob.WalletAddress())

	// Cleanup.
	alice.Shutdown()
	bob.Shutdown()

	log.Println("Done.")
}
