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

package ethereumUtil

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nervosnetwork/ckb-sdk-go/v2/indexer"
	ckbtypes "github.com/nervosnetwork/ckb-sdk-go/v2/types"
	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	"github.com/perun-network/perun-eth-backend/wallet/simple"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	"log"
	"math/big"
	"perun.network/perun-ckb-backend/wallet/address"
	"perun.network/perun-examples/payment-channel-ckb-eth/client"
)

// DeployContracts deploys the Perun smart contracts on the specified ledger.
func DeployContracts(nodeURL string, chainID uint64, privateKey string) (adj, ah common.Address) {
	k, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}
	w := swallet.NewWallet(k)
	ethClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		panic(err)
	}
	cb := ethchannel.NewContractBackend(
		ethClient,
		ethchannel.MakeChainID(big.NewInt(int64(chainID))),
		simple.NewTransactor(
			w,
			types.NewLondonSigner(big.NewInt(int64(chainID))),
		),
		1,
	)
	acc := accounts.Account{Address: crypto.PubkeyToAddress(k.PublicKey)}

	// Deploy adjudicator.
	adj, err = ethchannel.DeployAdjudicator(context.TODO(), cb, acc)
	if err != nil {
		panic(err)
	}

	// Deploy asset holder.
	ah, err = ethchannel.DeployETHAssetholder(context.TODO(), cb, adj, acc)
	if err != nil {
		panic(err)
	}

	return adj, ah
}

type balanceLogger struct {
	ethClient *ethclient.Client
}

// NewBalanceLogger creates a new balance logger for the specified ledger.
func NewBalanceLogger(chainURL string) balanceLogger {
	c, err := ethclient.Dial(chainURL)
	if err != nil {
		panic(err)
	}
	return balanceLogger{ethClient: c}
}

// LogBalances prints the balances of the specified accounts.
func (l balanceLogger) LogBalances(accounts ...common.Address) {
	bals := make([]*big.Float, len(accounts))
	for i, c := range accounts {
		bal, err := l.ethClient.BalanceAt(context.TODO(), c, nil)
		if err != nil {
			log.Fatal(err)
		}
		bals[i] = client.WeiToEth(bal)
	}
	log.Println("Client balances (ETH):", bals)
}

type CkbBalanceLogger struct {
	dialer indexer.Client
}

func NewCKBBalanceLogger(indexerURL string) CkbBalanceLogger {
	dialer, err := indexer.Dial(indexerURL)
	if err != nil {
		panic(err)
	}
	return CkbBalanceLogger{dialer: dialer}
}

func (l CkbBalanceLogger) LogBalances(participant *address.Participant) {
	ctx := context.Background()
	lockScript := &ckbtypes.Script{
		CodeHash: participant.UnlockScript.CodeHash, // secp256k1_blake160_sighash_all
		HashType: participant.UnlockScript.HashType,
		Args:     participant.UnlockScript.Args, // from public key
	}

	// Query CKB balance
	capacityResp, err := l.dialer.GetCellsCapacity(ctx, &indexer.SearchKey{
		Script:     lockScript,
		ScriptType: ckbtypes.ScriptTypeLock,
	})
	if err != nil {
		log.Fatalf("failed to get CKB balance: %v", err)
	}
	ckbAmount := capacityResp.Capacity
	log.Println("CKB balance", participant.PubKey, ckbAmount)
}
