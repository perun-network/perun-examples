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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/perun-eth-backend/wallet"
	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/log"
	"perun.network/perun-collateralized-channels/app"
)

func peerAddress(ch *client.Channel) common.Address {
	if len(ch.Params().Parts) != 2 {
		panic("not a two-party channel")
	}
	return wallet.AsEthAddr(ch.Params().Parts[1-int(ch.Idx())])
}

// HandleProposal handles channel proposals.
func (c *AppClient) HandleProposal(proposal client.ChannelProposal, responder *client.ProposalResponder) {
	log.Tracef("incoming channel proposal: %v", proposal)
	_proposal, ok := proposal.(*client.LedgerChannelProposalMsg)
	if !ok {
		responder.Reject(c.defaultContext(), "accepting only ledger channel proposals")
		return
	} else if _, ok = _proposal.App.(*app.CollateralApp); !ok {
		responder.Reject(c.defaultContext(), "accepting only collateralized channels")
		return
	}

	accept := _proposal.Accept(c.account, client.WithRandomNonce())
	ch, err := responder.Accept(c.defaultContext(), accept)
	if err != nil {
		log.Errorf("accepting channel proposal: %v", err)
		return
	}

	c.onNewChannel(ch)
}

// HandleUpdate handles channel updates.
func (c *AppClient) HandleUpdate(cur *channel.State, next client.ChannelUpdate, responder *client.UpdateResponder) {
	// Get channel and peer.
	peer, ch, ok := c.peerWithChannelForChannelID(cur.ID)
	if !ok {
		responder.Reject(c.defaultContext(), "unknown channel")
		return
	}

	// Get current and proposed balances.
	curData, ok := ch.state.Data.(*app.CollateralAppData)
	if !ok {
		responder.Reject(c.defaultContext(), "failed to parse current app data")
		return
	}
	propData, ok := next.State.Data.(*app.CollateralAppData)
	if !ok {
		responder.Reject(c.defaultContext(), "failed to parse proposed app data")
		return
	}
	curBal, _ := curData.Balance(ch.Params().Parts, c.WalletAddress())
	propBal, _ := propData.Balance(ch.Params().Parts, c.WalletAddress())
	// curPeerBal, _ := curData.Balance(ch.Params().Parts, peer)
	propPeerBal, _ := propData.Balance(ch.Params().Parts, peer)

	if propBal.Cmp(curBal) < 0 {
		// Reject updates that lower our balance.
		responder.Reject(c.defaultContext(), "invalid payment update: balance must increase")
		return
	} else if new(big.Int).Add(propBal, propPeerBal).Sign() != 0 {
		// Reject updates where the credit does not equal the debit.
		responder.Reject(c.defaultContext(), "invalid payment update: credit must equal debit")
		return
	}

	// Get peer collateral.
	peerCollateral, err := c.peerCollateral(peer)
	if err != nil {
		log.Infof("failed to get peer collateral: %v", err)
		responder.Reject(c.defaultContext(), "error getting peer collateral balance")
		return
	}

	// Get peer channel funding.
	peerFunding, err := c.ChannelFunding(peer)
	if err != nil {
		log.Infof("failed to get channel collateral: %v", err)
		responder.Reject(c.defaultContext(), "error getting channel collateral balance")
		return
	}

	// Get overdrawing history.
	hasOverdrawn, err := c.hasCollateralOverdrawnEvents(peer)
	if err != nil {
		log.Infof("failed to get insufficient collateral events: %v", err)
		responder.Reject(c.defaultContext(), "failed to assert clean settlement history")
		return
	}

	paymentAmount := new(big.Int).Sub(propBal, curBal)
	ok = c.updatePolicy(paymentAmount, peerCollateral, peerFunding, propPeerBal, hasOverdrawn)
	if !ok {
		responder.Reject(c.defaultContext(), "update reject by policy")
		return
	}

	err = responder.Accept(c.defaultContext())
	if err != nil {
		log.Panicf("Failed to accept channel update: %v", err)
	}
	ch.state = next.State.Clone()
}

// HandleAdjudicatorEvent handles adjudicator events.
func (c *AppClient) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	switch e := e.(type) {
	case *channel.ConcludedEvent:

		p, ch, ok := c.peerWithChannelForChannelID(e.ID())
		if !ok {
			log.Panicf("channel %v not found", e.ID())
		}

		ch.Close()
		delete(c.channels, p)
	}
}

func (c *AppClient) peerWithChannelForChannelID(cID channel.ID) (common.Address, *Channel, bool) {
	for p, ch := range c.channels {
		if ch.ID() == cID {
			return p, ch, true
		}
	}
	return common.Address{}, nil, false
}

func (c *AppClient) hasCollateralOverdrawnEvents(peer common.Address) (bool, error) {
	it, err := c.assetHolder.FilterCollateralOverdrawn(nil, []common.Address{peer})
	if err != nil {
		return false, errors.WithMessage(err, "creating filter")
	}

	for it.Next() {
		return true, nil
	}
	return false, nil
}
