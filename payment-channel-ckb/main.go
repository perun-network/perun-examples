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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
	"github.com/perun-network/perun-libp2p-wire/p2p"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/channel/persistence/keyvalue"
	"perun.network/perun-ckb-backend/channel/asset"
	"perun.network/perun-ckb-backend/wallet"
	"perun.network/perun-examples/payment-channel-ckb/client"
	"perun.network/perun-examples/payment-channel-ckb/deployment"
	"polycry.pt/poly-go/sortedkv/memorydb"
)

const (
	rpcNodeURL = "http://localhost:8114"
	Network    = types.NetworkTest
)

func main() {
	//Setup devnet environment
	log.Println("Deploying Devnet")
	sudtOwnerLockArg, err := parseSUDTOwnerLockArg("./devnet/accounts/sudt-owner-lock-hash.txt")
	if err != nil {
		log.Fatalf("error getting SUDT owner lock arg: %v", err)
	}
	d, _, err := deployment.GetDeployment("./devnet/contracts/migrations/dev/", "./devnet/system_scripts", sudtOwnerLockArg)
	if err != nil {
		log.Fatalf("error getting deployment: %v", err)
	}

	//Setup wallets
	log.Println("Creating wallets")
	w := wallet.NewEphemeralWallet()

	keyAlice, err := deployment.GetKey("./devnet/accounts/alice.pk")
	if err != nil {
		log.Fatalf("error getting alice's private key: %v", err)
	}
	keyBob, err := deployment.GetKey("./devnet/accounts/bob.pk")
	if err != nil {
		log.Fatalf("error getting bob's private key: %v", err)
	}
	aliceAccount := wallet.NewAccountFromPrivateKey(keyAlice)
	bobAccount := wallet.NewAccountFromPrivateKey(keyBob)

	err = w.AddAccount(aliceAccount)
	if err != nil {
		log.Fatalf("error adding alice's account: %v", err)
	}
	err = w.AddAccount(bobAccount)
	if err != nil {
		log.Fatalf("error adding bob's account: %v", err)
	}

	aliceWireAcc := p2p.NewRandomAccount(rand.New(rand.NewSource(time.Now().UnixNano())))
	aliceNet, err := p2p.NewP2PBus(aliceWireAcc)
	if err != nil {
		log.Fatalf("creating p2p net: %v", err)
	}
	aliceBus := aliceNet.Bus
	aliceListener := aliceNet.Listener
	go aliceBus.Listen(aliceListener)

	log.Println("Setting up payment channel clients")

	//Setup Payment Clients
	prAlice := keyvalue.NewPersistRestorer(memorydb.NewDatabase())

	alice, err := client.NewPaymentClient(
		"Alice",
		Network,
		d,
		rpcNodeURL,
		aliceAccount,
		*keyAlice,
		w,
		prAlice,
		aliceWireAcc.Address(),
		aliceNet,
	)
	if err != nil {
		log.Fatalf("error creating alice's client: %v", err)
	}

	prBob := keyvalue.NewPersistRestorer(memorydb.NewDatabase())
	bobWireAcc := p2p.NewRandomAccount(rand.New(rand.NewSource(time.Now().UnixNano())))
	bobNet, err := p2p.NewP2PBus(bobWireAcc)
	if err != nil {
		log.Fatalf("creating p2p net: %v", err)
	}
	bobBus := bobNet.Bus
	bobListener := bobNet.Listener
	go bobBus.Listen(bobListener)

	bob, err := client.NewPaymentClient(
		"Bob",
		Network,
		d,
		rpcNodeURL,
		bobAccount,
		*keyBob,
		w,
		prBob,
		bobWireAcc.Address(),
		bobNet,
	)
	if err != nil {
		log.Fatalf("error creating bob's client: %v", err)
	}

	ckbAsset := asset.Asset{
		IsCKBytes: true,
		SUDT:      nil,
	}

	fmt.Println("Alice Balance:", alice.GetBalances())
	fmt.Println("Bob Balance:", bob.GetBalances())

	//Open Channel between Alice and Bob
	log.Println("Opening channel and depositing funds")
	chAlice := alice.OpenChannel(bob.WireAddress(), bob.PeerID(), map[channel.Asset]float64{
		&asset.Asset{
			IsCKBytes: true,
			SUDT:      nil,
		}: 100.0,
	})

	log.Println("Alice sent proposal")
	//Bob accepts channel
	chBob := bob.AcceptedChannel()
	log.Println("Bob accepted proposal")

	printBalances(chAlice, ckbAsset)

	log.Println("Sending payments....")

	//Alice sends payment
	chAlice.SendPayment(map[channel.Asset]float64{
		&ckbAsset: 10.0,
	})
	log.Println("Alice sent Bob a payment")
	printBalances(chAlice, ckbAsset)

	//Bob sends payment
	chBob.SendPayment(map[channel.Asset]float64{
		&ckbAsset: 10.0,
	})
	log.Println("Bob sent Alice a payment")
	printBalances(chAlice, ckbAsset)

	chAlice.SendPayment(map[channel.Asset]float64{
		&ckbAsset: 10.0,
	})
	log.Println("Alice sent Bob a payment")
	printBalances(chAlice, ckbAsset)

	log.Println("Payments completed")

	//Settling channels
	log.Println("Settle channels")
	chAlice.Settle()

	//Cleanup
	alice.Shutdown()
	bob.Shutdown()
	log.Println("Clients shutdown, exiting method")

}

func printBalances(ch *client.PaymentChannel, asset asset.Asset) {
	chAlloc := ch.State().Allocation

	// Constants for formatting CKBytes
	const ckbyteConversionFactor = 100_000_000 // 1 CKByte = 100,000,000 smallest units

	// Log general information
	log.Println("=== Allocation Balances ===")

	// Get Alice's balance (participant 0)
	aliceBalance := chAlloc.Balance(0, &asset)
	aliceBalanceCKBytes := new(big.Float).Quo(new(big.Float).SetInt(aliceBalance), big.NewFloat(ckbyteConversionFactor))

	// Get Bob's balance (participant 1)
	bobBalance := chAlloc.Balance(1, &asset)
	bobBalanceCKBytes := new(big.Float).Quo(new(big.Float).SetInt(bobBalance), big.NewFloat(ckbyteConversionFactor))

	// Print Alice's balance
	log.Printf("Alice's allocation: %s CKBytes", aliceBalanceCKBytes.Text('f', 2))

	// Print Bob's balance
	log.Printf("Bob's allocation: %s CKBytes", bobBalanceCKBytes.Text('f', 2))

	// Calculate the total balance
	totalBalance := new(big.Int).Add(aliceBalance, bobBalance)
	totalBalanceCKBytes := new(big.Float).Quo(new(big.Float).SetInt(totalBalance), big.NewFloat(ckbyteConversionFactor))

	// Print the total channel balance
	log.Printf("Total channel balance: %s CKBytes", totalBalanceCKBytes.Text('f', 2))

	log.Println("===========================")
}

func parseSUDTOwnerLockArg(pathToSUDTOwnerLockArg string) (string, error) {
	b, err := ioutil.ReadFile(pathToSUDTOwnerLockArg)
	if err != nil {
		return "", fmt.Errorf("reading sudt owner lock arg from file: %w", err)
	}
	sudtOwnerLockArg := string(b)
	if sudtOwnerLockArg == "" {
		return "", errors.New("sudt owner lock arg not found in file")
	}
	return sudtOwnerLockArg, nil
}
