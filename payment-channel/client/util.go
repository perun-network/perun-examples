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

package client

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"perun.network/go-perun/backend/ethereum/channel"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
)

func CreateContractBackend(
	nodeURL string,
	chainID uint64,
	w *swallet.Wallet,
) (channel.ContractBackend, error) {
	signer := types.NewEIP155Signer(new(big.Int).SetUint64(chainID))
	transactor := swallet.NewTransactor(w, signer) //TODO transactor should be spawnable from Wallet: Add method "NewTransactor"

	ethClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return channel.ContractBackend{}, err
	}

	return channel.NewContractBackend(ethClient, transactor, txFinalityDepth), nil
}

func (c *Client) Logf(format string, v ...interface{}) {
	log.Printf("%v: %s", c.Name, fmt.Sprintf(format, v...))
}
