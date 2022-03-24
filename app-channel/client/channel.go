package client

import (
	"context"
	"fmt"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/perun-examples/app-channel/app"
)

// TicTacToeChannel is a wrapper for a Perun channel for the Tic-tac-toe app use case.
type TicTacToeChannel struct {
	ch *client.Channel
}

// newTicTacToeChannel creates a new tic-tac-toe app channel.
func newTicTacToeChannel(ch *client.Channel) *TicTacToeChannel {
	return &TicTacToeChannel{ch: ch}
}

// Set sends a game move to the channel peer.
func (g *TicTacToeChannel) Set(x, y int) {
	err := g.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		app, ok := state.App.(*app.TicTacToeApp)
		if !ok {
			return fmt.Errorf("invalid app type: %T", app)
		}

		return app.Set(state, x, y, g.ch.Idx())
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

// ForceSet registers a game move on-chain.
func (g *TicTacToeChannel) ForceSet(x, y int) {
	err := g.ch.ForceUpdate(context.TODO(), func(state *channel.State) {
		err := func() error {
			app, ok := state.App.(*app.TicTacToeApp)
			if !ok {
				return fmt.Errorf("invalid app type: %T", app)
			}

			return app.Set(state, x, y, g.ch.Idx())
		}()
		if err != nil {
			panic(err)
		}
	})
	if err != nil {
		panic(err)
	}
}

// Settle settles the app channel and withdraws the funds.
func (g *TicTacToeChannel) Settle() {
	// Channel should be finalized through last ("winning") move.
	// No need to set `isFinal` here.
	err := g.ch.Settle(context.TODO(), false)
	if err != nil {
		panic(err)
	}

	// Cleanup.
	g.ch.Close()
}
