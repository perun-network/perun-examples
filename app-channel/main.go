// Copyright 2021 PolyCrypt GmbH, Germany
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
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/perun-examples/app-channel/client"
	"perun.network/perun-examples/app-channel/eth"
)

type Role int

const (
	RoleDeployer Role = 0
	RoleAlice    Role = 1
	RoleBob      Role = 2
)

type config struct {
	chainURL    string                     // Url of the Ethereum node.
	chainID     *big.Int                   // ID of the targeted chain
	hosts       map[Role]string            // Hosts for incoming connections.
	privateKeys map[Role]*ecdsa.PrivateKey // Private key of each role.
	addrs       map[Role]*wallet.Address   // Wallet address of each role.
}

// Test parameters
var (
	blockTime             = 1 * time.Second
	defaultContextTimeout = 30 * time.Second
	stake                 = eth.EthToWei(big.NewFloat(10))
	playerTimeout         = 30 * time.Second
	cfg                   config
)

// main performs an end-to-end run of the tic tac toe channel application.
func main() {
	alice, bob := setup()
	logAccountBalance(alice, bob)

	ctx, cancel := context.WithTimeout(context.Background(), playerTimeout)
	defer cancel() // releases resources if slowOperation completes before timeout elapses

	// Start Alice's routine
	aliceAction := make(chan interface{}, 1)
	aliceErr := make(chan error, 1)
	go playerRoutine(ctx, "Alice", alice, aliceAction, aliceErr)

	// Start Bob's routine
	bobAction := make(chan interface{}, 1)
	bobErr := make(chan error, 1)
	go playerRoutine(ctx, "Bob", bob, bobAction, bobErr)

	// Alice proposing game
	aliceAction <- propose{bob.Address(), stake}
	// Bob accepting game
	bobAction <- accept{alice.Address(), stake}

	if err := <-aliceErr; err != nil {
		panic(fmt.Errorf("alice proposing game: %v", err))
	}
	if err := <-bobErr; err != nil {
		panic(fmt.Errorf("bob accepting game: %v", err))
	}

	// Alice set (2, 0)
	aliceAction <- set{2, 0}
	if err := <-aliceErr; err != nil {
		panic(fmt.Errorf("alice set: %v", err))
	}

	// Bob set (0, 0)
	bobAction <- set{0, 0}
	if err := <-bobErr; err != nil {
		panic(fmt.Errorf("bob set: %v", err))
	}

	// Alice set (0, 2)
	aliceAction <- set{0, 2}
	if err := <-aliceErr; err != nil {
		panic(fmt.Errorf("alice set: %v", err))
	}

	// Bob set (1, 1)
	bobAction <- set{1, 1}
	if err := <-bobErr; err != nil {
		panic(fmt.Errorf("bob set: %v", err))
	}

	// Alice set (2, 2)
	aliceAction <- set{2, 2}
	if err := <-aliceErr; err != nil {
		panic(fmt.Errorf("alice set: %v", err))
	}

	// Bob set (2, 1)
	bobAction <- set{2, 1}
	if err := <-bobErr; err != nil {
		panic(fmt.Errorf("bob set: %v", err))
	}

	// Alice set (1, 2)
	aliceAction <- set{1, 2}
	if err := <-aliceErr; err != nil {
		panic(fmt.Errorf("alice set: %v", err))
	}

	// Alice concludes
	aliceAction <- conclude{}
	if err := <-aliceErr; err != nil {
		panic(fmt.Errorf("alice conclude game: %v", err))
	}

	// Bob concludes
	bobAction <- conclude{}
	if err := <-bobErr; err != nil {
		panic(fmt.Errorf("bob conclude game: %v", err))
	}

	logAccountBalance(alice, bob)
}

type (
	propose struct {
		opponent common.Address
		stake    *big.Int
	}
	accept struct {
		opponent common.Address
		stake    *big.Int
	}
	set struct {
		x, y int
	}
	conclude struct{}
)

