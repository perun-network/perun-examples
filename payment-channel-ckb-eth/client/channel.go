// Copyright 2025 PolyCrypt GmbH
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
	"math/big"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"time"
)

// PaymentChannel is a wrapper for a Perun channel for the payment use case.
type PaymentChannel struct {
	ch         *client.Channel
	currencies []channel.Asset
}

func (c PaymentChannel) State() *channel.State {
	return c.ch.State().Clone()
}

func (c *PaymentChannel) GetChannel() *client.Channel {
	return c.ch
}
func (c *PaymentChannel) GetChannelParams() *channel.Params {
	return c.ch.Params()
}

func (c *PaymentChannel) GetChannelState() *channel.State {
	return c.ch.State()
}

// newPaymentChannel creates a new payment channel.
func newPaymentChannel(ch *client.Channel, currencies []channel.Asset) *PaymentChannel {
	return &PaymentChannel{
		ch:         ch,
		currencies: currencies,
	}
}

// PerformSwap performs a swap by "swapping" the balances of the two
// participants for both assets.
func (c PaymentChannel) PerformSwap() {
	err := c.ch.Update(context.TODO(), func(state *channel.State) { // We use context.TODO to keep the code simple.
		// We simply swap the balances for the two assets.
		state.Balances = channel.Balances{
			{state.Balances[0][1], state.Balances[0][0]},
			{state.Balances[1][1], state.Balances[1][0]},
		}

		// Set the state to final because we do not expect any other updates
		// than this swap.
		state.IsFinal = true
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

// SendEthPayment sends a payment to the channel peer.
func (c PaymentChannel) SendEthPayment(amount float64) {
	// Transfer the given amount from us to peer.
	// Use UpdateBy to update the channel state.
	err := c.ch.Update(context.TODO(), func(state *channel.State) { // We use context.TODO to keep the code simple.
		ethAmount := EthToWei(big.NewFloat(amount))
		actor := c.ch.Idx()
		peer := 1 - actor
		state.Allocation.TransferBalance(actor, peer, c.currencies[0], ethAmount)
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

// SendCKBPayment sends a payment to the channel peer.
func (c PaymentChannel) SendCKBPayment(amount int64) {
	// Transfer the given amount from us to peer.
	// Use UpdateBy to update the channel state.
	err := c.ch.Update(context.TODO(), func(state *channel.State) {
		actor := c.ch.Idx()
		peer := 1 - actor
		shannonAmount := CKByteToShannon(big.NewFloat(float64(amount)))
		state.Allocation.TransferBalance(actor, peer, c.currencies[1], shannonAmount)

	})
	if err != nil {
		panic(err)
	}

}

// Settle settles the payment channel and withdraws the funds.
func (c PaymentChannel) Settle() {
	// Finalize the channel to enable fast settlement.
	if !c.ch.State().IsFinal {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
		defer cancel()
		err := c.ch.Update(ctx, func(state *channel.State) {
			state.IsFinal = true
		})
		if err != nil {
			panic(err)
		}
	}

	// Settle concludes the channel and withdraws the funds.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	err := c.ch.Settle(ctx, false)
	if err != nil {
		panic(err)
	}

	// Close frees up channel resources.
	c.ch.Close()
}
