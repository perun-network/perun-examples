package client

import (
	"context"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/perun-examples/app-channel/app"
)

type Game struct {
	ch *client.Channel
}

func newGame(ch *client.Channel) *Game {
	return &Game{ch: ch}
}

func (g *Game) String() string {
	return g.ch.State().Clone().Data.(*app.TicTacToeAppData).String()
}

func (g *Game) Set(x, y int) {

	err := g.ch.UpdateBy(context.TODO(), func(s *channel.State) error {
		_app := s.App.(*app.TicTacToeApp)
		return _app.Set(s, x, y, g.ch.Idx())
	})

	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
}

func (g *Game) Settle() {
	// Finalize the channel to enable fast settlement.
	err := g.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		state.IsFinal = true // TODO:question - Can we finalize app channels like this?
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = g.ch.Settle(context.TODO(), false)
	if err != nil {
		panic(err)
	}

	g.ch.Close()
}
