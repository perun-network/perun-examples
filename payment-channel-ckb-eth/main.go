// Copyright 2025 PolyCrypt GmbH
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
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nervosnetwork/ckb-sdk-go/v2/rpc"
	ckbsigner "github.com/nervosnetwork/ckb-sdk-go/v2/transaction/signer"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	"io/ioutil"
	"log"
	"math/big"
	"perun.network/go-perun/wire"
	"perun.network/perun-ckb-backend/backend"
	"perun.network/perun-ckb-backend/channel/adjudicator"
	ckbasset "perun.network/perun-ckb-backend/channel/asset"
	"perun.network/perun-ckb-backend/channel/funder"
	ckbclient "perun.network/perun-ckb-backend/client"
	"perun.network/perun-ckb-backend/wallet"
	ckbaddress "perun.network/perun-ckb-backend/wallet/address"
	"perun.network/perun-examples/payment-channel-ckb-eth/client"
	"perun.network/perun-examples/payment-channel-ckb-eth/deployment"
	"perun.network/perun-examples/payment-channel-ckb-eth/ethereumUtil"
)

const (
	chainURL = "ws://127.0.0.1:8545"

	// Private keys.
	keyDeployer = "79ea8f62d97bc0591a4224c1725fca6b00de5b2cea286fe2e0bb35c5e76be46e"
	keyAlice    = "1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f"
	keyBob      = "f63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e"
)

