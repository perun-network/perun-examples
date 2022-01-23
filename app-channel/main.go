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

//todo:tutorial Mention that we use context.TODO and panic(err) to keep the code in simple, but in production code one should always use proper context and handle error appropriately.

// main runs a demo of the game client. It assumes that a blockchain node is
// available at `chainURL` and that the accounts corresponding to the specified
// secret keys are provided with sufficient funds.
func main() {
	// Deploy contracts.
	log.Println("Deploying contracts.")
	adjudicator, assetHolder, tikTakToeApp := deployContracts(chainURL, chainID, keyDeployer)
	asset := client.NewAsset(assetHolder)

	// Setup clients.
	log.Println("Setting up clients.")
	bus := wire.NewLocalBus() // Message bus used for off-chain communication.	//TODO:tutorial Extension that explains tcp/ip bus.
	alice := setupGameClient(bus, chainURL, adjudicator, assetHolder, keyAlice)
	bob := setupGameClient(bus, chainURL, adjudicator, assetHolder, keyBob)

	// Print balances before transactions.
	l := newBalanceLogger(chainURL)
	l.LogBalances(alice, bob)

	// Open app channel, play, close.
	log.Println("Opening channel.")
	_alice, chID := alice.ProposeGame(bob, asset, tikTakToeApp, 10)
	_bob := bob.GetGame(chID)

	// Alice set (2, 0)
	_alice.Set(2, 0)

	// Bob set (0, 0)
	_bob.Set(0, 0)

	// Alice set (0, 2)
	_alice.Set(0, 2)

	// Bob set (1, 1)
	_bob.Set(1, 1)

	// Alice set (2, 2)
	_alice.Set(2, 2)

	// Bob set (2, 1)
	_bob.Set(2, 1)

	// Alice set (1, 2)
	_alice.Set(1, 2)

	// Alice concludes
	_alice.Settle()

	alice.Shutdown()
	bob.Shutdown()
}
