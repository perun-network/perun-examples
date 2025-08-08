package main

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"perun.network/perun-ckb-backend/backend"
	"perun.network/perun-ckb-backend/channel/asset"
	"perun.network/perun-ckb-backend/channel/test"
	"perun.network/perun-ckb-backend/wallet"
	ckbwallet "perun.network/perun-ckb-backend/wallet"
	"perun.network/perun-examples/payment-channel-ckb/client"
	"perun.network/perun-examples/payment-channel-ckb/deployment"
)

const (
	sudtMaxCapacity = 200_00_000_000 // 200 ckb
)

// Setup contains all the necessary information for CKB payment channel setup.
type Setup struct {
	Deployment backend.Deployment
	SUDTInfo   deployment.SUDTInfo
	Wallets    []*ckbwallet.EphemeralWallet
	WalletAccs []*ckbwallet.Account
	CKBAsset   *asset.Asset
	SudtAsset  *asset.Asset
	AccKeys    []*secp256k1.PrivateKey
}

// NewSetup creates a new Setup instance with the provided parameters.
func NewSetup() *Setup {
	log.Println("Initializing CKB payment channel setup")

	sudtOwnerLockArg, err := parseSUDTOwnerLockArg("./devnet/accounts/sudt-owner-lock-hash.txt")
	if err != nil {
		log.Fatalf("error getting SUDT owner lock arg: %v", err)
	}

	d, sudtInfo, err := deployment.GetDeployment("./devnet/contracts/migrations/dev/", "./devnet/contracts/migrations_vc/dev/", "./devnet/system_scripts", sudtOwnerLockArg)
	if err != nil {
		log.Fatalf("error getting deployment: %v", err)
	}

	//Setup wallets
	log.Println("Creating wallets")
	wAlice := wallet.NewEphemeralWallet()
	wBob := wallet.NewEphemeralWallet()

	keyAlice, err := test.GetKey("./devnet/accounts/alice.pk")
	if err != nil {
		log.Fatalf("error getting alice's private key: %v", err)
	}
	keyBob, err := test.GetKey("./devnet/accounts/bob.pk")
	if err != nil {
		log.Fatalf("error getting bob's private key: %v", err)
	}

	aliceAccount := wallet.NewAccountFromPrivateKey(keyAlice)
	bobAccount := wallet.NewAccountFromPrivateKey(keyBob)

	err = wAlice.AddAccount(aliceAccount)
	if err != nil {
		log.Fatalf("error adding alice's account: %v", err)
	}
	err = wBob.AddAccount(bobAccount)
	if err != nil {
		log.Fatalf("error adding bob's account: %v", err)
	}

	ckbAsset := &asset.Asset{
		IsCKBytes: true,
		SUDT:      nil,
	}

	sudtAsset := &asset.Asset{
		IsCKBytes: false,
		SUDT:      asset.NewSUDT(*sudtInfo.Script, uint64(sudtMaxCapacity)),
	}
	return &Setup{
		Deployment: d,
		SUDTInfo:   sudtInfo,
		Wallets:    []*ckbwallet.EphemeralWallet{wAlice, wBob},
		WalletAccs: []*ckbwallet.Account{aliceAccount, bobAccount},
		CKBAsset:   ckbAsset,
		SudtAsset:  sudtAsset,
		AccKeys:    []*secp256k1.PrivateKey{keyAlice, keyBob},
	}
}

func parseSUDTOwnerLockArg(pathToSUDTOwnerLockArg string) (string, error) {
	b, err := os.ReadFile(pathToSUDTOwnerLockArg)
	if err != nil {
		return "", fmt.Errorf("reading sudt owner lock arg from file: %w", err)
	}
	sudtOwnerLockArg := string(b)
	if sudtOwnerLockArg == "" {
		return "", errors.New("sudt owner lock arg not found in file")
	}
	return sudtOwnerLockArg, nil
}

func printBalances(ch *client.PaymentChannel, assets []asset.Asset) {
	chAlloc := ch.State().Allocation

	// Constants for formatting CKBytes
	const ckbyteConversionFactor = 100_000_000 // 1 CKByte = 100,000,000 smallest units
	const sudtConversionFactor = 1             // SUDT is already in the smallest unit

	// Log general information
	log.Println("=== Allocation Balances ===")
	for _, asset := range assets {
		if asset.IsCKBytes {

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
		} else {
			// Get Alice's balance (participant 0)
			aliceBalance := chAlloc.Balance(0, &asset)
			aliceBalanceSUDTs := new(big.Float).Quo(new(big.Float).SetInt(aliceBalance), big.NewFloat(sudtConversionFactor))

			// Get Bob's balance (participant 1)
			bobBalance := chAlloc.Balance(1, &asset)
			bobBalanceSUDTs := new(big.Float).Quo(new(big.Float).SetInt(bobBalance), big.NewFloat(sudtConversionFactor))

			// Print Alice's balance
			log.Printf("Alice's allocation: %s SUDT", aliceBalanceSUDTs.Text('f', 2))

			// Print Bob's balance
			log.Printf("Bob's allocation: %s SUDT", bobBalanceSUDTs.Text('f', 2))

			// Calculate the total balance
			totalBalance := new(big.Int).Add(aliceBalance, bobBalance)
			totalBalanceSUDTs := new(big.Float).Quo(new(big.Float).SetInt(totalBalance), big.NewFloat(sudtConversionFactor))

			// Print the total channel balance
			log.Printf("Total channel balance: %s SUDT", totalBalanceSUDTs.Text('f', 2))

			log.Println("===========================")
		}
	}
}
