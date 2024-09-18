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
	"fmt"
	"log"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/wire"
	"perun.network/perun-icp-backend/client"
	"perun.network/perun-icp-backend/wallet"
)

const (
	Host              = "http://127.0.0.1"
	Port              = 4943
	perunPrincipal    = "be2us-64aaa-aaaaa-qaabq-cai"
	ledgerPrincipal   = "bkyz2-fmaaa-aaaaa-qaaaq-cai"
	userAId           = "97520b79b03e38d3f6b38ce5026a813ccc9d1a3e830edb6df5970e6ca6ad84be"
	userBId           = "40fd2dc85bc7d264b31f1fa24081d7733d303b49b7df84e3d372338f460aa678"
	userAPemPath      = "./userdata/identities/usera_identity.pem"
	userBPemPath      = "./userdata/identities/userb_identity.pem"
	channelCollateral = 50000
)

func main() {

	log.Println("Setting up wallets for Alice and Bob")
	perunWalletAlice := wallet.NewWallet()
	perunWalletBob := wallet.NewWallet()

	log.Println("Create communication channel between Alice and Bob")
	bus := wire.NewLocalBus()

	log.Println("Setting up Payment Clients")
	alice, err := client.SetupPaymentClient("Alice", perunWalletAlice, bus, perunPrincipal, ledgerPrincipal, Host, Port, userAPemPath)
	if err != nil {
		panic(err)
	}

	bob, err := client.SetupPaymentClient("Bob", perunWalletBob, bus, perunPrincipal, ledgerPrincipal, Host, Port, userBPemPath)
	if err != nil {
		panic(err)
	}

	log.Println("Alice opens Channel with Bob")
	alice.OpenChannel(bob.WireAddress(), channelCollateral)
	achan := alice.Channel
	log.Println("Alice opened channel")
	bob.AcceptedChannel()
	log.Println("Bob accepts Channel from Alice")
	bchan := bob.Channel

	log.Println("Initial Balances in the Channel")
	printBalances(achan.GetChannelState().Balances)

	// sending payment/s

	log.Println("Sending payments...")
	achan.SendPayment(1000)

	log.Println("Balance after first Payment:")
	printBalances(achan.GetChannelState().Balances)

	bchan.SendPayment(2000)

	log.Println("Balance after second Payment:")
	printBalances(achan.GetChannelState().Balances)

	log.Println("Settle channel")
	achan.Settle()

	log.Println("Final Balance after settling:")
	printBalances(achan.GetChannelState().Balances)

	log.Println("Shutdown Channel")
	alice.Shutdown()
	bob.Shutdown()
	log.Println("Done")
}

func printBalances(balances channel.Balances) {
	log.Println("Balances:")
	for i, assetBalances := range balances {
		fmt.Printf("Asset %d:\n", i+1)
		fmt.Printf("  Alice: %s\n", assetBalances[0].String()) // Access balance for Alice
		fmt.Printf("  Bob: %s\n", assetBalances[1].String())   // Access balance for Bob
	}
}
