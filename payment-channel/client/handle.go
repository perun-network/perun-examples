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
	perunio "perun.network/go-perun/pkg/io"
)

func (c *Client) HandleProposal(p client.ChannelProposal, r *client.ProposalResponder) {
	// Ensure that we got a ledger channel proposal.
	lcp, ok := p.(*client.LedgerChannelProposal)
	if !ok {
		fmt.Printf("Wrong channel proposal type: %T\n", p)
		return
	}

	// Check that our balance is 0.
	perunio.EqualEncoding(c.asset, p.Base().InitBals.Assets) //TODO check correct assets. (one asset with correct address)
	initBal := lcp.FundingAgreement[assetIdx][partIdx]
	if initBal.Cmp(big.NewInt(0)) != 0 {
		fmt.Printf("Wrong initial balance: %v\n", initBal)
		return
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

	c.HandleNewChannel(ch)
}

func (c *Client) HandleUpdate(state *channel.State, update client.ChannelUpdate, responder *client.UpdateResponder) {
	// We will accept every update that increases our balance.
	if err := responder.Accept(context.TODO()); err != nil {
		fmt.Printf("%s: Could not accept update: %v\n", c.RoleAsString(), err)
	}
}

func (c *Client) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	fmt.Printf("%s: HandleAdjudicatorEvent\n", c.RoleAsString())
	if _, ok := e.(*channel.ConcludedEvent); ok && c.role == RoleAlice {
		err := c.CloseChannel()
		if err != nil {
			log.Error(err)
		}
	}
}

func (c *Client) HandleNewChannel(ch *client.Channel) {
	fmt.Printf("%s: HandleNewChannel with id 0x%x\n", c.RoleAsString(), ch.ID())
	c.Channel = ch
	// Start the on-chain watcher.
	go func() {
		err := ch.Watch(c)
		if err != nil {
			fmt.Printf("%s: Watcher returned with: %s", c.RoleAsString(), err)
		}
	}()
}
