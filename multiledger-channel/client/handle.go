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
	"log"
	"math/big"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

// HandleProposal is the callback for incoming channel proposals.
func (c *SwapClient) HandleProposal(p client.ChannelProposal, r *client.ProposalResponder) {
	lcp, err := func() (*client.LedgerChannelProposalMsg, error) {
		// Ensure that we got a ledger channel proposal.
		lcp, ok := p.(*client.LedgerChannelProposalMsg)
		if !ok {
			return nil, fmt.Errorf("Invalid proposal type: %T\n", p)
		}

		// Check that we have the correct number of participants.
		if lcp.NumPeers() != 2 {
			return nil, fmt.Errorf("Invalid number of participants: %d", lcp.NumPeers())
		}

		// Check that the channel has the expected assets and funding balances.
		const assetIdx, peerIdx = 0, 1
		if err := channel.AssertAssetsEqual(lcp.InitBals.Assets, c.currencies[:]); err != nil {
			return nil, fmt.Errorf("Invalid assets: %v\n", err)
		} else if lcp.FundingAgreement[assetIdx][peerIdx].Cmp(big.NewInt(0)) != 0 {
			return nil, fmt.Errorf("Invalid funding balance")
		}
		return lcp, nil
	}()
	if err != nil {
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
	}

	// Create a channel accept message and send it.
	accept := lcp.Accept(
		c.account,                // The account we use in the channel.
		client.WithRandomNonce(), // Our share of the channel nonce.
	)
	ch, err := r.Accept(context.TODO(), accept)
	if err != nil {
		fmt.Printf("Error accepting channel proposal: %v\n", err)
		return
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)

	// Store channel.
	c.channels <- newSwapChannel(ch, c.currencies)
}

// HandleUpdate is the callback for incoming channel updates.
func (c *SwapClient) HandleUpdate(cur *channel.State, next client.ChannelUpdate, r *client.UpdateResponder) {
	// We accept every update that increases our balance.
	err := func() error {
		err := channel.AssertAssetsEqual(cur.Assets, next.State.Assets)
		if err != nil {
			return fmt.Errorf("Invalid assets: %v", err)
		}

		// Check the swap of the two currencies.
		for _, currency := range c.currencies {
			myIdx := 1 - next.ActorIdx // This works because we are in a two-party channel.
			// Get our and the peer's current balance.
			ourCurBal := cur.Balance(myIdx, currency)
			peerCurBal := cur.Balance(next.ActorIdx, currency)
			// Get our and the peer's balance of the next proposed state.
			ourNextBal := next.State.Balance(myIdx, currency)
			peerNextBal := next.State.Balance(next.ActorIdx, currency)

			// Check that the balances are "swapped".
			if !(ourCurBal.Cmp(peerNextBal) == 0 && ourNextBal.Cmp(peerCurBal) == 0) {
				return fmt.Errorf("invalid swap of balances")
			}
		}

		if !next.State.IsFinal {
			return fmt.Errorf("swap update must be final")
		}

		return nil
	}()
	if err != nil {
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
	}

	// Send the acceptance message.
	err = r.Accept(context.TODO())
	if err != nil {
		panic(err)
	}
}

// HandleAdjudicatorEvent is the callback for smart contract events.
func (c *SwapClient) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	log.Printf("Adjudicator event: type = %T, client = %v", e, c.account)
}
