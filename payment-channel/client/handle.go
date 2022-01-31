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
func (c *PaymentClient) HandleProposal(p client.ChannelProposal, r *client.ProposalResponder) {
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
		err := channel.AssetsAssertEqual(lcp.InitBals.Assets, []channel.Asset{c.currency})
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
	c.channels <- newPaymentChannel(ch, c.currency)
}

// HandleUpdate is the callback for incoming channel updates.
func (c *PaymentClient) HandleUpdate(cur *channel.State, next client.ChannelUpdate, r *client.UpdateResponder) {
	// We accept every update that increases our balance.
	err := func() error {
		err := channel.AssetsAssertEqual(cur.Assets, next.State.Assets) //TODO:go-perun move assets to parameters to disallow changing the assets until there is a use case for that?
		if err != nil {
			return fmt.Errorf("Invalid assets: %v", err)
		}

		//TODO:go-perun bug, machine.go: `validTransition` checks whether balances per asset index are preserved, but does not check whether assets are the same.
		curBal := cur.Allocation.Balance(receiverIdx, c.currency)
		nextBal := next.State.Allocation.Balance(receiverIdx, c.currency)
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
		panic(err)
	}
}

// HandleAdjudicatorEvent is the callback for smart contract events.
func (c *PaymentClient) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) { //TODO:go-perun provide channel with event. expose channel registry?
	log.Printf("Adjudicator event: type = %T, client = %v", e, c.account)
}
