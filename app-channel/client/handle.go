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
	"perun.network/perun-examples/app-channel/app"
)

func (c *Client) HandleProposal(proposal client.ChannelProposal, responder *client.ProposalResponder) {
	c.Lock()
	defer c.Unlock()

	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	ctx, cancel := c.defaultContextWithTimeout()
	defer cancel()

	_proposal, ok := proposal.(*client.LedgerChannelProposal)
	if !ok {
		err = responder.Reject(ctx, "accepting only ledger channel proposals")
		return
	} else if _, ok = _proposal.App.(*app.TicTacToeApp); !ok {
		err = responder.Reject(ctx, "accepting only collateralized channels")
		return
	}

	prop := &GameProposal{
		ch:       proposal,
		response: make(chan bool),
		result:   make(chan *ProposalResult),
	}
	c.gameProposals <- prop
	select {
	case <-ctx.Done():
		err = responder.Reject(ctx, "proposal response timeout")
		return
	case r := <-prop.response:
		if !r {
			err = responder.Reject(ctx, "game proposal rejected")
			return
		}
	}

	accept := _proposal.Accept(c.PerunAddress(), client.WithRandomNonce())
	ch, err := responder.Accept(ctx, accept)
	prop.result <- &ProposalResult{c.newGame(ch), err}
}

func (c *Client) HandleUpdate(cur *channel.State, update client.ChannelUpdate, responder *client.UpdateResponder) {
	c.Lock()
	defer c.Unlock()

	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	ctx, cancel := c.defaultContextWithTimeout()
	defer cancel()

	g, ok := c.games[update.State.ID]
	if !ok {
		err = responder.Reject(ctx, "unknown channel")
		return
	}

	_app, ok := g.ch.Params().App.(*app.TicTacToeApp)
	if !ok {
		err = responder.Reject(ctx, "invalid app")
		return
	}

	err = _app.ValidTransition(g.ch.Params(), g.state, update.State, update.ActorIdx)
	if err != nil {
		err = responder.Reject(ctx, fmt.Sprintf("invalid action: %v", err))
		return
	}

	err = responder.Accept(ctx)
	if err != nil {
		g.errs <- err
	}

	g.state = update.State.Clone()
}

func (c *Client) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	log.Info("Adjudicator event: %+v", e)
	switch e := e.(type) {
	case *channel.ConcludedEvent:
		c.Lock()
		defer c.Unlock()

		g, ok := c.games[e.ID()]
		if !ok {
			// We may see the concluded event twice.
			log.Warnf("channel %v not found", e.ID())
			return
		}

		g.ch.Close()
		delete(c.games, e.ID())
	}
}
