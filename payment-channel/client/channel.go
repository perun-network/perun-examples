package client

import (
	"context"
	"math/big"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

// PaymentChannel is a wrapper for a Perun channel for the payment use case.
type PaymentChannel struct {
	ch *client.Channel
}

// newPaymentChannel creates a new payment channel.
func newPaymentChannel(ch *client.Channel) *PaymentChannel {
	return &PaymentChannel{ch: ch}
}

// SendPayment sends a payment to the channel peer.
func (c PaymentChannel) SendPayment(asset channel.Asset, amount uint64) {
	// Transfer the given amount from us to peer.
	// Use UpdateBy to update the channel state.
	err := c.ch.UpdateBy(context.TODO(), func(state *channel.State) error { // We use context.TO DO to keep the code simple.
		ethAmount := new(big.Int).SetUint64(amount)
		state.Allocation.TransferBalance(proposerIdx, receiverIdx, asset, ethAmount)
		return nil
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

// Settle settles the payment channel and withdraws the funds.
func (c PaymentChannel) Settle() {
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
