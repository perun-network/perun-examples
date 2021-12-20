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
	"fmt"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/log"
)

func (c *Client) HandleProposal(proposal client.ChannelProposal, responder *client.ProposalResponder) {
	// Check that we got a ledger channel proposal.
	_proposal, ok := proposal.(*client.LedgerChannelProposal)
	if !ok {
		fmt.Printf("%s: Received a proposal that was not for a ledger channel.", c.RoleAsString())
		return
	}
	fmt.Printf("%s: Received channel proposal\n", c.RoleAsString())

	ctx, cancel := c.defaultContextWithTimeout()
	defer cancel()

	// Create a channel accept message and send it.
	accept := _proposal.Accept(c.PerunAddress(), client.WithRandomNonce())
	ch, err := responder.Accept(ctx, accept)
	c.HandleNewChannel(ch) // TODO: 1/2 Check with MG why this is needed here (and not needed in App Channel example)

	if err != nil {
		fmt.Printf("%s: Accepting channel: %w\n", c.RoleAsString(), err)
	} else {
		fmt.Printf("%s: Accepted channel with id 0x%x\n", c.RoleAsString(), ch.ID())
	}
}

func (c *Client) HandleUpdate(state *channel.State, update client.ChannelUpdate, responder *client.UpdateResponder) {
	fmt.Printf("%s: HandleUpdate\n", c.RoleAsString())
	ctx, cancel := c.defaultContextWithTimeout()
	defer cancel()

	// We will accept every update
	if err := responder.Accept(ctx); err != nil {
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
