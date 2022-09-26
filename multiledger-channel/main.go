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
	"math/big"

	echannel "github.com/perun-network/perun-eth-backend/channel"

	"perun.network/go-perun/wire"
	"perun.network/perun-examples/multiledger-channel/client"
)

const (
	// The URLs and IDs of the two chains.
	chainAURL = "ws://127.0.0.1:8545"
	chainAID  = 1337
	chainBURL = "ws://127.0.0.1:8546"
	chainBID  = 1338

	// Private keys.
	keyDeployer = "79ea8f62d97bc0591a4224c1725fca6b00de5b2cea286fe2e0bb35c5e76be46e"
	keyAlice    = "1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f"
	keyBob      = "f63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e"

	initialTokenAmount = 100
)

// main runs a demo of a multi-ledger channel. It assumes that two blockchain
// nodes are available at `chainAURL` and `chainBURL` and that the accounts
// corresponding to the specified secret keys are provided with sufficient funds.
func main() {
	// Deploy contracts.
	log.Println("Deploying contracts.")

	chainA := client.ChainConfig{
		ChainID:  echannel.MakeChainID(big.NewInt(chainAID)),
		ChainURL: chainAURL,
	}
	chainB := client.ChainConfig{
		ChainID:  echannel.MakeChainID(big.NewInt(chainBID)),
		ChainURL: chainBURL,
	}

	chains := [2]client.ChainConfig{chainA, chainB}

	// deployContracts will set the contract addresses of the chain
	// configurations after it deployed them.
	deployContracts(chains[:])

	// Setup clients.
	log.Println("Setting up clients.")
	bus := wire.NewLocalBus() // Message bus used for off-chain communication.
	alice := setupPaymentClient(bus, keyAlice, chains)
	bob := setupPaymentClient(bus, keyBob, chains)

	// Print balances before transactions.
	l := newBalanceLogger(chains[0], chains[1])
	l.LogBalances(alice.WalletAddress(), bob.WalletAddress())

	// Open channel, transact, close.
	log.Println("Opening channel and depositing funds.")

	// The balances that each party puts in the channel from the specific chain.
	chainABalances := [2]float64{15, 0} // Alice puts 5 PRN from chain A in the channel.
	chainBBalances := [2]float64{0, 42} // Bob puts 13.37 PRN from chain B in the channel.

	chAlice := alice.OpenChannel(bob.WireAddress(), chainABalances, chainBBalances)
	chBob := bob.AcceptedChannel()

	log.Println("Sending payments...")
	chAlice.SendPayment(15, client.ChainAIdx)
	chBob.SendPayment(42, client.ChainBIdx)

	log.Println("Settling channel.")
	chAlice.Settle() // Conclude and withdraw.
	chBob.Settle()   // Withdraw.

	// Print balances after transactions.
	l.LogBalances(alice.WalletAddress(), bob.WalletAddress())

	// Cleanup.
	alice.Shutdown()
	bob.Shutdown()
}
