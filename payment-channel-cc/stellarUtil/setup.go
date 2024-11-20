// Copyright 2024 PolyCrypt GmbH
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

package stellarUtil

import (
	"context"
	"crypto/elliptic"
	"fmt"
	"log"
	mathrand "math/rand"
	"path/filepath"
	pwallet "perun.network/go-perun/wallet"
	"runtime"
	"time"

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

func NewExampleSetup(sk []string, addr [][20]byte) (*Setup, error) {

	_, kpsToFundR, _ := test.MakeRandPerunAccsWallets(3)
	accs, kpsToFund, ws := MakePerunAccsWallets(2, sk, addr)
	log.Println("secret keys: ", sk)
	pubKey := make([][]byte, 2)
	for i := 0; i < 2; i++ {
		pk := accs[i].Participant().StellarPubKey
		pubKey[i] = elliptic.Marshal(pk.Curve, pk.X, pk.Y)
	}
	log.Println("publicKeys: ", pubKey)
	log.Println("addresses:", addr)

	if err := test.CreateFundStellarAccounts(kpsToFund, initLumensBalance); err != nil {
		return nil, fmt.Errorf("error funding accounts: %w", err)
	}
	if err := test.CreateFundStellarAccounts(kpsToFundR, initLumensBalance); err != nil {
		return nil, fmt.Errorf("error funding accounts: %w", err)
	}

	depTokenOneKp := kpsToFundR[0]
	depTokenTwoKp := kpsToFundR[1]

	depTokenKps := []*keypair.Full{depTokenOneKp, depTokenTwoKp}

	epPerunKp := kpsToFundR[2]

	relPathPerun, err := getDataFilePath(PerunContractPath)
	if err != nil {
		return nil, fmt.Errorf("error getting Perun contract path: %w", err)
	}
	relPathAsset, err := getDataFilePath(StellarAssetContractPath)
	if err != nil {
		return nil, fmt.Errorf("error getting asset contract path: %w", err)
	}

	perunAddress, _ := Deploy(epPerunKp, relPathPerun)

	tokenAddressOne, _ := Deploy(depTokenOneKp, relPathAsset)
	tokenAddressTwo, _ := Deploy(depTokenTwoKp, relPathAsset)

	tokenAddresses := []xdr.ScAddress{tokenAddressOne, tokenAddressTwo}
	tokenVector, err := test.MakeCrossAssetVector(tokenAddresses)
	if err != nil {
		return nil, err
	}

	test.InitTokenContract(depTokenOneKp, tokenAddressOne)
	test.InitTokenContract(depTokenTwoKp, tokenAddressTwo)

	SetupAccountsAndContracts(depTokenKps, kpsToFund, tokenAddresses, initTokenBalance)

	var assetContractIDs []pchannel.Asset
	for _, tokenAddress := range tokenAddresses {
		assetContractID, err := types.NewStellarAssetFromScAddress(tokenAddress)
		if err != nil {
			return nil, err
		}
		assetContractIDs = append(assetContractIDs, assetContractID)
	}

	cbs := test.NewContractBackendsFromKeys(kpsToFund, []pwallet.Account{accs[0], accs[1]})

	aliceCB := cbs[0]
	aliceWallet := ws[0]

	bobCB := cbs[1]
	bobWallet := ws[1]

	channelAccs := []*wallet.Account{accs[0], accs[1]}
	channelCBs := []*client.ContractBackend{aliceCB, bobCB}
	channelWallets := []*wallet.EphemeralWallet{aliceWallet, bobWallet}

	funders, adjs := test.CreateFundersAndAdjudicators(channelAccs, cbs, perunAddress, tokenVector)

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
			MintToken(deployerKps[i], tokenAddresses[i], tokenBalance, addr)
		}
	}
}

func NewContractBackendFromKey(kp *keypair.Full) *client.ContractBackend {
	trConfig := client.TransactorConfig{}
	trConfig.SetKeyPair(kp)
	return client.NewContractBackend(&trConfig)
}

func MakePerunAccsWallets(count int, sk []string, addr [][20]byte) ([]*wallet.Account, []*keypair.Full, []*wallet.EphemeralWallet) {
	accs := make([]*wallet.Account, count)
	kps := make([]*keypair.Full, count)
	ws := make([]*wallet.EphemeralWallet, count)

	for i := 0; i < count; i++ {
		acc, kp, w := MakePerunAccWallet(sk[i], addr[i])
		accs[i] = acc
		kps[i] = kp
		ws[i] = w
	}
	return accs, kps, ws
}

func MakePerunAccWallet(sk string, addr [20]byte) (*wallet.Account, *keypair.Full, *wallet.EphemeralWallet) {
	w := wallet.NewEphemeralWallet()

	kp, err := keypair.Random()
	if err != nil {
		panic(err)
	}
	acc := wallet.NewAccount(sk, *kp.FromAddress(), addr)
	err = w.AddAccount(acc)
	if err != nil {
		panic(err)
	}
	return acc, kp, w
}

func (s *Setup) NewCtx(testTimeout float64) context.Context {
	timeout := time.Duration(testTimeout * float64(time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	if cancel != nil {
		return nil
	}
	return ctx
}
