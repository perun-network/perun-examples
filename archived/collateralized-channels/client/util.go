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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	ewallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
	"perun.network/perun-collateralized-channels/app"
)

func zeroBalance(asset common.Address) *channel.Allocation {
	return &channel.Allocation{
		Assets:   []channel.Asset{ewallet.AsWalletAddr(asset)},
		Balances: [][]*big.Int{{big.NewInt(0), big.NewInt(0)}},
	}
}

func (c *Client) PerunAddress() wallet.Address {
	return c.perunClient.Account.Address()
}

func (c *Client) Address() common.Address {
	return c.perunClient.Account.EthAddress()
}

func (c *Client) defaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), c.contextTimeout)
	return ctx
}

func (c *Client) settle(ch *Channel) (err error) {
	err = ch.UpdateBy(c.defaultContext(), func(s *channel.State) error {
		s.IsFinal = true
		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "final update")
	}

	err = ch.Register(c.defaultContext())
	if err != nil {
		return errors.WithMessage(err, "registering")
	}

	err = ch.Settle(c.defaultContext(), false)
	return errors.WithMessage(err, "settling channel")
}

func (ch *Channel) UpdateBy(ctx context.Context, update func(s *channel.State) error) error {
	err := ch.Channel.UpdateBy(ctx, update)
	if err != nil {
		return err
	}
	ch.state = ch.State().Clone()
	return nil
}

func (c *Client) Logf(format string, v ...interface{}) {
	log.Printf("Client %v: %v", c.Address(), fmt.Sprintf(format, v...))
}

func (c *Client) OnChainBalance() (b *big.Int, err error) {
	return c.perunClient.EthClient.BalanceAt(c.defaultContext(), c.Address(), nil)
}

func (c *Client) PeerCollateral() (b *big.Int, err error) {
	return c.peerCollateral(c.Address())
}

func (c *Client) peerCollateral(peer common.Address) (b *big.Int, err error) {
	b, err = c.assetHolder.Holdings(nil, calcPeerCollateralID(peer))
	return b, errors.WithMessage(err, "reading peer collateral")
}

func (c *Client) ChannelFunding(peer common.Address) (b *big.Int, err error) {
	ch, ok := c.channels[peer]
	if !ok {
		return nil, errors.New("channel not found")
	}
	b, err = c.assetHolder.Holdings(nil, calcChannelCollateralID(ch.ID(), c.Address()))
	return b, errors.WithMessage(err, "reading peer collateral")
}

func (c *Client) ChannelBalances() (balances map[common.Address]*big.Int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	balances = make(map[common.Address]*big.Int)
	for p, ch := range c.channels {
		bal, err := app.ChannelBalance(ch.Params().Parts, ch.State().Data, c.Address())
		if err != nil {
			return nil, errors.WithMessagef(err, "reading channel balance: %v", ch)
		}
		balances[p] = bal
	}
	return
}

func calcPeerCollateralID(peer common.Address) [32]byte {
	abiAddress, err := abi.NewType("address", "", nil)
	if err != nil {
		log.Panicf("failed to create abi address type: %v", err)
	}
	bytes, err := abi.Arguments{{Type: abiAddress}}.Pack(peer)
	if err != nil {
		log.Panicf("failed to encode peer address: %v", err)
	}
	return crypto.Keccak256Hash(bytes)
}

func calcChannelCollateralID(channelID channel.ID, peer common.Address) [32]byte {
	abiBytes32, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		log.Panicf("failed to create abi type: %v", err)
	}
	abiAddress, err := abi.NewType("address", "", nil)
	if err != nil {
		log.Panicf("failed to create abi type: %v", err)
	}
	bytes, err := abi.Arguments{{Type: abiBytes32}, {Type: abiAddress}}.Pack(channelID, peer)
	if err != nil {
		log.Panicf("failed to encode peer address: %v", err)
	}
	return crypto.Keccak256Hash(bytes)
}
