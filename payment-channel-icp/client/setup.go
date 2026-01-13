// Copyright 2025 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"perun.network/go-perun/client"
	pwallet "perun.network/go-perun/wallet"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"
	"perun.network/perun-icp-backend/channel"
	chanconn "perun.network/perun-icp-backend/channel/connector"
	"perun.network/perun-icp-backend/wallet"
	icwallet "perun.network/perun-icp-backend/wallet"
)

func SetupPaymentClient(
	name string,
	w *wallet.FsWallet, // w is the wallet used to resolve addresses to accounts for channels.
	wireAcc wire.Account,
	bus wire.Bus,
	perunID string,
	ledgerID string,
	host string,
	port int,
	accountPath string,
) (*PaymentClient, error) {

	acc := w.NewAccount()

	// Connect to Perun pallet and get funder + adjudicator from it.

	perunConn := chanconn.NewICConnector(perunID, ledgerID, accountPath, host, port)

	funder := channel.NewFunder(acc, perunConn)
	adj := channel.NewAdjudicator(acc, perunConn)

	// Setup dispute watcher.
	watcher, err := local.NewWatcher(adj)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}

	// Setup Perun client.
	wallet := map[pwallet.BackendID]pwallet.Wallet{
		icwallet.ICPBackendID: w,
	}
	wireAddr := map[pwallet.BackendID]wire.Address{
		icwallet.ICPBackendID: wireAcc.Address(),
	}
	perunClient, err := client.New(wireAddr, bus, funder, adj, wallet, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	// Create client and start request handler.
	c := &PaymentClient{
		Name:        name,
		perunClient: perunClient,
		account:     &acc,
		currency:    channel.Asset,
		channels:    make(chan *PaymentChannel, 1),
		ICConn:      perunConn,
		wAddr:       wireAddr,
		balance:     big.NewInt(0),
	}

	go c.PollBalances()
	go perunClient.Handle(c, c)
	return c, nil
}
