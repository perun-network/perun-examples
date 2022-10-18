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
	"context"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

// SwapChannel is a wrapper for a Perun channel for the swap use case.
type SwapChannel struct {
	ch     *client.Channel
	assets [2]channel.Asset
}

// newSwapChannel creates a new channel for swaps.
func newSwapChannel(ch *client.Channel, currencies [2]channel.Asset) *SwapChannel {
	return &SwapChannel{
		ch:     ch,
		assets: currencies,
	}
}

// PerformSwap performs a swap by "swapping" the balances of the two
// participants for both assets.
func (c SwapChannel) PerformSwap() {
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

// Settle settles the channel and withdraws the funds.
func (c SwapChannel) Settle() {
	// Settle concludes the channel and withdraws the funds.
	err := c.ch.Settle(context.TODO(), false)
	if err != nil {
		panic(err)
	}

	// Close frees up channel resources.
	c.ch.Close()
}
