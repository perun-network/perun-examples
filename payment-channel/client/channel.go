package client

import (
	"context"
	"math/big"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

type Channel struct {
	ch *client.Channel
}

func newChannel(ch *client.Channel) *Channel {
	return &Channel{ch: ch}
}

//TODO document all exported functions.
func (c Channel) SendPayment(asset channel.Asset, amount uint64) {
	// Transfer the given amount from us to peer.
	// Use UpdateBy to update the channel state.
	err := c.ch.UpdateBy(context.TODO(), func(state *channel.State) error { //TODO mention that we always use context.TODO for simplicity.
		ethAmount := new(big.Int).SetUint64(amount)
		state.Allocation.TransferBalance(proposerIdx, receiverIdx, asset, ethAmount)
		return nil
	})
	if err != nil {
		panic(err) //TODO mention that we always panic on error for simplicity.
	}
}

func (c Channel) Close() {
	// Finalize the channel to enable fast settlement.
	err := c.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		state.IsFinal = true
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Settle concludes the channel and withdraws the funds.
	err = c.ch.Settle(context.TODO(), false)
	if err != nil {
		panic(err)
	}

	// Close frees up channel resources.
	c.ch.Close()
}
