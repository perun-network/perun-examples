// Copyright 2025 PolyCrypt GmbH
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
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

// HandleProposal is the callback for incoming channel proposals.
func (c *PaymentClient) HandleProposal(p client.ChannelProposal, r *client.ProposalResponder) {
	lcp, err := func() (*client.LedgerChannelProposalMsg, error) {
		// Ensure that we got a ledger channel proposal.
		lcp, ok := p.(*client.LedgerChannelProposalMsg)
		if !ok {
			return nil, errors.New("invalid proposal type: expected *client.LedgerChannelProposalMsg")
		}

		// Check that we have the correct number of participants.
		if lcp.NumPeers() != 2 {
			return nil, fmt.Errorf("invalid number of participants: %d", lcp.NumPeers())
		}

		// Check that the channel has the expected assets and funding balances.
		const assetIdx, clientIdx, peerIdx = 0, 0, 1

		pAssets := make([]channel.Asset, len(c.currencies))
		for i, asset := range c.currencies {
			pAssets[i] = channel.Asset(asset)
		}

		if err := channel.AssertAssetsEqual(lcp.InitBals.Assets, pAssets); err != nil {
			return nil, fmt.Errorf("invalid assets: %v", err)
		} else if lcp.FundingAgreement[assetIdx][clientIdx].Cmp(big.NewInt(0)) == 0 && lcp.FundingAgreement[assetIdx][peerIdx].Cmp(big.NewInt(0)) == 0 {
			return nil, errors.New("invalid funding balance: both client and peer cannot have zero funding")
		}
		return lcp, nil
	}()
	if err != nil {
		errReject := r.Reject(context.TODO(), err.Error())
		if errReject != nil {
			// Log the error or take other action as needed
			log.Printf("error rejecting proposal: %v\n", errReject)
		}
	}

	// Create a channel accept message and send it.
	accept := lcp.Accept(
		c.account,                // The account we use in the channel.
		client.WithRandomNonce(), // Our share of the channel nonce.
	)
	ch, err := r.Accept(context.TODO(), accept)
	if err != nil {
		log.Printf("error accepting channel proposal: %v\n", err)
		return
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)

	// Store channel.
	c.channels <- newPaymentChannel(ch, c.currencies)

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
		// curBal0 := cur.Allocation.Balance(receiverIdx, c.currencies[0])
		// nextBal0 := next.State.Allocation.Balance(receiverIdx, c.currencies[0])
		// if nextBal0.Cmp(curBal0) < 0 {
		// 	return fmt.Errorf("Invalid balance: %v", nextBal0)
		// }
		for _, currency := range c.currencies {
			curBal := cur.Allocation.Balance(receiverIdx, currency)
			nextBal := next.State.Allocation.Balance(receiverIdx, currency)
			if nextBal.Cmp(curBal) < 0 {
				return fmt.Errorf("invalid balance for asset %v: %v", currency, nextBal)
			}
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
func (c *PaymentClient) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	log.Printf("Adjudicator event: type = %T, client = %v", e, c.account)
}

func (c *PaymentClient) GetChannel() (*PaymentChannel, error) {
	select {
	case channel := <-c.channels:
		c.channels <- channel // Put the channel back into the channels channel.
		return channel, nil
	case <-time.After(time.Second): // Set a timeout duration (e.g., 1 second).
		return nil, errors.New("no channel available")
	}
}
