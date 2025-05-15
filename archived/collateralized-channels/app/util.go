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

package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ewallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

func (a *CollateralApp) ZeroBalance() channel.Data {
	return CollateralAppData{
		[][]*big.Int{{big.NewInt(0), big.NewInt(0)}},
	}
}

func IsZeroBalances(bals [][]*big.Int) bool {
	for _, bals2 := range bals {
		for _, bal := range bals2 {
			if bal.Sign() != 0 {
				return false
			}
		}
	}
	return true
}

func Transfer(peers []wallet.Address, s *channel.State, from common.Address, to common.Address, amount *big.Int) error {
	fromIdx, ok := peerIndex(peers, from)
	if !ok {
		return errors.Errorf("unknown address: %v", from)
	}
	toIdx, ok := peerIndex(peers, to)
	if !ok {
		return errors.Errorf("unknown address: %v", from)
	}

	d, ok := s.Data.(CollateralAppData)
	if !ok {
		return errors.New("invalid type")
	}

	// Update from balance and to balance
	d.balances[assetIdx][fromIdx] = new(big.Int).Sub(d.balances[assetIdx][fromIdx], amount)
	d.balances[assetIdx][toIdx] = new(big.Int).Add(d.balances[assetIdx][toIdx], amount)

	return nil
}

func ChannelBalance(peers []wallet.Address, data channel.Data, account common.Address) (*big.Int, error) {
	d, ok := data.(CollateralAppData)
	if !ok {
		return nil, errors.New("invalid type")
	}
	b, ok := d.Balance(peers, account)
	if !ok {
		return nil, errors.New("invalid peer")
	}
	return b, nil
}

func (d CollateralAppData) Balance(peers []wallet.Address, account common.Address) (*big.Int, bool) {
	peerIdx, ok := peerIndex(peers, account)
	if !ok {
		return nil, false
	}
	return new(big.Int).Set(d.balances[assetIdx][peerIdx]), true
}

func peerIndex(peers []wallet.Address, addr common.Address) (int, bool) {
	walletAddr := ewallet.AsWalletAddr(addr)
	for i, p := range peers {
		if p.Cmp(walletAddr) == 0 {
			return i, true
		}
	}
	return 0, false
}
