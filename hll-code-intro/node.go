// Copyright (c) 2021, PolyCrypt GmbH, Germany. All rights reserved.
// This file is part of perun-tutorial. Use of this source code is
// governed by the Apache 2.0 license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"math/big"

	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/backend/ethereum/wallet/hd"
	"perun.network/go-perun/channel"
	perun "perun.network/go-perun/client"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire"
)

// node holds all value that we need to conveniently interact with the
// go-perun client.
type node struct {
	// account signs on- and off-chain data.
	account wallet.Account
	// transactor is used for signing Ethereum transactions.
	transactor *hd.Transactor
	// chain connects to the Ethereum chain.
	chain ethchannel.ContractBackend
	// client is the central object for interacting with go-perun.
	client *perun.Client
}

// setupNode sets up everything go-perun and returns as `node`.
func setupNode() *node {
	// Set up a wallet and account.
	account, wallet := setupWallet()
	// Set up a transactor (Ethereum specific).
	transactor := setupTransactor(wallet)
	// Connect to the Ethereum node.
	chain := connectToChain(transactor)
	// Validate the smart-contracts.
	adjudicator, assetholder := setupContracts(chain, account.Account)
	// Setup go-perun networking.
	bus := setupNetworking(account)
	// Setup a funder for the initial channel funding.
	funder := setupFunder(chain, account.Account, assetholder)

	// Create the go-perun client.
	client, err := perun.New(cfg.bobAddr, bus, funder, adjudicator, wallet)
	noError(err)

	// Create the node that defines all event handlers for go-perun.
	node := &node{account, transactor, chain, client}
	// Start Proposal- and UpdateHandlers.
	go client.Handle(node, node)

	return node
}

// HandleProposal is called by go-perun for each channel that is being proposed
// to Bob. Bob can use the responder to accept or reject the proposal.
func (n *node) HandleProposal(perun.ChannelProposal, *perun.ProposalResponder) {}

// HandleUpdate is called by go-perun for every update of a channel that
// Alice wants to do. Bob can then accept or reject it with the `responder`.
func (n *node) HandleUpdate(*channel.State, perun.ChannelUpdate, *perun.UpdateResponder) {}

// HandleAdjudicatorEvent is called by go-perun when it detects an event in the
// Adjudicator smart-contract.
func (n *node) HandleAdjudicatorEvent(channel.AdjudicatorEvent) {}

// openChannel opens a payment-channel between Alice and Bob.
// The passed amount is the initial balance of Alice and Bob in the channel.
func (n *node) openChannel(initAmount *big.Int) *perun.Channel {
	fmt.Println("Opening channel")
	// Set the initial balances for each peer and each asset in the channel.
	initBals := &channel.Allocation{
		Assets:   []channel.Asset{ethwallet.AsWalletAddr(cfg.assetAddr)},
		Balances: [][]*big.Int{{ /*Bob*/ initAmount, initAmount /*Alice*/}},
	}
	// Define the participants of the channel.
	// Bob goes first since he proposes the channel.
	parts := []wire.Address{
		cfg.bobAddr,
		cfg.aliceAddr,
	}
	// Create a ledger channel proposal.
	proposal, err := perun.NewLedgerChannelProposal(10, cfg.bobAddr, initBals, parts)
	noError(err)
	// Propose the channel.
	channel, err := n.client.ProposeChannel(context.Background(), proposal)
	noError(err)
	// Start the on-chain watcher.
	go channel.Watch(n)

	fmt.Printf("Opened channel with id 0x%x \n", channel.ID())
	return channel
}

// sendBalance sends the specified `amount` to Alice using the payment-channel.
func sendBalance(ch *perun.Channel, amount *big.Int) {
	err := ch.UpdateBy(context.Background(), func(state *channel.State) error {
		// Subtract the amount from Bobs balance.
		bobBal := new(big.Int).Sub(state.Balances[assetIdx][bobIndex], amount)
		state.Balances[assetIdx][bobIndex] = bobBal
		// Add the amount Alice' balance.
		aliceBal := new(big.Int).Add(state.Balances[assetIdx][aliceIndex], amount)
		state.Balances[assetIdx][aliceIndex] = aliceBal

		return nil
	})
	noError(err)
}

// closeChannel closes the channel with Alice and withdraws all funds.
func closeChannel(ch *perun.Channel) {
	fmt.Println("Closing the channel")
	// Send a final update to Alice.
	err := ch.UpdateBy(context.Background(), func(state *channel.State) error {
		state.IsFinal = true
		return nil
	})
	noError(err)
	// Settle the channel and withdraw on-chain funds.
	err = ch.Settle(context.Background(), false)
	noError(err)
	// Free remaining go and go-perun resources.
	err = ch.Close()
	noError(err)
}
