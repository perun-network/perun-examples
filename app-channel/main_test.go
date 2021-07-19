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
	"log"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/perun-examples/app-channel/client"
	"perun.network/perun-examples/app-channel/eth"
	"perun.network/perun-examples/app-channel/ganache"
	"perun.network/perun-examples/app-channel/perun"
)

// Test parameters
var (
	blockTimeInSeconds        = 1 * time.Second
	defaultContextTimeout     = 30 * time.Second
	collateralWithdrawalDelay = 10 * blockTimeInSeconds
	hostClient1               = "127.0.0.1:8546"
	hostClient2               = "127.0.0.1:8547"
	stake                     = eth.EthToWei(big.NewFloat(10))
	playerTimeout             = 30 * time.Second
)

// TestTicTacToeApp is an end-to-end test of the tic tac toe channel application.
func TestTicTacToeApp(t *testing.T) {
	c1, c2 := setupClients(t)
	logAccountBalance(c1, c2)

	require := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), playerTimeout)
	defer cancel()

	// Start Player1 routine
	player1Action := make(chan interface{}, 1)
	player1Err := make(chan error, 1)
	go playerRoutine(ctx, "Player1", c1, player1Action, player1Err)

	// Start Player2 routine
	player2Action := make(chan interface{}, 1)
	player2Err := make(chan error, 1)
	go playerRoutine(ctx, "Player2", c2, player2Action, player2Err)

	// Player1 propose game, Player2 accept game
	player1Action <- propose{c2.Address(), stake}
	player2Action <- accept{c1.Address(), stake}
	require.NoError(<-player1Err)
	require.NoError(<-player2Err)

	// Player1 set (2, 0)
	player1Action <- set{2, 0}
	require.NoError(<-player1Err)

	// Player2 set (0, 0)
	player2Action <- set{0, 0}
	require.NoError(<-player2Err)

	// Player1 set (0, 2)
	player1Action <- set{0, 2}
	require.NoError(<-player1Err)

	// Player2 set (1, 1)
	player2Action <- set{1, 1}
	require.NoError(<-player2Err)

	// Player1 set (2, 2)
	player1Action <- set{2, 2}
	require.NoError(<-player1Err)

	// Player2 set (2, 1)
	player2Action <- set{2, 1}
	require.NoError(<-player2Err)

	// Player1 set (1, 2)
	player1Action <- set{1, 2}
	require.NoError(<-player1Err)

	// Player1 conclude
	player1Action <- conclude{}
	require.NoError(<-player1Err)

	// Player2 conclude
	player2Action <- conclude{}
	require.NoError(<-player2Err)

	logAccountBalance(c1, c2)
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

func setupClients(t *testing.T) (*client.Client, *client.Client) {
	require := require.New(t)

	// Ganache config
	ganacheCfg := createGanacheConfig()

	// Start ganache blockchain with prefunded accounts
	log.Print("Starting local blockchain...")
	ganache, err := ganache.StartGanacheWithPrefundedAccounts(ganacheCfg)
	require.NoError(err, "starting ganache")
	t.Cleanup(func() {
		err := ganache.Shutdown()
		if err != nil {
			log.Print("shutting down ganaceh:", err)
		}
	})

	// Deploy contracts
	log.Print("Deploying contracts...")
	nodeURL := ganacheCfg.NodeURL()
	deploymentKey := ganache.Accounts[0].PrivateKey
	contracts, err := deployContracts(nodeURL, ganacheCfg.ChainID, deploymentKey, defaultContextTimeout, collateralWithdrawalDelay)
	require.NoError(err, "deploying contracts")

	log.Print("Setting up clients...")
	// Setup Client1
	clientConfig1 := createClientConfig(
		nodeURL, contracts,
		ganache.Accounts[1].PrivateKey, hostClient1,
		ganache.Accounts[2].Address(), hostClient2,
	)
	c1, err := client.StartClient(clientConfig1)
	require.NoError(err, "Client1 setup")

	// Setup Client2
	clientConfig2 := createClientConfig(
		nodeURL, contracts,
		ganache.Accounts[2].PrivateKey, hostClient2,
		ganache.Accounts[1].Address(), hostClient1,
	)
	c2, err := client.StartClient(clientConfig2)
	require.NoError(err, "Client2 setup")
	log.Print("Setup done.")

	return c1, c2
}

func createGanacheConfig() ganache.GanacheConfig {
	// Accounts
	accountFunding := []ganache.KeyWithBalance{
		{PrivateKey: "0x50b4713b4ba55b6fbcb826ae04e66c03a12fc62886a90ca57ab541959337e897", BalanceEth: 10},  // Contract Deployer
		{PrivateKey: "0x1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f", BalanceEth: 100}, // Client1
		{PrivateKey: "0xf63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e", BalanceEth: 200}, // Client2
	}

	ganacheCmd := os.Getenv("GANACHE_CMD")
	if len(ganacheCmd) == 0 {
		ganacheCmd = "ganache-cli"
	}
	return ganache.GanacheConfig{
		Cmd:         ganacheCmd,
		Host:        "127.0.0.1",
		Port:        8545,
		BlockTime:   blockTimeInSeconds,
		Funding:     accountFunding,
		StartupTime: 3 * time.Second,
		ChainID:     big.NewInt(1337),
	}
}

// genClientConfig is a helper function for client setup.
func createClientConfig(nodeURL string, contracts ContractAddresses, privateKey *ecdsa.PrivateKey, host string, peerAddress common.Address, peerHost string) client.ClientConfig {
	return client.ClientConfig{
		ClientConfig: perun.ClientConfig{
			PrivateKey:      privateKey,
			Host:            host,
			ETHNodeURL:      nodeURL,
			AdjudicatorAddr: contracts.AdjudicatorAddr,
			AssetHolderAddr: contracts.AssetHolderAddr,
			DialerTimeout:   1 * time.Second,
			PeerAddresses: []perun.PeerWithAddress{
				{
					Peer:    wallet.AsWalletAddr(peerAddress),
					Address: peerHost,
				},
			},
		},
		ChallengeDuration: collateralWithdrawalDelay / 2,
		AppAddress:        contracts.AppAddr,
		ContextTimeout:    defaultContextTimeout,
	}
}