func playerRoutine(ctx context.Context, name string, c *client.Client, actions chan interface{}, errors chan error) {
	var g *client.Game
	var err error
	logPrintf := func(format string, v ...interface{}) {
		log.Printf("%v: %v", name, fmt.Sprintf(format, v...))
	}

	for {
		select {
		case <-ctx.Done():
			errors <- ctx.Err()
			return
		case a := <-actions:
			switch a := a.(type) {
			case propose:
				logPrintf("Proposing game to %v with stake %v", a.opponent, toEth(a.stake))
				g, err = c.ProposeGame(a.opponent, a.stake)
				errors <- err
			case accept:
				proposal, err := c.NextGameProposal()
				if err != nil {
					errors <- err
					return
				}
				logPrintf("Accepting game proposal from %v with stake %v", a.opponent, toEth(a.stake))
				g, err = proposal.Accept()
				errors <- err
			case set:
				logPrintf("Setting (%d, %d)", a.x, a.y)
				errors <- g.Set(a.x, a.y)
				fmt.Println(g.String())
			case conclude:
				logPrintf("Concluding")
				errors <- g.Conclude()
			}
		}
	}
}

func setup() (*client.Client, *client.Client) {
	initConfig()

	// Deploy contracts
	log.Print("Deploying contracts...")
	nodeURL := cfg.chainURL
	deploymentKey := cfg.privateKeys[RoleDeployer]

	contracts, err := deployContracts(nodeURL, cfg.chainID, deploymentKey, defaultContextTimeout)
	if err != nil {
		panic(fmt.Errorf("deploying contracts: %v", err))
	}

	log.Print("Setting up clients...")
	// Setup Alice
	clientConfig1 := createClientConfig(
		nodeURL, contracts,
		cfg.privateKeys[RoleAlice], cfg.hosts[RoleAlice],
		cfg.addrs[RoleBob], cfg.hosts[RoleBob],
	)
	c1, err := client.StartClient(clientConfig1)
	if err != nil {
		panic(fmt.Errorf("alice setup: %v", err))
	}

	// Setup Bob
	clientConfig2 := createClientConfig(
		nodeURL, contracts,
		cfg.privateKeys[RoleBob], cfg.hosts[RoleBob],
		cfg.addrs[RoleAlice], cfg.hosts[RoleAlice],
	)
	c2, err := client.StartClient(clientConfig2)
	if err != nil {
		panic(fmt.Errorf("bob setup: %v", err))
	}
	log.Print("Setup done.")

	return c1, c2
}

func initConfig() {
	cfg.chainURL = "ws://127.0.0.1:8545"
	cfg.chainID = big.NewInt(1337)
	cfg.hosts = map[Role]string{
		RoleDeployer: "0.0.0.0:8401",
		RoleAlice:    "0.0.0.0:8402",
		RoleBob:      "0.0.0.0:8403",
	}

	// convert the private key strings ...
	rawKeys := []string{
		"50b4713b4ba55b6fbcb826ae04e66c03a12fc62886a90ca57ab541959337e897", // Deployer
		"1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f", // Alice
		"f63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e", // Bob
	}

	// ... to ECDSA keys:
	privateKeys := make(map[Role]*ecdsa.PrivateKey, len(rawKeys))
	for index, key := range rawKeys {
		privateKey, _ := crypto.HexToECDSA(key)
		privateKeys[Role(index)] = privateKey
	}
	cfg.privateKeys = privateKeys

	// Fix the on-chain addresses of the App Deployer, Alice and Bob.
	addresses := make(map[Role]*wallet.Address, len(rawKeys))
	for index, key := range cfg.privateKeys {
		commonAddress := crypto.PubkeyToAddress(key.PublicKey)
		addresses[index] = wallet.AsWalletAddr(commonAddress)
	}
	cfg.addrs = addresses
}

// createClientConfig is a helper function for client setup.
func createClientConfig(nodeURL string, contracts ContractAddresses, privateKey *ecdsa.PrivateKey, host string, peerAddress *wallet.Address, peerHost string) client.ClientConfig {
	return client.ClientConfig{
		SetupClientConfig: client.SetupClientConfig{
			PrivateKey:      privateKey,
			Host:            host,
			ETHNodeURL:      nodeURL,
			AdjudicatorAddr: contracts.AdjudicatorAddr,
			AssetHolderAddr: contracts.AssetHolderAddr,
			DialerTimeout:   1 * time.Second,
			PeerAddresses: []client.PeerWithAddress{
				{
					Peer:    peerAddress,
					Address: peerHost,
				},
			},
		},
		ChallengeDuration: 5 * blockTime,
		AppAddress:        contracts.AppAddr,
		ContextTimeout:    defaultContextTimeout,
	}
}
