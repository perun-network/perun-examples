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

	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire"
)

// WalletAddress returns the wallet address of the client.
func (c *PaymentClient) WalletAddress() wallet.Address {
	return c.account.Address()
}

// WireAddress returns the wire address of the client.
func (c *PaymentClient) WireAddress() wire.Address {
	return c.waddress
}

// DotToPlanck converts a given amount in Dot to Planck.
func DotToPlanck(d *big.Float) *big.Int {
	plankFloat := new(big.Float).Mul(d, new(big.Float).SetFloat64(dot.PlanckPerDot))
	plank, _ := plankFloat.Int(nil)
	return plank
}

// PlanckToDot converts a given amount in Planck to Dot.
func PlanckToDot(d *big.Int) *dot.Dot {
	return dot.NewDotFromPlanck(d)
}
