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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"perun.network/go-perun/wallet"
)

func (c *Client) PerunAddress() wallet.Address {
	return ethwallet.
}

func (c *Client) Address() common.Address {
	return c.PerunClient.Account.Account.Address
}

func (c *Client) Logf(format string, v ...interface{}) {
	fmt.Printf("Client %v: %v", c.Address(), fmt.Sprintf(format, v...))
}

func (c *Client) OnChainBalance() (b *big.Int, err error) {
	ctx, cancel := c.defaultContextWithTimeout()
	defer cancel()
	return c.PerunClient.EthClient.BalanceAt(ctx, c.Address(), nil)
}

func (c *Client) RoleAsString() (name string) {
	if c.role == RoleAlice {
		return "Alice"
	} else {
		return "Bob"
	}
}

func (c *Client) PeerRoleAsString() (name string) {
	if 1-c.role == RoleAlice {
		return "Alice"
	} else {
		return "Bob"
	}
}
