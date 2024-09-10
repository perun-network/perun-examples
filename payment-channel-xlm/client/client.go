package client

import (
	"fmt"
	"log"
	"math/big"

	"perun.network/go-perun/wire/net/simple"

	"context"
	"errors"

	pchannel "perun.network/go-perun/channel"
	pclient "perun.network/go-perun/client"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"
	"perun.network/perun-stellar-backend/channel"
	"perun.network/perun-stellar-backend/wallet"
)

type PaymentClient struct {
	perunClient *pclient.Client
	account     *wallet.Account
	currencies  []pchannel.Asset
	channels    chan *PaymentChannel
	Channel     *PaymentChannel
	wAddr       wire.Address
	balance     *big.Int
}

func SetupPaymentClient(
	w *wallet.EphemeralWallet,
	acc *wallet.Account,
	stellarTokenIDs []pchannel.Asset,
	bus *wire.LocalBus,
	funder *channel.Funder,
	adj *channel.Adjudicator,

) (*PaymentClient, error) {
	watcher, err := local.NewWatcher(adj)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}
	// Setup Perun client.
	wireAddr := simple.NewAddress(acc.Address().String())
	perunClient, err := pclient.New(wireAddr, bus, funder, adj, w, watcher)
	if err != nil {
		return nil, errors.New("creating client")
	}

	c := &PaymentClient{
		perunClient: perunClient,
		account:     acc,
		currencies:  stellarTokenIDs,
		channels:    make(chan *PaymentChannel, 1),
		wAddr:       wireAddr,
		balance:     big.NewInt(0),
	}

	go perunClient.Handle(c, c)
	return c, nil
}

// startWatching starts the dispute watcher for the specified channel.
func (c *PaymentClient) startWatching(ch *pclient.Channel) {
	go func() {
		err := ch.Watch(c)
		if err != nil {
			log.Printf("watcher returned with error: %v", err)
		}
	}()
}

func (c *PaymentClient) OpenChannel(peer wire.Address, balances pchannel.Balances) {
	// We define the channel participants. The proposer has always index 0. Here
	// we use the on-chain addresses as off-chain addresses, but we could also
	// use different ones.

	participants := []wire.Address{c.WireAddress(), peer}

	initAlloc := pchannel.NewAllocation(2, c.currencies...)
	initAlloc.Balances = balances

	// Prepare the channel proposal by defining the channel parameters.
	challengeDuration := uint64(10) // On-chain challenge duration in seconds.
	proposal, err := pclient.NewLedgerChannelProposal(
		challengeDuration,
		c.account.Address(),
		initAlloc,
		participants,
	)
	if err != nil {
		panic(err)
	}

	// Send the proposal.
	ch, err := c.perunClient.ProposeChannel(context.TODO(), proposal)
	if err != nil {
		panic(err)
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)
	c.Channel = newPaymentChannel(ch, c.currencies)
}

func (p *PaymentClient) WireAddress() wire.Address {
	return p.wAddr
}

func (c *PaymentClient) AcceptedChannel() *PaymentChannel {
	return <-c.channels
}

// Shutdown gracefully shuts down the client.
func (c *PaymentClient) Shutdown() {
	c.perunClient.Close()
}
