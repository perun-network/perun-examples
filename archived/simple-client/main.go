// Copyright 2021 PolyCrypt GmbH, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"perun.network/go-perun/client"
)

func main() {
	// Setup Alice and Bob.
	alice, bob := setup(RoleAlice), setup(RoleBob)
	// Run our example protocol: Bob Opens, Updates and Closes.
	if err := bob.openChannel(); err != nil {
		panic(fmt.Errorf("opening channel: %w", err))
	}
	if err := bob.updateChannel(); err != nil {
		panic(fmt.Errorf("updating channel: %w", err))
	}
	if err := bob.closeChannel(); err != nil {
		panic(fmt.Errorf("closing channel: %w", err))
	}
	// Wait for both nodes to stop.
	fmt.Println("Waiting for Alice")
	<-alice.done
	fmt.Println("Waiting for Bob")
	<-bob.done
}

func setup(role Role) *node {
	fmt.Println("Starting ", role)
	account, wallet, err := setupWallet(role)
	if err != nil {
		panic(fmt.Sprintf("setting up wallet: %v", err))
	}
	transactor := createTransactor(wallet)

	_, contractBackend, err := connectToChain(transactor)
	if err != nil {
		panic(fmt.Sprintf("connecting to chain: %v", err))
	}

	adjudicator, assetholder, err := setupContracts(role, contractBackend, account.Account)
	if err != nil {
		panic(fmt.Errorf("setting up contracts: %w", err))
	}

	listener, bus, err := setupNetwork(role, account)
	if err != nil {
		panic(fmt.Errorf("setting up network: %w", err))
	}

	funder := setupFunder(contractBackend, account.Account, assetholder)
	cl, err := client.New(cfg.addrs[role], bus, funder, adjudicator, wallet)
	if err != nil {
		panic(fmt.Errorf("creating client: %w", err))
	}
	// Create the node that defines all event handlers for go-perun.
	node := &node{role: role, account: account, transactor: transactor,
		contractBackend: contractBackend, assetholder: assetholder, listener: listener,
		bus: bus, client: cl, ch: nil, done: make(chan struct{})}
	// Set the NewChannel handler.
	cl.OnNewChannel(node.HandleNewChannel)
	// Start Proposal- and UpdateHandlers.
	go cl.Handle(node, node)
	// Listen on incoming connections.
	go bus.Listen(listener)
	return node
}
