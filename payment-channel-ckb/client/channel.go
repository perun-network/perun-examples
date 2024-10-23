package client

import (
	"context"
	"fmt"
	"math/big"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

type PaymentChannel struct {
	ch     *client.Channel
	assets []channel.Asset
}

// newPaymentChannel creates a new payment channel.
func newPaymentChannel(ch *client.Channel, assets []channel.Asset) *PaymentChannel {
	return &PaymentChannel{
		ch:     ch,
		assets: assets,
	}
}

func (c PaymentChannel) State() *channel.State {
	return c.ch.State().Clone()
}

func (c PaymentChannel) SendPayment(amounts map[channel.Asset]float64) {
	// Transfer the given amount from us to peer.
	// Use UpdateBy to update the channel state.
	err := c.ch.Update(context.TODO(), func(state *channel.State) {
		actor := c.ch.Idx()
		peer := 1 - actor
		//fmt.Println("Send payment handler called")
		for a, amount := range amounts {
			fmt.Println(a)
			fmt.Println(amount)
			if amount < 0 {
				continue
			}
			/*
				switch a := a.(type) {
				case *asset.Asset:
					if a.IsCKBytes {
						fmt.Println("inside condition isCKBytes")
						shannonAmount := CKByteToShannon(big.NewFloat(amount))
						state.Allocation.TransferBalance(actor, peer, a, shannonAmount)
					} else {
						fmt.Println("inside if conditional !isCKBytes")
						intAmount := new(big.Int).SetUint64(uint64(amount))
						state.Allocation.TransferBalance(actor, peer, a, intAmount)
					}
				}
			*/

			shannonAmount := CKByteToShannon(big.NewFloat(amount))
			state.Allocation.TransferBalance(actor, peer, a, shannonAmount)

		}

	})
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

// Settle settles the payment channel and withdraws the funds.
func (c PaymentChannel) Settle() {
	// Finalize the channel to enable fast settlement.
	if !c.ch.State().IsFinal {
		err := c.ch.Update(context.TODO(), func(state *channel.State) {
			state.IsFinal = true
		})
		if err != nil {
			panic(err)
		}
	}

	// Settle concludes the channel and withdraws the funds.
	err := c.ch.Settle(context.TODO())
	if err != nil {
		panic(err)
	}

	// Close frees up channel resources.
	c.ch.Close()
}
