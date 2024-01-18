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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
	"perun.network/perun-collateralized-channels/app"
)

func zeroBalance(asset channel.Asset) *channel.Allocation {
	initAlloc := channel.NewAllocation(2, asset)
	initAlloc.SetAssetBalances(asset, []channel.Bal{
		big.NewInt(0), // Our initial balance.
		big.NewInt(0), // Peer's initial balance.
	})
	return initAlloc
}

// WalletAddress returns the wallet address of the client.
func (c *AppClient) WalletAddress() common.Address {
	return common.Address(*c.account.(*ethwallet.Address))
}

func (c *AppClient) defaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), c.contextTimeout)
	return ctx
}

func (c *AppClient) settle(ch *Channel) (err error) {
	// Finalize the channel to enable fast settlement.
	if !ch.State().IsFinal {
		err := ch.Update(context.TODO(), func(state *channel.State) {
			state.IsFinal = true
		})
		if err != nil {
			return errors.WithMessage(err, "final update")
		}
	}

	err = ch.Settle(c.defaultContext(), false)
	if err != nil {
		return errors.WithMessage(err, "settling channel")

	}

	return ch.Close()
}

// Update updates the channel state and the internal state.
func (ch *Channel) Update(ctx context.Context, update func(s *channel.State)) error {
	err := ch.Channel.Update(ctx, update)
	if err != nil {
		return err
	}
	ch.state = ch.State().Clone()
	return nil
}

// Logf logs a message with the client's address.
func (c *AppClient) Logf(format string, v ...interface{}) {
	log.Printf("Client %v: %v", c.WalletAddress(), fmt.Sprintf(format, v...))
}

// OnChainBalance returns the on-chain balance of the client.
func (c *AppClient) OnChainBalance() (b *big.Int, err error) {
	return c.ethClient.BalanceAt(c.defaultContext(), c.WalletAddress(), nil)
}

// PeerCollateral returns the collateral of the peer.
func (c *AppClient) PeerCollateral() (b *big.Int, err error) {
	return c.peerCollateral(c.WalletAddress())
}

func (c *AppClient) peerCollateral(peer common.Address) (b *big.Int, err error) {
	b, err = c.assetHolder.Holdings(nil, calcPeerCollateralID(peer))
	return b, errors.WithMessage(err, "reading peer collateral")
}

// ChannelFunding returns the channel funding of the peer.
func (c *AppClient) ChannelFunding(peer common.Address) (b *big.Int, err error) {
	ch, ok := c.channels[peer]
	if !ok {
		return nil, errors.New("channel not found")
	}
	b, err = c.assetHolder.Holdings(nil, calcChannelCollateralID(ch.ID(), c.WalletAddress()))
	return b, errors.WithMessage(err, "reading peer collateral")
}

// ChannelBalances returns the channel balances of the client.
func (c *AppClient) ChannelBalances() (balances map[common.Address]*big.Int, err error) {
	balances = make(map[common.Address]*big.Int)
	for p, ch := range c.channels {
		bal, err := app.ChannelBalance(ch.Params().Parts, ch.State().Data, c.WalletAddress())
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

// CreateContractBackend creates a new contract backend.
func CreateContractBackend(
	nodeURL string,
	chainID uint64,
	wallet *swallet.Wallet,
) (*ethclient.Client, ethchannel.ContractBackend, *swallet.Transactor, error) {
	signer := types.LatestSignerForChainID(new(big.Int).SetUint64(chainID))
	transactor := swallet.NewTransactor(wallet, signer)
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, ethchannel.ContractBackend{}, nil, nil
	}

	return client, ethchannel.NewContractBackend(client, ethchannel.MakeChainID(big.NewInt(int64(chainID))), transactor, txFinalityDepth), transactor, nil
}
