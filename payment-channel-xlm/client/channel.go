package client

import (
	"context"
	"math/big"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

type PaymentChannel struct {
	ch         *client.Channel
	currencies []channel.Asset
}

func (c *PaymentChannel) GetChannel() *client.Channel {
	return c.ch
}
func (c *PaymentChannel) GetChannelParams() *channel.Params {
	return c.ch.Params()
}

func (c *PaymentChannel) GetChannelState() *channel.State {
	return c.ch.State()
}

func newPaymentChannel(ch *client.Channel, currencies []channel.Asset) *PaymentChannel {
	return &PaymentChannel{
		ch:         ch,
		currencies: currencies,
	}
}

// SendPayment sends a payment to the channel peer.
func (c PaymentChannel) SendPayment(amount int64, assetIdx int) {
	// Transfer the given amount from us to peer.
	// Use UpdateBy to update the channel state.
	err := c.ch.Update(context.TODO(), func(state *channel.State) {
		icp := big.NewInt(amount)
		actor := c.ch.Idx()
		peer := 1 - actor
		state.Allocation.TransferBalance(actor, peer, c.currencies[assetIdx], icp)
	})
	if err != nil {
		panic(err)
	}
}

// PerformSwap performs a swap by "swapping" the balances of the two
// participants for both assets.
func (c PaymentChannel) PerformSwap() {
	err := c.ch.Update(context.TODO(), func(state *channel.State) { // We use context.TODO to keep the code simple.
		// We simply swap the balances for the two assets.
		state.Balances = channel.Balances{
			{state.Balances[0][1], state.Balances[0][0]},
			{state.Balances[1][1], state.Balances[1][0]},
		}

		// Set the state to final because we do not expect any other updates
		// than this swap.
		state.IsFinal = true
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

// Settle settles the payment channel and withdraws the funds.
func (c PaymentChannel) Settle() {
	// If the channel is not finalized: Finalize the channel to enable fast settlement.

	if !c.ch.State().IsFinal {
		err := c.ch.Update(context.TODO(), func(state *channel.State) {
			state.IsFinal = true
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
