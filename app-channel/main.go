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

	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/wire"
	"perun.network/perun-examples/app-channel/app"
	"perun.network/perun-examples/app-channel/client"
)

const (
	chainURL = "ws://127.0.0.1:8545"
	chainID  = 1337

	// Private keys.
	keyDeployer = "79ea8f62d97bc0591a4224c1725fca6b00de5b2cea286fe2e0bb35c5e76be46e"
	keyAlice    = "1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f"
	keyBob      = "f63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e"
)

//todo:tutorial Mention that we use context.TODO and panic(err) to keep the code simple, but in production code one should always use proper context and handle error appropriately.

// main runs a demo of the game client. It assumes that a blockchain node is
// available at `chainURL` and that the accounts corresponding to the specified
// secret keys are provided with sufficient funds.
func main() {
	// Deploy contracts.
	log.Println("Deploying contracts.")
	adjudicator, assetHolder, appAddress := deployContracts(chainURL, chainID, keyDeployer)
	asset := client.NewAsset(assetHolder)
	app := app.NewTicTacToeApp(wallet.AsWalletAddr(appAddress))

	// Setup clients.
	log.Println("Setting up clients.")
	bus := wire.NewLocalBus() // Message bus used for off-chain communication.	//TODO:tutorial Extension that explains tcp/ip bus.
	stake := big.NewInt(10)
	alice := setupGameClient(bus, chainURL, adjudicator, assetHolder, keyAlice, app, stake)
	bob := setupGameClient(bus, chainURL, adjudicator, assetHolder, keyBob, app, stake)

	// Print balances before transactions.
	l := newBalanceLogger(chainURL)
	l.LogBalances(alice, bob)

	// Open app channel, play, close.
	log.Println("Opening channel.")
	appAlice := alice.ProposeAppChannel(bob, asset)
	appBob := bob.AcceptedGame()

	log.Println("Start playing.")
	log.Println("Alice's turn.")
	// Alice set (2, 0)
	appAlice.Set(2, 0)

	log.Println("Bob's turn.")
	// Bob set (0, 0)
	appBob.Set(0, 0)

	log.Println("Alice's turn.") //TODO:question Can we make the turn's appear "less instant" in the console output
	// Alice set (0, 2)
	appAlice.Set(0, 2)

	log.Println("Bob's turn.")
	// Bob set (1, 1)
	appBob.Set(1, 1)

	log.Println("Alice's turn.")
	// Alice set (2, 2)
	appAlice.Set(2, 2)

	log.Println("Bob's turn.")
	// Bob set (2, 1)
	appBob.Set(2, 1)

	log.Println("Alice's turn.")
	// Alice set (1, 2)
	appAlice.Set(1, 0)

	log.Println("Bob's turn.")
	// Bob set (2, 1)
	appBob.Set(0, 1)

	log.Println("Bob wins.")
	log.Println("Payout.")

	//TODO Include dispute? Or separate tutorial?

	// Payout.
	appAlice.Settle()
	appBob.Settle()

	// Print balances after transactions.
	l.LogBalances(alice, bob)

	// Cleanup.
	alice.Shutdown()
	bob.Shutdown()
}
