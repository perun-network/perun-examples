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

package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	"github.com/perun-network/perun-eth-backend/wire"
)

// CreateContractBackend creates a new contract backend.
func CreateContractBackend(
	nodeURL string,
	chainID *big.Int,
	w *swallet.Wallet,
) (ethchannel.ContractBackend, error) {
	signer := types.LatestSignerForChainID(chainID)
	transactor := swallet.NewTransactor(w, signer)

	ethClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return ethchannel.ContractBackend{}, err
	}

	return ethchannel.NewContractBackend(ethClient, ethchannel.MakeChainID(chainID), transactor, txFinalityDepth), nil
}

// WalletAddress returns the wallet address of the client.
func (c *SwapClient) WalletAddress() common.Address {
	return common.Address(*c.account.(*ethwallet.Address))
}

// WireAddress returns the wire address of the client.
func (c *SwapClient) WireAddress() *wire.Address {
	return &wire.Address{Address: ethwallet.AsWalletAddr(c.WalletAddress())}
}
