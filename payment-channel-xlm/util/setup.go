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

package util

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	mathrand "math/rand"
	"path/filepath"
	"runtime"
	"time"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	pchannel "perun.network/go-perun/channel"
	"perun.network/perun-stellar-backend/channel"
	"perun.network/perun-stellar-backend/channel/types"
	"perun.network/perun-stellar-backend/client"
	"perun.network/perun-stellar-backend/wallet"
	"perun.network/perun-stellar-backend/wire/scval"
)

const (
	PerunContractPath        = "payment-channel-xlm/testdata/perun_soroban_multi_contract.wasm"
	StellarAssetContractPath = "payment-channel-xlm/testdata/perun_soroban_token.wasm"
	initLumensBalance        = "10000000"
	initTokenBalance         = uint64(2000000)
	DefaultTestTimeout       = 30
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

	accs, kpsToFund, ws := MakeRandPerunAccsWallets(5)

	if err := CreateFundStellarAccounts(kpsToFund, initLumensBalance); err != nil {
		return nil, fmt.Errorf("error funding accounts: %w", err)
	}

	depTokenOneKp := kpsToFund[2]
	depTokenTwoKp := kpsToFund[3]

	depTokenKps := []*keypair.Full{depTokenOneKp, depTokenTwoKp}

	epPerunKp := kpsToFund[4]

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

	InitTokenContract(depTokenOneKp, tokenAddressOne)
	InitTokenContract(depTokenTwoKp, tokenAddressTwo)

	SetupAccountsAndContracts(depTokenKps, kpsToFund[:2], tokenAddresses, initTokenBalance)

	var assetContractIDs []pchannel.Asset

	for _, tokenAddress := range tokenAddresses {
		assetContractID, err := types.NewStellarAssetFromScAddress(tokenAddress)
		if err != nil {
			return nil, err
		}
		assetContractIDs = append(assetContractIDs, assetContractID)
	}

	cbs := NewContractBackendsFromKeys(kpsToFund[:2])

	aliceCB := cbs[0]
	aliceWallet := ws[0]

	bobCB := cbs[1]
	bobWallet := ws[1]

	channelAccs := []*wallet.Account{accs[0], accs[1]}
	channelCBs := []*client.ContractBackend{aliceCB, bobCB}
	channelWallets := []*wallet.EphemeralWallet{aliceWallet, bobWallet}

	funders, adjs := CreateFundersAndAdjudicators(channelAccs, cbs, perunAddress, tokenAddresses)

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
func CreateFundersAndAdjudicators(accs []*wallet.Account, cbs []*client.ContractBackend, perunAddress xdr.ScAddress, tokenScAddresses []xdr.ScAddress) ([]*channel.Funder, []*channel.Adjudicator) {
	funders := make([]*channel.Funder, len(accs))
	adjs := make([]*channel.Adjudicator, len(accs))

	tokenVecAddresses := scval.MakeScVecFromScAddresses(tokenScAddresses)

	for i, acc := range accs {
		funders[i] = channel.NewFunder(acc, cbs[i], perunAddress, tokenVecAddresses)
		adjs[i] = channel.NewAdjudicator(acc, cbs[i], perunAddress, tokenVecAddresses)
	}
	return funders, adjs
}

func NewContractBackendsFromKeys(kps []*keypair.Full) []*client.ContractBackend {
	cbs := make([]*client.ContractBackend, len(kps))
	// generate Configs
	for i, kp := range kps {
		cbs[i] = NewContractBackendFromKey(kp)
	}
	return cbs
}

func NewContractBackendFromKey(kp *keypair.Full) *client.ContractBackend {
	trConfig := client.TransactorConfig{}
	trConfig.SetKeyPair(kp)
	return client.NewContractBackend(&trConfig)
}

func MakeRandPerunAccsWallets(count int) ([]*wallet.Account, []*keypair.Full, []*wallet.EphemeralWallet) {
	accs := make([]*wallet.Account, count)
	kps := make([]*keypair.Full, count)
	ws := make([]*wallet.EphemeralWallet, count)

	for i := 0; i < count; i++ {
		acc, kp, w := MakeRandPerunAccWallet()
		accs[i] = acc
		kps[i] = kp
		ws[i] = w
	}
	return accs, kps, ws
}

func MakeRandPerunAcc() (*wallet.Account, *keypair.Full) {
	w := wallet.NewEphemeralWallet()

	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}

	seed := binary.LittleEndian.Uint64(b[:])

	r := mathrand.New(mathrand.NewSource(int64(seed)))

	acc, kp, err := w.AddNewAccount(r)
	if err != nil {
		panic(err)
	}
	return acc, kp
}

func MakeRandPerunAccWallet() (*wallet.Account, *keypair.Full, *wallet.EphemeralWallet) {
	w := wallet.NewEphemeralWallet()

	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}

	seed := binary.LittleEndian.Uint64(b[:])

	r := mathrand.New(mathrand.NewSource(int64(seed)))

	acc, kp, err := w.AddNewAccount(r)
	if err != nil {
		panic(err)
	}
	return acc, kp, w
}

func CreateFundStellarAccounts(pairs []*keypair.Full, initialBalance string) error {

	numKps := len(pairs)

	masterClient := client.NewHorizonMasterClient()
	masterHzClient := masterClient.GetHorizonClient()
	sourceKey := keypair.Root(client.NETWORK_PASSPHRASE)

	hzClient := client.NewHorizonClient()

	ops := make([]txnbuild.Operation, numKps)

	accReq := horizonclient.AccountRequest{AccountID: masterClient.GetAddress()}
	sourceAccount, err := masterHzClient.AccountDetail(accReq)
	if err != nil {
		panic(err)
	}

	masterAccount := txnbuild.SimpleAccount{
		AccountID: masterClient.GetAddress(),
		Sequence:  sourceAccount.Sequence,
	}

	for i := 0; i < numKps; i++ {
		pair := pairs[i]

		ops[i] = &txnbuild.CreateAccount{
			SourceAccount: masterAccount.AccountID,
			Destination:   pair.Address(),
			Amount:        initialBalance,
		}
	}

	txParams := client.GetBaseTransactionParamsWithFee(&masterAccount, txnbuild.MinBaseFee, ops...)

	txSigned, err := client.CreateSignedTransactionWithParams([]*keypair.Full{sourceKey}, txParams)

	if err != nil {
		panic(err)
	}
	_, err = hzClient.SubmitTransaction(txSigned)
	if err != nil {
		panic(err)
	}

	accounts := make([]txnbuild.Account, numKps)
	for i, kp := range pairs {
		request := horizonclient.AccountRequest{AccountID: kp.Address()}
		account, err := hzClient.AccountDetail(request)
		if err != nil {
			panic(err)
		}

		accounts[i] = &account
	}

	for _, keys := range pairs {
		log.Printf("Funded Stellar L1 account %s (%s) with %s XLM.\n",
			keys.Seed(), keys.Address(), initialBalance)
	}

	return nil
}

func (s *Setup) NewCtx(testTimeout float64) context.Context {
	timeout := time.Duration(testTimeout * float64(time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	if cancel != nil {
		return nil
	}
	return ctx
}
