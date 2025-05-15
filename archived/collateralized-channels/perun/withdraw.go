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

package perun

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	pwallet "perun.network/go-perun/wallet"
	"perun.network/perun-collateralized-channels/contracts/generated/collateralAssetHolderETH"
)

func NewWithdrawalAuth(channelID channel.ID, acc pwallet.Account, receiver pwallet.Address, amount *big.Int) (collateralAssetHolderETH.AssetHolderWithdrawalAuth, []byte, error) {
	auth := collateralAssetHolderETH.AssetHolderWithdrawalAuth{
		ChannelID:   channelID,
		Participant: wallet.AsEthAddr(acc.Address()),
		Receiver:    wallet.AsEthAddr(receiver),
		Amount:      amount,
	}
	enc, err := encodeAssetHolderWithdrawalAuth(auth)
	if err != nil {
		return collateralAssetHolderETH.AssetHolderWithdrawalAuth{}, nil, errors.WithMessage(err, "encoding withdrawal auth")
	}

	sig, err := acc.SignData(enc)
	return auth, sig, errors.WithMessage(err, "sign data")
}

func encodeAssetHolderWithdrawalAuth(auth collateralAssetHolderETH.AssetHolderWithdrawalAuth) ([]byte, error) {
	var (
		abiUint256, _ = abi.NewType("uint256", "", nil)
		abiAddress, _ = abi.NewType("address", "", nil)
		abiBytes32, _ = abi.NewType("bytes32", "", nil)
	)

	// encodeAssetHolderWithdrawalAuth encodes the AssetHolderWithdrawalAuth as with abi.encode() in the smart contracts.
	args := abi.Arguments{
		{Type: abiBytes32},
		{Type: abiAddress},
		{Type: abiAddress},
		{Type: abiUint256},
	}
	enc, err := args.Pack(
		auth.ChannelID,
		auth.Participant,
		auth.Receiver,
		auth.Amount,
	)
	return enc, errors.WithStack(err)
}
