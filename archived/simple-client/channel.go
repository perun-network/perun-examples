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
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"

	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire/net"
	"perun.network/go-perun/wire/net/simple"
)

func setupFunder(contractBackend ethchannel.ContractBackend, account accounts.Account, assetHolder common.Address) channel.Funder {
	ethDepositor := new(ethchannel.ETHDepositor)
	accounts := map[ethchannel.Asset]accounts.Account{ethwallet.Address(assetHolder): account}
	depositors := map[ethchannel.Asset]ethchannel.Depositor{ethwallet.Address(assetHolder): ethDepositor}
	return ethchannel.NewFunder(contractBackend, accounts, depositors)
}

func setupNetwork(role Role, account wallet.Account) (listener net.Listener, bus *net.Bus, err error) {
	dialer := simple.NewTCPDialer(10 * time.Second)
	dialer.Register(cfg.addrs[1-role], cfg.hosts[1-role])

	listener, err = simple.NewTCPListener(cfg.hosts[role])
	fmt.Printf("Setting up listener for %s\n", cfg.hosts[role])
	if err != nil {
		err = fmt.Errorf("creating listener: %w", err)
		return
	}

	bus = net.NewBus(account, dialer)
	return
}
