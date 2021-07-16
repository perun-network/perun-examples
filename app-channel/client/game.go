package client

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/perun-examples/app-channel/app"
)

type (
	Game struct {
		sync.Mutex
		c     *Client
		ch    *client.Channel
		state *channel.State
		errs  chan error
	}

	GameProposal struct {
		ch       client.ChannelProposal
		response chan bool
		result   chan *ProposalResult
	}

	ProposalResult struct {
		g   *Game
		err error
	}
)

func (p *GameProposal) Accept() (*Game, error) {
	p.response <- true
	r := <-p.result
	return r.g, r.err
}

func (g *Game) String() string {
	return g.state.Data.(*app.TicTacToeAppData).String()
}

func (g *Game) Set(x, y int) error {
	g.Lock()
	defer g.Unlock()

	ctx, cancel := g.c.defaultContextWithTimeout()
	defer cancel()

	err := g.ch.UpdateBy(ctx, func(s *channel.State) error {
		_app, ok := s.App.(*app.TicTacToeApp)
		if !ok {
			return fmt.Errorf("invalid app type: %T", _app)
		}

		return _app.Set(s, x, y, g.ch.Idx())
	})
	if err != nil {
		return err
	}
	g.state = g.ch.State().Clone()
	return nil
}

func (g *Game) Conclude() error {
	g.Lock()
	defer g.Unlock()

	ctx, cancel := g.c.defaultContextWithTimeout()
	defer cancel()

	err := g.ch.Register(ctx)
	if err != nil {
		return errors.WithMessage(err, "registering")
	}

	err = g.ch.Settle(ctx, false)
	return errors.WithMessage(err, "settling channel")
}
