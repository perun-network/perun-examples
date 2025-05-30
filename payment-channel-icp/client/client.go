// Copyright 2025 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	"math/big"
	"sync"

	pchannel "perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wire"
	"perun.network/perun-icp-backend/channel"
	chanconn "perun.network/perun-icp-backend/channel/connector"
	"perun.network/perun-icp-backend/channel/connector/icperun"
	icwallet "perun.network/perun-icp-backend/wallet"
)

// PaymentClient is a payment channel client.
type PaymentClient struct {
	perunClient   *client.Client
	account       *icwallet.Account
	currency      pchannel.Asset
	channels      chan *PaymentChannel
	Channel       *PaymentChannel
	ICConn        *chanconn.Connector
	observerMutex sync.Mutex
	balanceMutex  sync.Mutex
	Name          string
	wAddr         wire.Address
	balance       *big.Int
}

// startWatching starts the dispute watcher for the specified channel.
func (c *PaymentClient) startWatching(ch *client.Channel) {
	go func() {
		err := ch.Watch(c)
		if err != nil {
			fmt.Printf("Watcher returned with error: %v", err)
		}
	}()
}

// OpenChannel opens a new channel with the specified peer and funding.
func (c *PaymentClient) OpenChannel(peer wire.Address, amount float64) { //*PaymentChannel
	// We define the channel participants. The proposer has always index 0. Here
	// we use the on-chain addresses as off-chain addresses, but we could also
	// use different ones.

	participants := []wire.Address{c.WireAddress(), peer}

	// We create an initial allocation which defines the starting balances.
	initBal := big.NewInt(int64(amount))

	initAlloc := pchannel.NewAllocation(2, channel.Asset)
	initAlloc.SetAssetBalances(channel.Asset, []pchannel.Bal{
		initBal, // Our initial balance.
		initBal, // Peer's initial balance.
	})

	// Prepare the channel proposal by defining the channel parameters.
	challengeDuration := uint64(10) // On-chain challenge duration in seconds.
	proposal, err := client.NewLedgerChannelProposal(
		challengeDuration,
		c.account.Address(),
		initAlloc,
		participants,
	)
	if err != nil {
		panic(err)
	}

	// Send the proposal.
	ch, err := c.perunClient.ProposeChannel(context.TODO(), proposal)
	if err != nil {
		panic(err)
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)
	c.Channel = newPaymentChannel(ch, c.currency)
}

func (p *PaymentClient) WireAddress() wire.Address {
	return p.wAddr
}

// AcceptedChannel returns the next accepted channel.
func (c *PaymentClient) AcceptedChannel() *PaymentChannel {
	c.Channel = <-c.channels
	return c.Channel

}

// Shutdown gracefully shuts down the client.
func (c *PaymentClient) Shutdown() {
	c.perunClient.Close()
}

func (c *PaymentClient) GetChannelBalance() (*big.Int, error) {
	if c.Channel == nil {
		return big.NewInt(0), nil
	}

	chanParams := c.Channel.GetChannelParams()
	cid := chanParams.ID()
	addr := chanParams.Parts[c.Channel.ch.Idx()]
	addrBytes, err := addr.MarshalBinary()
	if err != nil {
		panic(err)
	}

	queryBalArgs := icperun.Funding{
		Channel:     cid,
		Participant: addrBytes,
	}

	balNat, err := c.ICConn.PerunAgent.QueryHoldings(queryBalArgs)
	if err != nil {
		panic(err)
	}

	if (*balNat) == nil {
		return big.NewInt(0), nil
	}

	return (*balNat).BigInt(), nil
}

func (p *PaymentClient) DisplayAddress() string {
	addr := p.account.Address().String()

	return addr
}

func (p *PaymentClient) DisplayName() string {
	return p.Name
}

func (p *PaymentClient) HasOpenChannel() bool {
	return p.Channel != nil
}
