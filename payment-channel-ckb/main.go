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
	"math/rand"
	"time"

	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
	"github.com/perun-network/perun-libp2p-wire/p2p"
	"perun.network/go-perun/channel"
	ckbchannel "perun.network/perun-ckb-backend/channel"
	"perun.network/perun-ckb-backend/channel/asset"
	"perun.network/perun-examples/payment-channel-ckb/client"
)

const (
	rpcNodeURL = "http://localhost:8114"
	Network    = types.NetworkTest // Network to use for the CKB client
)

func main() {
	//Setup devnet environment
	setup := NewSetup()

	log.Println("Setting up payment channel clients")
	aliceWireAcc := p2p.NewRandomAccount(rand.New(rand.NewSource(time.Now().UnixNano())))
	aliceNet, err := p2p.NewP2PBus(ckbchannel.CKBBackendID, aliceWireAcc)
	if err != nil {
		log.Fatalf("creating p2p net: %v", err)
	}
	aliceBus := aliceNet.Bus
	aliceListener := aliceNet.Listener
	go aliceBus.Listen(aliceListener)

	alice, err := client.NewPaymentClient(
		"Alice",
		Network,
		setup.Deployment,
		rpcNodeURL,
		setup.WalletAccs[0],
		*setup.AccKeys[0],
		setup.Wallets[0],
		aliceWireAcc.Address(),
		aliceNet,
	)
	if err != nil {
		log.Fatalf("error creating alice's client: %v", err)
	}

	bobWireAcc := p2p.NewRandomAccount(rand.New(rand.NewSource(time.Now().UnixNano())))
	bobNet, err := p2p.NewP2PBus(ckbchannel.CKBBackendID, bobWireAcc)
	if err != nil {
		log.Fatalf("creating p2p net: %v", err)
	}
	bobBus := bobNet.Bus
	bobListener := bobNet.Listener
	go bobBus.Listen(bobListener)

	bob, err := client.NewPaymentClient(
		"Bob",
		Network,
		setup.Deployment,
		rpcNodeURL,
		setup.WalletAccs[1],
		*setup.AccKeys[1],
		setup.Wallets[1],
		bobWireAcc.Address(),
		bobNet,
	)
	if err != nil {
		log.Fatalf("error creating bob's client: %v", err)
	}

	fmt.Println("Alice Balance:", alice.GetBalances())
	fmt.Println("Bob Balance:", bob.GetBalances())

	//Open Channel between Alice and Bob
	log.Println("Opening channel and depositing funds")
	chAlice := alice.OpenChannel(bob.WireAddress(), bob.PeerID(), map[channel.Asset]float64{
		setup.CKBAsset: 100.0,
	})

	log.Println("Alice sent proposal")
	//Bob accepts channel
	chBob := bob.AcceptedChannel()
	log.Println("Bob accepted proposal")

	assets := []asset.Asset{*setup.CKBAsset}
	printBalances(chAlice, assets)

	log.Println("Sending payments....")

	//Alice sends payment
	chAlice.SendPayment(map[channel.Asset]float64{
		setup.CKBAsset: 10.0,
	})
	log.Println("Alice sent Bob a payment")
	printBalances(chAlice, assets)

	//Bob sends payment
	chBob.SendPayment(map[channel.Asset]float64{
		setup.CKBAsset: 5.0,
	})
	log.Println("Bob sent Alice a payment")
	printBalances(chAlice, assets)

	log.Println("Payments completed")

	//Settling channels
	log.Println("Settle channels")
	chAlice.Settle()

	//Cleanup
	alice.Shutdown()
	bob.Shutdown()
	log.Println("Clients shutdown, exiting method")

}
