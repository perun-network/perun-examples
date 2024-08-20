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
	"io"
	"math/big"

	"github.com/pkg/errors"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

const assetIdx = 0

// CollateralApp is a channel app enabling collateralized channels.
type CollateralApp struct {
	Addr wallet.Address
}

func NewCollateralApp(addr wallet.Address) *CollateralApp {
	return &CollateralApp{
		Addr: addr,
	}
}

// Def returns the app address.
func (a *CollateralApp) Def() wallet.Address {
	return a.Addr
}

// DecodeData decodes the channel data.
func (a *CollateralApp) DecodeData(r io.Reader) (channel.Data, error) {
	balances, err := readTupleInt256ArrayArray(r)
	if err != nil {
		return nil, errors.WithMessage(err, "reading (int256[][])")
	}
	return CollateralAppData{balances: balances}, nil
}

// CollateralAppData is the app data struct.
type CollateralAppData struct {
	balances [][]*big.Int
}

// Encode encodes app data onto an io.Writer.
func (d CollateralAppData) Encode(w io.Writer) (err error) {
	err = writeTupleInt256ArrayArray(w, d.balances)
	return errors.WithMessage(err, "writing (int256[][])")
}

// Clone returns a deep copy of the app data.
func (d CollateralAppData) Clone() channel.Data {
	balances := make([][]*big.Int, len(d.balances))
	for i := range balances {
		balances[i] = make([]*big.Int, len(d.balances[i]))
		for j := range balances[i] {
			balances[i][j] = new(big.Int).Set(d.balances[i][j])
		}
	}
	return CollateralAppData{balances: balances}
}

// ValidTransition checks that the data of the `to` state is of type Invoice.
func (a *CollateralApp) ValidTransition(_ *channel.Params, _, to *channel.State, _ channel.Index) error {
	if !IsZeroBalances(to.Balances) {
		return errors.New("must have zero balances")
	}
	return nil
}

// ValidInit checks that the initial state is valid.
func (a *CollateralApp) ValidInit(p *channel.Params, s *channel.State) error {
	d, ok := s.Data.(CollateralAppData)
	if !ok {
		return errors.New("failed to cast app data")
	} else if len(d.balances) != len(s.Assets) {
		return errors.New("invalid balance length")
	} else if !IsZeroBalances(s.Balances) {
		return errors.New("must have zero balances")
	}
	return nil
}
