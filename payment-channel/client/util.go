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
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"perun.network/go-perun/backend/ethereum/channel"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire"
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

// dummyAccount represents a wire account that does not support data signing.
type dummyAccount struct {
	addr wire.Address
}

// Address used by this account.
func (a dummyAccount) Address() wallet.Address {
	return a.addr
}

// SignData requests a signature from this account.
// It returns the signature or an error.
func (a dummyAccount) SignData(data []byte) ([]byte, error) {
	panic("unsupported")
}
