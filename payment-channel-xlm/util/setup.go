// Copyright 2025 PolyCrypt GmbH
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

package util

import (
	"fmt"
	mathrand "math/rand"
	"path/filepath"
	"runtime"

	pwallet "perun.network/go-perun/wallet"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
	pchannel "perun.network/go-perun/channel"
	"perun.network/perun-stellar-backend/channel"
	"perun.network/perun-stellar-backend/channel/test"
	"perun.network/perun-stellar-backend/channel/types"
	"perun.network/perun-stellar-backend/client"
	"perun.network/perun-stellar-backend/wallet"
)

const (
	PerunContractPath        = "payment-channel-cc/testdata/perun_soroban_contract.wasm"
	StellarAssetContractPath = "payment-channel-cc/testdata/perun_soroban_token.wasm"
	initLumensBalance        = "10000000"
	initTokenBalance         = uint64(2000000)
)

type Setup struct {
	accs     []*wallet.Account
	ws       []*wallet.EphemeralWallet
	cbs      []*client.ContractBackend
	Rng      *mathrand.Rand
	funders  []*channel.Funder
	adjs     []*channel.Adjudicator
	assetIDs []pchannel.Asset
}

func (s *Setup) GetStellarClients() []*client.ContractBackend {
	return s.cbs
}

func (s *Setup) GetFunders() []*channel.Funder {
	return s.funders
}

func (s *Setup) GetAdjudicators() []*channel.Adjudicator {
	return s.adjs
}

func (s *Setup) GetTokenAsset() []pchannel.Asset {
	return s.assetIDs
}

func (s *Setup) GetAccounts() []*wallet.Account {
	return s.accs
}

func (s *Setup) GetWallets() []*wallet.EphemeralWallet {
	return s.ws
}

func getProjectRoot() (string, error) {
	_, b, _, _ := runtime.Caller(1)
	basepath := filepath.Dir(b)

	fp, err := filepath.Abs(filepath.Join(basepath, "../..")) //filepath.Abs(filepath.Join(basepath, "../.."))
	return fp, err
}

func getDataFilePath(filename string) (string, error) {
	root, err := getProjectRoot()
	if err != nil {
		return "", err
	}

	fp := filepath.Join(root, "", filename)
	return fp, nil
}

func NewExampleSetup() (*Setup, error) {

	accs, kpsToFund, ws := test.MakeRandPerunAccsWallets(5)

	if err := test.CreateFundStellarAccounts(kpsToFund, initLumensBalance); err != nil {
		return nil, fmt.Errorf("error funding accounts: %w", err)
	}

	depTokenOneKp := kpsToFund[2]
	depTokenTwoKp := kpsToFund[3]

	depTokenKps := []*keypair.Full{depTokenOneKp, depTokenTwoKp}

	depPerunKp := kpsToFund[4]

	relPathPerun, err := getDataFilePath(PerunContractPath)
	if err != nil {
		return nil, fmt.Errorf("error getting Perun contract path: %w", err)
	}
	relPathAsset, err := getDataFilePath(StellarAssetContractPath)
	if err != nil {
		return nil, fmt.Errorf("error getting asset contract path: %w", err)
	}

	perunAddress, _ := Deploy(depPerunKp, relPathPerun)

	tokenAddressOne, _ := Deploy(depTokenOneKp, relPathAsset)
	tokenAddressTwo, _ := Deploy(depTokenTwoKp, relPathAsset)

	tokenAddresses := []xdr.ScAddress{tokenAddressOne, tokenAddressTwo}
	tokenVector, err := test.MakeCrossAssetVector(tokenAddresses)
	if err != nil {
		return nil, err
	}

	if err = test.InitTokenContract(depTokenOneKp, tokenAddressOne, client.HorizonURL); err != nil {
		return nil, err
	}

	if err = test.InitTokenContract(depTokenTwoKp, tokenAddressTwo, client.HorizonURL); err != nil {
		return nil, err
	}

	SetupAccountsAndContracts(depTokenKps, kpsToFund[:2], tokenAddresses, initTokenBalance)

	var assetContractIDs []pchannel.Asset
	for _, tokenAddress := range tokenAddresses {
		assetContractID, err := types.NewStellarAssetFromScAddress(tokenAddress)
		if err != nil {
			return nil, err
		}
		assetContractIDs = append(assetContractIDs, assetContractID)
	}

	cbs := test.NewContractBackendsFromKeys(kpsToFund[:2], []pwallet.Account{accs[0], accs[1]}, client.HorizonURL)

	aliceCB := cbs[0]
	aliceWallet := ws[0]

	bobCB := cbs[1]
	bobWallet := ws[1]

	channelAccs := []*wallet.Account{accs[0], accs[1]}
	channelCBs := []*client.ContractBackend{aliceCB, bobCB}
	channelWallets := []*wallet.EphemeralWallet{aliceWallet, bobWallet}

	funders, adjs := test.CreateFundersAndAdjudicators(channelAccs, cbs, perunAddress, tokenVector, false)

	setup := Setup{
		accs:     channelAccs,
		ws:       channelWallets,
		cbs:      channelCBs,
		funders:  funders,
		adjs:     adjs,
		assetIDs: assetContractIDs,
	}

	return &setup, nil
}

func SetupAccountsAndContracts(deployerKps []*keypair.Full, kps []*keypair.Full, tokenAddresses []xdr.ScAddress, tokenBalance uint64) {

	for i := range deployerKps {
		for _, kp := range kps {
			addr, err := types.MakeAccountAddress(kp)
			if err != nil {
				panic(err)
			}
			test.MintToken(deployerKps[i], tokenAddresses[i], tokenBalance, addr, client.HorizonURL)
		}
	}
}
