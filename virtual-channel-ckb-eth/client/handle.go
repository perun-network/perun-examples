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
	"time"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

// HandleProposal is the callback for incoming channel proposals.
func (c *PaymentClient) HandleProposal(p client.ChannelProposal, r *client.ProposalResponder) {
	log.Println("Received channel proposal")
	switch p := p.(type) {
	case *client.LedgerChannelProposalMsg:
		lcp, err := func() (*client.LedgerChannelProposalMsg, error) {
			// Ensure that we got a ledger channel proposal.
			lcp := p

			// Check that we have the correct number of participants.
			if lcp.NumPeers() != 2 {
				return nil, fmt.Errorf("invalid number of participants: %d", lcp.NumPeers())
			}

			// Check that the channel has the expected assets and funding balances.
			const assetIdx, peerIdx = 0, 1
			if err := channel.AssertAssetsEqual(lcp.InitBals.Assets, c.currency); err != nil {
				return nil, fmt.Errorf("invalid assets: %v", err)
			} else if lcp.FundingAgreement[assetIdx][peerIdx].Cmp(big.NewInt(0)) != 0 {
				// return nil, fmt.Errorf("invalid funding balance")
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
		return
	case *client.VirtualChannelProposalMsg:
		vcp, err := func() (*client.VirtualChannelProposalMsg, error) {
			vcp := p
			log.Printf("Received virtual channel proposal: %v", vcp)

			// Check that we have the correct number of participants.
			if vcp.NumPeers() != 2 {
				return nil, fmt.Errorf("invalid number of participants: %d", vcp.NumPeers())
			}

			// Check that the channel has the expected assets and funding balances.
			const assetIdx, peerIdx = 0, 1
			if err := channel.AssertAssetsEqual(vcp.InitBals.Assets, c.currency); err != nil {
				return nil, fmt.Errorf("invalid assets: %v", err)
			} else if vcp.FundingAgreement[assetIdx][peerIdx].Cmp(big.NewInt(0)) != 0 {
				// return nil, fmt.Errorf("invalid funding balance")
			}
			return vcp, nil
		}()
		if err != nil {
			log.Println("Rejecting proposal: ", err)
			r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
		}

		// Create a channel accept message and send it.
		accept := vcp.Accept(
			c.account, // The account we use in the channel.
		)

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
		defer cancel()
		log.Println("Accepting proposal: ", accept)
		ch, err := r.Accept(ctx, accept)
		if err != nil {
			fmt.Printf("Error accepting channel proposal: %v\n", err)
			return
		}

		// Start the on-chain event watcher. It automatically handles disputes.
		c.startWatching(ch)

		// Store channel.
		c.channels <- newPaymentChannel(ch, c.currency)
		return
	default:
		err := fmt.Errorf("invalid proposal type: %T", p)
		log.Println("Rejecting proposal: ", err)
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
	}
}

// HandleUpdate is the callback for incoming channel updates.
func (c *PaymentClient) HandleUpdate(cur *channel.State, next client.ChannelUpdate, r *client.UpdateResponder) {
	log.Println("Received channel update")
	// Send the acceptance message.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()
	err := r.Accept(ctx)
	if err != nil {
		panic(err)
	}
}

// HandleAdjudicatorEvent is the callback for smart contract events.
func (c *PaymentClient) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	log.Printf("Adjudicator event: type = %T, client = %v", e, c.account)
}
