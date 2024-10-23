// Copyright 2024 PolyCrypt GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

	"perun.network/perun-examples/payment-channel-xlm/client"
	"perun.network/perun-examples/payment-channel-xlm/util"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/wire"
)

func main() {
	// Initialize setup (replaces the test setup)
	log.Println("Starting initial setup")
	setup, err := util.NewExampleSetup()
	if err != nil {
		panic(err)
	}

	payment_example(setup)
}

func payment_example(setup *util.Setup) {
	log.Println("Creating Accounts for Alice and Bob")
	accAlice := setup.GetAccounts()[0]
	accBob := setup.GetAccounts()[1]
	wAlice := setup.GetWallets()[0]
	wBob := setup.GetWallets()[1]
	funderAlice := setup.GetFunders()[0]
	funderBob := setup.GetFunders()[1]
	adjAlice := setup.GetAdjudicators()[0]
	adjBob := setup.GetAdjudicators()[1]

	log.Println("Initializing a connection between Alice and Bob")
	bus := wire.NewLocalBus()
	log.Println("Setup payment clients for Alice and Bob")
	alicePerun, err := client.SetupPaymentClient(wAlice, accAlice, setup.GetTokenAsset(), bus, funderAlice, adjAlice)
	if err != nil {
		panic(err)
	}
	bobPerun, err := client.SetupPaymentClient(wBob, accBob, setup.GetTokenAsset(), bus, funderBob, adjBob)
	if err != nil {
		panic(err)
	}

	log.Println("Setting initial balances")
	balances := channel.Balances{
		{big.NewInt(1000), big.NewInt(100)},
		{big.NewInt(0), big.NewInt(1000)},
	}

	log.Println("Alice opens a channel with Bob")
	alicePerun.OpenChannel(bobPerun.WireAddress(), balances)
	aliceChannel := alicePerun.Channel
	bobChannel := bobPerun.AcceptedChannel()

	printBalances(alicePerun.Channel.GetChannelState().Balances)

	log.Println("Alice sends payment to Bob")

	aliceChannel.SendPayment(500, 1)
	printBalances(alicePerun.Channel.GetChannelState().Balances)

	log.Println("Bob sends payment to Alice")

	bobChannel.SendPayment(250, 1)
	printBalances(alicePerun.Channel.GetChannelState().Balances)

	log.Println("Channel is being settled")
	aliceChannel.Settle()
	bobChannel.Settle()

	alicePerun.Shutdown()
	bobPerun.Shutdown()

	log.Println("Done")
}

func printBalances(balances channel.Balances) {
	log.Println("Channel Balances:")

	// Manually print for Asset 1
	log.Printf("Asset 1:\n")
	log.Printf("  Alice: %s\n", balances[0][0].String())
	log.Printf("  Bob: %s\n", balances[0][1].String())

	// Manually print for Asset 2
	log.Printf("Asset 2:\n")
	log.Printf("  Alice: %s\n", balances[1][0].String())
	log.Printf("  Bob: %s\n", balances[1][1].String())
}
