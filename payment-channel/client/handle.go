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
	"math/big"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/log"
)

const proposerIdx = 0 //TODO go-perun: expose channel.ProposerIdx

func (c *Client) HandleProposal(p client.ChannelProposal, r *client.ProposalResponder) {
	// Ensure that we got a ledger channel proposal.
	lcp, ok := p.(*client.LedgerChannelProposal)
	if !ok {
		fmt.Printf("Wrong channel proposal type: %T\n", p)
		r.Reject(context.TODO(), "Invalid proposal type.")
		return
	}

	// Check that we do not need to fund anything.
	zeroBal := big.NewInt(0)
	for _, bals := range lcp.FundingAgreement {
		for i, bal := range bals {
			if i != proposerIdx && bal.Cmp(zeroBal) != 0 {
				fmt.Printf("Expected funding balance 0, got %v\n", bal)
				r.Reject(context.TODO(), "Invalid funding agreement.")
				return
			}
		}
	}

	// Create a channel accept message and send it.
	accept := lcp.Accept(
		c.AccountAddress,         // The account we use in the channel.
		client.WithRandomNonce(), // Our share of the channel nonce.
	)
	ch, err := r.Accept(context.TODO(), accept)
	if err != nil {
		fmt.Printf("Error accepting channel proposal: %v\n", err)
		return
	}

	c.startWatcher(ch)
}

func (c *Client) HandleUpdate(cur *channel.State, next client.ChannelUpdate, r *client.UpdateResponder) {
	// We accept every update that increases our balance.

	// Ensure that the assets did not change.
	err := channel.AssetsAssertEqual(cur.Assets, next.State.Assets)
	if err != nil {
		r.Reject(context.TODO(), "Invalid assets.")
		return
	}

	// Ensure that our balance has increased.
	err = next.State.Balances.AssertGreaterOrEqual(cur.Balances)
	for i, bals := range next.State.Balances {
		for j, nextBal := range bals {
			curBal := cur.Balances[i][j]
			if i != proposerIdx && nextBal.Cmp(curBal) > 0 {
				r.Reject(context.TODO(), "Invalid balances.")
				return
			}
		}
	}

	// Send the acceptance message.
	err = r.Accept(context.TODO())
	if err != nil {
		fmt.Printf("Error accepting update: %v\n", err)
		return
	}
}

func (c *Client) startWatcher(ch *client.Channel) {
	go func() {
		err := ch.Watch(c)
		if err != nil {
			// Panic because if the watcher is not running, we are no longer
			// protected against registration of old states.
			panic("Watcher returned with error: %v", err)
		}
	}()
}

func (c *Client) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) { //TODO provide channel with event.
	c.Log("Received Adjudicator event: %T", e)
	if _, ok := e.(*channel.ConcludedEvent); ok {
		err := c.CloseChannel()
		if err != nil {
			log.Error(err)
		}
	}
}