// main runs a demo of the payment client. It assumes that a blockchain node is
// available at `chainURL` and that the accounts corresponding to the specified
// secret keys are provided with sufficient funds.
func main() {

	// Configure log flags: date/time and file/line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Deploy contracts.
	log.Println("Deploying contracts.")
	ethadjudicator, assetHolder := ethereumUtil.DeployContracts(chainURL, 1337, keyDeployer)
	log.Println("Adjudicator:", ethadjudicator.Hex())
	log.Println("Asset holder:", assetHolder.Hex())

	asset := *ethwallet.AsWalletAddr(assetHolder)

	// Setup clients.
	sudtOwnerLockArg, err := parseSUDTOwnerLockArg("./devnet/accounts/sudt-owner-lock-hash1.txt")
	if err != nil {
		panic(err)
	}
	d, _, err := deployment.GetDeployment("./devnet/contracts/migrations/dev/", "./devnet/contracts/migrations_vc/dev/", "./devnet/system_scripts", []string{sudtOwnerLockArg})
	if err != nil {
		log.Fatalf("error getting deployment: %v", err)
	}

	//Setup wallets
	log.Println("Creating wallets")
	kAlice, err := crypto.HexToECDSA(keyAlice)
	if err != nil {
		panic(err)
	}
	kBob, err := crypto.HexToECDSA(keyBob)
	if err != nil {
		panic(err)
	}
	ckbWalletA := wallet.NewEphemeralWallet()
	ckbWalletB := wallet.NewEphemeralWallet()

	keyAliceCkb, err := deployment.GetKey("./devnet/accounts/alice.pk")
	if err != nil {
		log.Fatalf("error getting alice's private key: %v", err)
	}
	keyBobCkb, err := deployment.GetKey("./devnet/accounts/bob.pk")
	if err != nil {
		log.Fatalf("error getting bob's private key: %v", err)
	}
	ethWalletA := swallet.NewWallet(kAlice)
	accA := crypto.PubkeyToAddress(kAlice.PublicKey)
	eaddrA := ethwallet.AsWalletAddr(accA)
	ethWalletB := swallet.NewWallet(kBob)
	accB := crypto.PubkeyToAddress(kBob.PublicKey)
	eaddrB := ethwallet.AsWalletAddr(accB)
	omniHash := d.OmniLockScript.CodeHash
	alicePart, authDataAlice, _ := ckbaddress.NewEthereumParticipantFromPublicKey(keyAliceCkb.PubKey(), omniHash)
	bobPart, authDataBob, _ := ckbaddress.NewEthereumParticipantFromPublicKey(keyBobCkb.PubKey(), omniHash)
	aliceAccount := wallet.NewAccountFromPrivateKey(keyAliceCkb, omniHash, false)
	bobAccount := wallet.NewAccountFromPrivateKey(keyBobCkb, omniHash, false)
	ckbWalletA.AddAccount(aliceAccount)
	ckbWalletB.AddAccount(bobAccount)

	ckbAsset := ckbasset.NewCKBytesAsset()
	id := ckbasset.MakeCCID(ckbasset.MakeContractID("03"))
	ca := ckbasset.NewNervosAsset(*ckbAsset, id)
	network := types.NetworkTest
	backendRPCClientAlice, err := rpc.Dial("http://localhost:8114")
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	backendRPCClientBob, err := rpc.Dial("http://localhost:8114")
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	signerA := backend.NewEVMSignerInstance(alicePart.ToCKBAddress(network), *keyAliceCkb, network, authDataAlice)
	txSignerA := signerA.Signer()
	txSignerA.RegisterLockSigner(d.OmniLockScript.CodeHash, &ckbsigner.OmnilockSigner{})
	ckbClientA, _ := ckbclient.NewClient(backendRPCClientAlice, signerA, d)
	ckbFunderA := funder.NewDefaultFunder(ckbClientA, d)
	ckbFunderA.MaxIterationsUntilAbort = 80
	ckbAdjudicatorA := adjudicator.NewAdjudicator(ckbClientA)

	signerB := backend.NewEVMSignerInstance(bobPart.ToCKBAddress(network), *keyBobCkb, network, authDataBob)
	txSignerB := signerB.Signer()
	txSignerB.RegisterLockSigner(d.OmniLockScript.CodeHash, &ckbsigner.OmnilockSigner{})
	ckbClientB, _ := ckbclient.NewClient(backendRPCClientBob, signerB, d)
	ckbFunderB := funder.NewDefaultFunder(ckbClientB, d)
	ckbFunderB.MaxIterationsUntilAbort = 80
	ckbAdjudicatorB := adjudicator.NewAdjudicator(ckbClientB)

	bus := wire.NewLocalBus() // Message bus used for off-chain communication.
	alice, _ := client.SetupPaymentClient(bus, ethWalletA, accA, eaddrA, chainURL, 1337, ethadjudicator, asset, ckbWalletA, aliceAccount, &ca, ckbFunderA, ckbAdjudicatorA)

	bob, _ := client.SetupPaymentClient(bus, ethWalletB, accB, eaddrB, chainURL, 1337, ethadjudicator, asset, ckbWalletB, bobAccount, &ca, ckbFunderB, ckbAdjudicatorB)

	log.Println("Participants: ", alice.WireAddress(), bob.WireAddress())

	// Print balances before transactions.
	l := ethereumUtil.NewBalanceLogger(chainURL)
	l.LogBalances(alice.WalletEthAddress(), bob.WalletEthAddress())
	ckbLogger := ethereumUtil.NewCKBBalanceLogger("http://localhost:8114")
	ckbLogger.LogBalances(aliceAccount.Address().(*ckbaddress.Participant))
	ckbLogger.LogBalances(bobAccount.Address().(*ckbaddress.Participant))

	// Open channel, transact, close.
	log.Println("Opening channel and depositing funds.")
	chAlice := alice.OpenChannel(bob.WireAddress(), 1, 50)
	log.Println("Channel accepted by Bob.")
	chBob := bob.AcceptedChannel()

	log.Println("Sending payments...")
	chAlice.SendEthPayment(1)
	chBob.SendCKBPayment(50)
	log.Println("Alice sent Bob a payment")
	printBalances(chAlice, ca)
	log.Println("Settling channel.")
	chAlice.Settle() // Settle

	// Print balances after transactions.
	l.LogBalances(alice.WalletEthAddress(), bob.WalletEthAddress())
	ckbLogger.LogBalances(aliceAccount.Address().(*ckbaddress.Participant))
	ckbLogger.LogBalances(bobAccount.Address().(*ckbaddress.Participant))

	// Cleanup.
	alice.Shutdown()
	bob.Shutdown()
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
func printBalances(ch *client.PaymentChannel, asset ckbasset.NervosAsset) {
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
