package client

import (
	"context"
	"fmt"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/perun-examples/app-channel/app"
)

type Game struct {
	ch    *client.Channel
	state *channel.State // TODO:question - Do we really need to store the state? Fetching it when needed did not work
}

func newGame(ch *client.Channel) *Game {
	return &Game{ch: ch, state: ch.State().Clone()}
}

func (g *Game) String() string {
	return g.state.Data.(*app.TicTacToeAppData).String()
}

func (g *Game) Set(x, y int) {
	err := g.ch.UpdateBy(context.TODO(), func(state *channel.State) error {
		_app, ok := state.App.(*app.TicTacToeApp)
		if !ok {
			return fmt.Errorf("invalid app type: %T", _app)
		}

		return _app.Set(state, x, y, g.ch.Idx())
	})
	if err != nil {
		panic(err) // We panic on error to keep the code simple.
	}
	g.state = g.ch.State().Clone()
}

func (g *Game) Settle() {
	// Channel should be finalized through last ("winning") move. No need to set .isFinal here
	err := g.ch.Settle(context.TODO(), false)
	if err != nil {
		panic(err)
	}

	g.ch.Close()
}
