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
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/perun-examples/payment-channel/client"
)

const (
	chainURL              = "ws://127.0.0.1:8545"
	chainID               = 1337
	defaultContextTimeout = 15 * time.Second
	deploymentKey         = "" // Key used for contract deployment. //TODO insert key

	// Alice
	hostAlice = "0.0.0.0:8401"
	keyAlice  = "1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f"

	// Bob
	hostBob = "0.0.0.0:8402"
	keyBob  = "f63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e"
)

func main() {
	contracts := deployContracts(chainURL, privateKey) //TODO

	alice := setupClient(hostAlice, keyAlice, chainURL, contracts)
	bob := setupClient(hostBob, keyBob, chainURL, contracts)

	logAccountBalance(alice, bob)

	alice.OpenChannel()
	alice.UpdateChannel()
	bob.UpdateChannel()
	bob.CloseChannel()

	logAccountBalance(alice, bob)
}

// setupClient sets up a new client with the given parameters.
func setupClient(
	nodeURL string,
	contracts ContractAddresses,
	host string,
	privateKey string,
) *client.Client {
	c, err := client.StartClient(clientConfig)

	return c1, c2
}

// initConfig initializes the config for a test run via ganache-cli
func initConfig() {
	// ... to ECDSA keys:
	privateKeys := make(map[client.Role]*ecdsa.PrivateKey, len(rawKeys))
	for index, key := range rawKeys {
		privateKey, _ := crypto.HexToECDSA(key)
		privateKeys[client.Role(index)] = privateKey
	}
	cfg.privateKeys = privateKeys

	// Fix the on-chain addresses of Alice and Bob.
	addresses := make(map[client.Role]*ethwallet.Address, len(rawKeys))
	for index, key := range cfg.privateKeys {
		commonAddress := crypto.PubkeyToAddress(key.PublicKey)
		addresses[index] = ethwallet.AsWalletAddr(commonAddress)
	}
	cfg.addrs = addresses
}

// createClientConfig is a helper function for client setup.
func createClientConfig(nodeURL string, contracts ContractAddresses, privateKey *ecdsa.PrivateKey, host string, peerAddress *ethwallet.Address, peerHost string) client.PaymentClientConfig {
	return client.PaymentClientConfig{
		PerunClientConfig: client.PerunClientConfig{
			PrivateKey:      privateKey,
			Host:            host,
			EthNodeURL:      nodeURL,
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
	}
}
