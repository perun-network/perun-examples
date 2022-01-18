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
	"fmt"
	"math/big"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

func (c *Client) HandleProposal(p client.ChannelProposal, r *client.ProposalResponder) {
	lcp, err := func() (*client.LedgerChannelProposal, error) {
		// Ensure that we got a ledger channel proposal.
		lcp, ok := p.(*client.LedgerChannelProposal)
		if !ok {
			return nil, fmt.Errorf("Invalid proposal type: %T\n", p)
		}

		// Check that we have the correct number of participants.
		if lcp.NumPeers() != 2 { //TODO:go-perun rename NumPeers to NumParts, Peers to Participants anywhere where all parties are referred to
			return nil, fmt.Errorf("Invalid number of participants: %d", lcp.NumPeers())
		}

		// Check that the channel has the expected assets.
		err := channel.AssetsAssertEqual(lcp.InitBals.Assets, []channel.Asset{c.asset})
		if err != nil {
			return nil, fmt.Errorf("Invalid assets: %v\n", err)
		}

		// Check that we do not need to fund anything.
		zeroBal := big.NewInt(0)
		for _, bals := range lcp.FundingAgreement {
			bal := bals[receiverIdx]
			if bal.Cmp(zeroBal) != 0 {
				return nil, fmt.Errorf("Invalid funding balance: %v", bal)
			}
		}
		return lcp, nil
	}()
	if err != nil {
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
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

	// Store channel.
	c.channelsMtx.Lock()
	c.channels[ch.ID()] = newChannel(ch)
	c.channelsMtx.Unlock()

	// Start the on-chain event watcher. It automatically handles disputes.
	go func() {
		err := ch.Watch(c)
		if err != nil {
			// Panic because if the watcher is not running, we are no longer
			// protected against registration of old states.
			panic(fmt.Sprintf("Watcher returned with error: %v", err))
		}
	}()
}

func (c *Client) HandleUpdate(cur *channel.State, next client.ChannelUpdate, r *client.UpdateResponder) {
	// We accept every update that increases our balance.
	err := func() error {
		err := channel.AssetsAssertEqual(cur.Assets, next.State.Assets) //TODO:go-perun move assets to parameters to disallow changing the assets until there is a use case for that?
		if err != nil {
			return fmt.Errorf("Invalid assets: %v", err)
		}

		//TODO comment: go-perun ensures that total balance stays the same. //TODO:go-perun bug, machine.go:validTransition does only check balances, but not assets.
		curBal := cur.Allocation.Balance(receiverIdx, c.asset)
		nextBal := next.State.Allocation.Balance(receiverIdx, c.asset)
		if nextBal.Cmp(curBal) < 0 {
			return fmt.Errorf("Invalid balance: %v", nextBal)
		}
		return nil
	}()
	if err != nil {
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
	}

	// Send the acceptance message.
	err = r.Accept(context.TODO())
	if err != nil {
		c.Logf("Error accepting update: %v\n", err)
		return
	}
}

func (c *Client) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) { //TODO:go-perun provide channel with event. expose channel registry?
	c.Logf("Received Adjudicator event: %T", e)

	switch e := e.(type) {
	case *channel.ConcludedEvent:
		c.channelsMtx.RLock()
		ch := c.channels[e.ID()]
		c.channelsMtx.RUnlock()

		err := ch.ch.Settle(context.TODO(), false)
		if err != nil {
			panic(err)
		}
	}
}
