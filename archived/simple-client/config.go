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

	"github.com/ethereum/go-ethereum/common"

	"perun.network/go-perun/backend/ethereum/wallet"
)

type config struct {
	mnemonic string                   // Ethereum wallet mnemonic.
	chainURL string                   // Url of the Ethereum node.
	hosts    map[Role]string          // Hosts for incoming connections.
	addrs    map[Role]*wallet.Address // Wallet addresses of both roles.
}

// cfg saves the global configuration.
var cfg config

type Role int

const (
	RoleAlice Role = iota
	RoleBob
)

func (r Role) String() string {
	switch r {
	case RoleAlice:
		return "Alice"
	case RoleBob:
		return "Bob"
	}
	return fmt.Sprintf("%d", r)
}

func init() {
	cfg.mnemonic = "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic"
	cfg.chainURL = "ws://127.0.0.1:8545"
	// Set the Host that the go-perun Client will bind to.
	cfg.hosts = map[Role]string{
		RoleAlice: "0.0.0.0:8401",
		RoleBob:   "0.0.0.0:8402",
	}
	// Fix the on-chain addresses of Alice and Bob.
	cfg.addrs = map[Role]*wallet.Address{
		RoleAlice: wallet.AsWalletAddr(common.HexToAddress("0x2EE1ac154435f542ECEc55C5b0367650d8A5343B")),
		RoleBob:   wallet.AsWalletAddr(common.HexToAddress("0x70765701b79a4e973dAbb4b30A72f5a845f22F9E")),
	}
}
