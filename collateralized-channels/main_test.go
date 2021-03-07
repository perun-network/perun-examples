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
	"log"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/perun-collateralized-channels/client"
	"perun.network/perun-collateralized-channels/eth"
	"perun.network/perun-collateralized-channels/ganache"
	"perun.network/perun-collateralized-channels/perun"
)

// TestCollateralizedChannels is an end-to-end test of collateral channels.
func TestCollateralizedChannels(t *testing.T) {
	require := require.New(t)

	// Accounts
	accountFunding := []struct {
		PrivateKey string
		BalanceEth uint
	}{
		{"0x50b4713b4ba55b6fbcb826ae04e66c03a12fc62886a90ca57ab541959337e897", 10},  // Contract Deployer
		{"0x1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f", 100}, // Client1
		{"0xf63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e", 200}, // Client2
	}

	// Ganache config
	ganacheCmd := os.Getenv("GANACHE_CMD")
	if len(ganacheCmd) == 0 {
		ganacheCmd = "ganache-cli"
	}
	ganacheCfg := ganache.GanacheConfig{
		Cmd:         ganacheCmd,
		Host:        "127.0.0.1",
		Port:        8545,
		BlockTime:   1 * time.Second,
		Funding:     accountFunding,
		StartupTime: 3 * time.Second,
	}

	// More test parameters
	var (
		defaultContextTimeout     = 30 * time.Second
		collateralWithdrawalDelay = 10 * ganacheCfg.BlockTime
		hostClient1               = "127.0.0.1:8546"
		hostClient2               = "127.0.0.1:8547"
		collateralClient1         = eth.EthToWei(big.NewFloat(50))
		payment1Client1ToClient2  = eth.EthToWei(big.NewFloat(5))
		channelFundingClient1     = eth.EthToWei(big.NewFloat(25))
		payment2Client1ToClient2  = eth.EthToWei(big.NewFloat(10))
	)

	// Start ganache blockchain with prefunded accounts
	log.Print("Starting local blockchain...")
	ganache, err := ganache.StartGanacheWithPrefundedAccounts(ganacheCfg)
	require.NoError(err, "starting ganache")
	defer ganache.Shutdown()

	// Deploy contracts
	log.Print("Deploying contracts...")
	nodeURL := fmt.Sprintf("ws://%s:%d", ganacheCfg.Host, ganacheCfg.Port)
	deploymentKey := ganache.Accounts[0].PrivateKey
	contracts, err := deployContracts(nodeURL, deploymentKey, defaultContextTimeout, collateralWithdrawalDelay)
	require.NoError(err, "deploying contracts")

	// Helper function for client setup
	genClientDef := func(privateKey *ecdsa.PrivateKey, host string, peerAddress common.Address, peerHost string) client.ClientConfig {
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

	log.Print("Setting up clients...")
	// Setup Client1
	clientDef1 := genClientDef(
		ganache.Accounts[1].PrivateKey, hostClient1,
		ganache.Accounts[2].Address(), hostClient2,
	)
	paymentAcceptancePolicy1 := func(
		amount, collateral, funding, balance *big.Int,
		hasOverdrawn bool,
	) (ok bool) {
		return true
	}
	c1, err := client.SetupClient(
		clientDef1,
		paymentAcceptancePolicy1,
	)
	require.NoError(err, "Client1 setup")

	// Setup Client2
	clientDef2 := genClientDef(ganache.
		Accounts[2].PrivateKey, hostClient2,
		ganache.Accounts[1].Address(), hostClient1,
	)
	paymentAcceptancePolicy2 := func(
		amount, collateral, funding, balance *big.Int,
		hasOverdrawn bool,
	) (ok bool) {
		// We reject unfunded payments if they exceed 10% of the collateral.
		balanceFundingDiff := new(big.Int).Sub(funding, balance)
		collateral10percent := new(big.Int).Div(collateral, big.NewInt(10))
		if balanceFundingDiff.Sign() < 0 && balanceFundingDiff.Cmp(collateral10percent) < 0 {
			return false
		}

		// We accept all other payments.
		return true
	}
	c2, err := client.SetupClient(
		clientDef2,
		paymentAcceptancePolicy2,
	)
	require.NoError(err, "Client2 setup")

	e := &Environment{map[common.Address]string{
		c1.Address(): "Client1",
		c2.Address(): "Client2",
	}}
	e.logAccountBalance(c1, c2)
	log.Print("Setup done.")

	// Deposit Client1 collateral
	log.Printf("Client1: Depositing %v as collateral...", toEth(collateralClient1))
	err = c1.IncreaseCollateral(collateralClient1)
	require.NoError(err, "increasing Client1 collateral")
	e.logAccountBalance(c1)

	// Send payment from Client1 to Client2
	log.Printf("Client1: Sending %v to Client2...", toEth(payment1Client1ToClient2))
	err = c1.SendPayment(c2.Address(), payment1Client1ToClient2) // open unfunded channel, handle channel proposal, transfer amount, handle update
	require.NoError(err, "Client1 sending payment to Client2")
	e.logChannelBalances(c1, c2)

	// Client1 deposits channel funding
	log.Printf("Client1: Depositing %v as channel funding...", toEth(channelFundingClient1))
	err = c1.IncreaseChannelCollateral(c2.Address(), channelFundingClient1)
	require.NoError(err, "Client1 increasing channel funding")
	e.logAccountBalance(c1)
	e.logChannelBalances(c1)

	// Client1 sends another payment to Client2
	log.Printf("Client1: Sending %v to Client2...", toEth(payment2Client1ToClient2))
	err = c1.SendPayment(c2.Address(), payment2Client1ToClient2) // send another payment
	require.NoError(err, "Client1 sending another payment to Client2")
	e.logChannelBalances(c1, c2)

	// Client2 settles the channel and withdraws the received payments
	log.Print("Client2: Settling channel...")
	err = c2.Settle(c1.Address()) // c2 settles channel with c1
	require.NoError(err, "Client2 settling the channel")
	e.logAccountBalance(c2)

	log.Print("Done.")
}
