package client

import (
	"context"
	"math/big"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

// PaymentChannel is a wrapper for a Perun channel for the payment use case.
type PaymentChannel struct {
	ch       *client.Channel
	currency channel.Asset
}

// newPaymentChannel creates a new payment channel.
func newPaymentChannel(ch *client.Channel, currency channel.Asset) *PaymentChannel {
	return &PaymentChannel{
		ch:       ch,
		currency: currency,
	}
}

// SendPayment sends a payment to the channel peer.
func (c PaymentChannel) SendPayment(amount float64) {
	// Transfer the given amount from us to peer.
	// Use UpdateBy to update the channel state.
	err := c.ch.UpdateBy(context.TODO(), func(state *channel.State) error { // We use context.TODO to keep the code simple.
		plancks := DotToPlanck(big.NewFloat(amount))
		actor := c.ch.Idx()
		peer := 1 - actor
		state.Allocation.TransferBalance(actor, peer, c.currency, plancks)
		return nil
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

// Settle settles the payment channel and withdraws the funds.
func (c PaymentChannel) Settle() {
	// Finalize the channel to enable fast settlement.
	if !c.ch.State().IsFinal {
		err := c.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
			state.IsFinal = true
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	// Settle concludes the channel and withdraws the funds.
	err := c.ch.Settle(context.TODO(), false)
	if err != nil {
		panic(err)
	}

	// Close frees up channel resources.
	c.ch.Close()
}
