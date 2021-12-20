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
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"time"

	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/perun-examples/payment-channel/client"
)

type config struct {
	chainURL    string                            // Url of the Ethereum node.
	chainID     *big.Int                          // ID of the targeted chain
	hosts       map[client.Role]string            // Hosts for incoming connections.
	privateKeys map[client.Role]*ecdsa.PrivateKey // Private keys.
	addrs       map[client.Role]*wallet.Address   // Wallet addresses.
}

// Test parameters
var (
	cfg                   config
	defaultContextTimeout = 15 * time.Second
)

// main performs the opening, updating and closing of a simple perun payment channel
func main() {
	alice, bob := setup()
	logAccountBalance(alice, bob)

	if err := bob.OpenChannel(cfg.addrs[client.RoleAlice]); err != nil {
		panic(fmt.Errorf("opening channel: %w", err))
	}
	if err := bob.UpdateChannel(); err != nil {
		panic(fmt.Errorf("updating channel: %w", err))
	}
	if err := bob.CloseChannel(); err != nil {
		panic(fmt.Errorf("closing channel: %w", err))
	}

	logAccountBalance(alice, bob)
}

// setup creates the two clients Alice and Bob
func setup() (*client.Client, *client.Client) {
	initConfig()

	// Deploy contracts (Bob deploys)
	fmt.Println("Deploying contracts...")
	nodeURL := cfg.chainURL
	contracts, err := deployContracts(nodeURL, cfg.chainID, cfg.privateKeys[client.RoleBob], defaultContextTimeout)
	if err != nil {
		panic(fmt.Errorf("deploying contracts: %v", err))
	}

	fmt.Println("Setting up clients...")
	// Setup Alice
	clientConfig1 := createClientConfig(
		client.RoleAlice, nodeURL, contracts,
		cfg.privateKeys[client.RoleAlice], cfg.hosts[client.RoleAlice],
		cfg.addrs[client.RoleBob], cfg.hosts[client.RoleBob],
	)

	c1, err := client.StartClient(clientConfig1)
	if err != nil {
		panic(fmt.Errorf("alice setup: %v", err))
	}

	// Setup Bob
	clientConfig2 := createClientConfig(
		client.RoleBob, nodeURL, contracts,
		cfg.privateKeys[client.RoleBob], cfg.hosts[client.RoleBob],
		cfg.addrs[client.RoleAlice], cfg.hosts[client.RoleAlice],
	)

	c2, err := client.StartClient(clientConfig2)

	if err != nil {
		panic(fmt.Errorf("bob setup: %v", err))
	}
	fmt.Println("Setup done.")

	return c1, c2
}

// initConfig initializes the config for a test run via ganache-cli
func initConfig() {
	cfg.chainURL = "ws://127.0.0.1:8545"
	cfg.chainID = big.NewInt(1337)
	cfg.hosts = map[client.Role]string{
		client.RoleAlice: "0.0.0.0:8400",
		client.RoleBob:   "0.0.0.0:8401",
	}

	// convert the private key strings ...
	rawKeys := []string{
		"1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f", // Alice
		"f63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e", // Bob
	}

	// ... to ECDSA keys:
	privateKeys := make(map[client.Role]*ecdsa.PrivateKey, len(rawKeys))
	for index, key := range rawKeys {
		privateKey, _ := crypto.HexToECDSA(key)
		privateKeys[client.Role(index)] = privateKey
	}
	cfg.privateKeys = privateKeys

	// Fix the on-chain addresses of Alice and Bob.
	addresses := make(map[client.Role]*wallet.Address, len(rawKeys))
	for index, key := range cfg.privateKeys {
		commonAddress := crypto.PubkeyToAddress(key.PublicKey)
		addresses[index] = wallet.AsWalletAddr(commonAddress)
	}
	cfg.addrs = addresses
}

// createClientConfig is a helper function for client setup.
func createClientConfig(role client.Role, nodeURL string, contracts ContractAddresses, privateKey *ecdsa.PrivateKey, host string, peerAddress *wallet.Address, peerHost string) client.ClientConfig {
	return client.ClientConfig{
		SetupClientConfig: client.SetupClientConfig{
			Role:            role,
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
		ContextTimeout: defaultContextTimeout,
	}
}
