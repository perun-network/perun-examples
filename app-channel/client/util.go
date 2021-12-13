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
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ewallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

func makeStakeAllocation(asset common.Address, stake *big.Int) *channel.Allocation {
	return &channel.Allocation{
		Assets:   []channel.Asset{ewallet.AsWalletAddr(asset)},
		Balances: [][]*big.Int{{new(big.Int).Set(stake), new(big.Int).Set(stake)}},
	}
}

func (c *Client) PerunAddress() wallet.Address {
	return c.perunClient.Account.Address()
}

func (c *Client) Address() common.Address {
	return c.perunClient.Account.Account.Address
}

func (c *Client) challengeDurationInSeconds() uint64 {
	return uint64(c.challengeDuration.Seconds())
}

func (c *Client) defaultContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.contextTimeout)
}

func (c *Client) Logf(format string, v ...interface{}) {
	log.Printf("Client %v: %v", c.Address(), fmt.Sprintf(format, v...))
}

func (c *Client) OnChainBalance() (b *big.Int, err error) {
	ctx, cancel := c.defaultContextWithTimeout()
	defer cancel()
	return c.perunClient.EthClient.BalanceAt(ctx, c.Address(), nil)
}
