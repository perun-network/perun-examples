package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/channel/persistence/keyvalue"
	"perun.network/go-perun/wire"
	"perun.network/perun-ckb-backend/channel/asset"
	"perun.network/perun-ckb-backend/wallet"
	"perun.network/perun-ckb-demo/client"
	"perun.network/perun-ckb-demo/deployment"
	"polycry.pt/poly-go/sortedkv/memorydb"
)

const (
	rpcNodeURL = "http://localhost:8114"
	Network    = types.NetworkTest
)

func SetLogFile(path string) {
	logFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(logFile)
}

func main() {
	SetLogFile("demo.log")
	sudtOwnerLockArg, err := parseSUDTOwnerLockArg("./devnet/accounts/sudt-owner-lock-hash.txt")
	if err != nil {
		log.Fatalf("error getting SUDT owner lock arg: %v", err)
	}
	d, _, err := deployment.GetDeployment("./devnet/contracts/migrations/dev/", "./devnet/system_scripts", sudtOwnerLockArg)
	if err != nil {
		log.Fatalf("error getting deployment: %v", err)
	}
	/*
		maxSudtCapacity := transaction.CalculateCellCapacity(types.CellOutput{
			Capacity: 0,
			Lock:     &d.DefaultLockScript,
			Type:     sudtInfo.Script,
		})
	*/
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

	// Setup clients.
	log.Println("Setting up clients.")
	bus := wire.NewLocalBus() // Message bus used for off-chain communication.
	prAlice := keyvalue.NewPersistRestorer(memorydb.NewDatabase())
	prBob := keyvalue.NewPersistRestorer(memorydb.NewDatabase())
	alice, err := client.NewPaymentClient(
		"Alice",
		Network,
		d,
		bus,
		rpcNodeURL,
		aliceAccount,
		*keyAlice,
		w,
		prAlice,
	)
	if err != nil {
		log.Fatalf("error creating alice's client: %v", err)
	}
	bob, err := client.NewPaymentClient(
		"Bob",
		Network,
		d,
		bus,
		rpcNodeURL,
		bobAccount,
		*keyBob,
		w,
		prBob,
	)
	if err != nil {
		log.Fatalf("error creating bob's client: %v", err)
	}

	//print balances before transaction

	fmt.Println("Balances of Alice and Bob before transaction")
	str := "'s account balance"
	fmt.Println(alice.Name, str, alice.GetBalances())
	fmt.Println(bob.Name, str, bob.GetBalances())

	ckbAsset := asset.Asset{
		IsCKBytes: true,
		SUDT:      nil,
	}

	/*
		sudtAsset := asset.Asset{
			IsCKBytes: false,
			SUDT: &asset.SUDT{
				TypeScript:  *sudtInfo.Script,
				MaxCapacity: maxSudtCapacity,
			},
		}
	*/

	fmt.Println("Opening channel and depositing funds")
	chAlice := alice.OpenChannel(bob.WireAddress(), map[channel.Asset]float64{
		&asset.Asset{
			IsCKBytes: true,
			SUDT:      nil,
		}: 100.0,
	})
	strAlice := "Alice"
	strBob := "Bob"
	fmt.Println(alice.Name, str, alice.GetBalances())
	fmt.Println(bob.Name, str, bob.GetBalances())

	fmt.Println("Alice sent proposal")
	chBob := bob.AcceptedChannel()
	fmt.Println("Bob accepted proposal")
	fmt.Println("Sending payments....")

	chAlice.SendPayment(map[channel.Asset]float64{
		&ckbAsset: 10.0,
	})
	fmt.Println("Alice sent Bob a payment")
	printAllocationBalances(chAlice, ckbAsset, strAlice)
	printAllocationBalances(chBob, ckbAsset, strBob)

	chBob.SendPayment(map[channel.Asset]float64{
		&ckbAsset: 10.0,
	})
	fmt.Println("Bob sent Alice a payment")
	printAllocationBalances(chAlice, ckbAsset, strAlice)
	printAllocationBalances(chBob, ckbAsset, strBob)

	chAlice.SendPayment(map[channel.Asset]float64{
		&ckbAsset: 10.0,
	})
	fmt.Println("Alice sent Bob a payment")
	printAllocationBalances(chAlice, ckbAsset, strAlice)
	printAllocationBalances(chBob, ckbAsset, strBob)

	fmt.Println("Payments completed")
	printAllocationBalances(chAlice, ckbAsset, strAlice)
	printAllocationBalances(chBob, ckbAsset, strBob)

	fmt.Println("Skip Settling Channel and force client shutdown")
	//chAlice.Settle()

	fmt.Println(alice.Name, str, alice.GetBalances())
	fmt.Println(bob.Name, str, bob.GetBalances())

	//cleanup
	alice.Shutdown()
	bob.Shutdown()
	fmt.Println("Clients shutdown, exiting method")

	fmt.Println("Creating clients again to see if channels can be restored")
	alice2, err := client.NewPaymentClient(
		"Alice",
		Network,
		d,
		bus,
		rpcNodeURL,
		aliceAccount,
		*keyAlice,
		w,
		prAlice,
	)
	if err != nil {
		log.Fatalf("error creating alice's client: %v", err)
	}
	bob2, err := client.NewPaymentClient(
		"Bob",
		Network,
		d,
		bus,
		rpcNodeURL,
		bobAccount,
		*keyBob,
		w,
		prBob,
	)
	if err != nil {
		log.Fatalf("error creating bob's client: %v", err)
	}

	chansAlice := alice2.Restore()
	chansBob := bob2.Restore()
	fmt.Println("Alice and Bob's channels successfully restored")

	// Print balances after transactions.
	fmt.Println(alice.Name, str, alice.GetBalances())
	fmt.Println(bob.Name, str, bob.GetBalances())

	fmt.Println("Alice sending payment to Bob")
	chansAlice[0].SendPayment(map[channel.Asset]float64{
		&ckbAsset: 10.0,
	})
	fmt.Println("Bob sending payment to Alice")
	chansBob[0].SendPayment(map[channel.Asset]float64{
		&ckbAsset: 10.0,
	})

	chansAlice[0].Settle()
	fmt.Println("Balances after settling channel")
	fmt.Println(alice.Name, str, alice.GetBalances())
	fmt.Println(bob.Name, str, bob.GetBalances())

}

func printAllocationBalances(ch *client.PaymentChannel, asset asset.Asset, name string) {
	chAlloc := ch.State().Allocation
	//_assets := chAlloc.Assets
	fmt.Println("Assets held by" + name)
	/*
		for _, a := range _assets {
			fmt.Println(a)
		}
	*/
	fmt.Println(name + "'s allocation in channel: " + chAlloc.Balance(1, &asset).String())
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
