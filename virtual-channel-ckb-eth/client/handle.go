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
	"fmt"
	"log"
	"math/big"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"time"
)

// HandleProposal is the callback for incoming channel proposals.
func (c *PaymentClient) HandleProposal(p client.ChannelProposal, r *client.ProposalResponder) {
	log.Println("Received channel proposal")
	lcp, err := func() (*client.LedgerChannelProposalMsg, error) {
		// Ensure that we got a ledger channel proposal.
		lcp, ok := p.(*client.LedgerChannelProposalMsg)
		if !ok {
			return nil, fmt.Errorf("invalid proposal type: %T", p)
		}

		// Check that we have the correct number of participants.
		if lcp.NumPeers() != 2 {
			return nil, fmt.Errorf("invalid number of participants: %d", lcp.NumPeers())
		}

		// Check that the channel has the expected assets and funding balances.
		const assetIdx, peerIdx = 0, 1
		if err := channel.AssertAssetsEqual(lcp.InitBals.Assets, c.currency); err != nil {
			return nil, fmt.Errorf("invalid assets: %v", err)
		} else if lcp.FundingAgreement[assetIdx][peerIdx].Cmp(big.NewInt(0)) != 0 {
			return nil, fmt.Errorf("invalid funding balance")
		}
		return lcp, nil
	}()
	if err != nil {
		log.Println("Rejecting proposal: ", err)
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
	}

	// Create a channel accept message and send it.
	accept := lcp.Accept(
		c.account,                // The account we use in the channel.
		client.WithRandomNonce(), // Our share of the channel nonce.
	)
	// ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	// defer cancel()
	log.Println("Accepting proposal: ", accept)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()
	ch, err := r.Accept(ctx, accept)
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
		err := channel.AssertAssetsEqual(cur.Assets, next.State.Assets)
		if err != nil {
			return fmt.Errorf("invalid assets: %v", err)
		}

		receiverIdx := 1 - next.ActorIdx // This works because we are in a two-party channel.
		curBal := cur.Allocation.Balance(receiverIdx, c.currency[0])
		nextBal := next.State.Allocation.Balance(receiverIdx, c.currency[0])
		if nextBal.Cmp(curBal) < 0 {
			return fmt.Errorf("invalid balance: %v", nextBal)
		}
		return nil
	}()
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
		defer cancel()
		r.Reject(ctx, err.Error()) //nolint:errcheck // It's OK if rejection fails.
	}

	// Send the acceptance message.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()
	err = r.Accept(ctx)
	if err != nil {
		panic(err)
	}
}

// HandleAdjudicatorEvent is the callback for smart contract events.
func (c *PaymentClient) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	log.Printf("Adjudicator event: type = %T, client = %v", e, c.account)
}
