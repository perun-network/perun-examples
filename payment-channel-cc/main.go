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
	"github.com/ethereum/go-ethereum/crypto"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	"log"
	"perun.network/go-perun/wire"
	"perun.network/perun-examples/payment-channel-cc/ethereumUtil"
	"perun.network/perun-examples/payment-channel-cc/stellarUtil"
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
	adjudicator, assetHolder := ethereumUtil.DeployContracts(chainURL, 1337, keyDeployer)
	log.Println("Adjudicator:", adjudicator.Hex())
	log.Println("Asset holder:", assetHolder.Hex())

	asset := *ethwallet.AsWalletAddr(assetHolder)

	// Setup clients.
	log.Println("Setting up clients.")
	kAlice, err := crypto.HexToECDSA(keyAlice)
	if err != nil {
		panic(err)
	}
	kBob, err := crypto.HexToECDSA(keyBob)
	if err != nil {
		panic(err)
	}
	setup, err := stellarUtil.NewExampleSetup([]string{keyAlice, keyBob}, [][20]byte{crypto.PubkeyToAddress(kAlice.PublicKey), crypto.PubkeyToAddress(kBob.PublicKey)})
	if err != nil {
		panic(err)
	}
	bus := wire.NewLocalBus() // Message bus used for off-chain communication.
	alice := ethereumUtil.SetupPaymentClient(bus, chainURL, adjudicator, asset, kAlice, setup.GetWallets()[0], setup.GetAccounts()[0], setup.GetTokenAsset()[0], setup.GetFunders()[0], setup.GetAdjudicators()[0])
	bob := ethereumUtil.SetupPaymentClient(bus, chainURL, adjudicator, asset, kBob, setup.GetWallets()[1], setup.GetAccounts()[1], setup.GetTokenAsset()[1], setup.GetFunders()[1], setup.GetAdjudicators()[1])

	log.Println("Participants: ", alice.WireAddress(), bob.WireAddress())
	// Print balances before transactions.
	l := ethereumUtil.NewBalanceLogger(chainURL)
	l.LogBalances(alice.WalletEthAddress(), bob.WalletEthAddress())

	// Open channel, transact, close.
	log.Println("Opening channel and depositing funds.")
	chAlice := alice.OpenChannel(bob.WireAddress(), 1, 50)
	chBob := bob.AcceptedChannel()

	log.Println("Sending payments...")
	chAlice.SendEthPayment(1)
	chBob.SendStellarPayment(50)

	log.Println("Settling channel.")
	chAlice.Settle() // Conclude and withdraw.
	chBob.Settle()   // Withdraw.

	// Print balances after transactions.
	l.LogBalances(alice.WalletEthAddress(), bob.WalletEthAddress())

	// Cleanup.
	alice.Shutdown()
	bob.Shutdown()
}
