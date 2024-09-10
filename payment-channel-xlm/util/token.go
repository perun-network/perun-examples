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
	"errors"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
	"perun.network/perun-stellar-backend/channel"
	"perun.network/perun-stellar-backend/channel/types"
	"perun.network/perun-stellar-backend/client"
	"perun.network/perun-stellar-backend/event"
	"perun.network/perun-stellar-backend/wire/scval"
)

const tokenDecimals = uint32(7)
const tokenName = "PerunToken"
const tokenSymbol = "PRN"

type TokenParams struct {
	decimals uint32
	name     string
	symbol   string
}

func NewTokenParams() *TokenParams {
	return &TokenParams{
		decimals: tokenDecimals,
		name:     tokenName,
		symbol:   tokenSymbol,
	}
}

func (t *TokenParams) GetDecimals() uint32 {
	return t.decimals
}

func (t *TokenParams) GetName() string {
	return t.name
}

func (t *TokenParams) GetSymbol() string {
	return t.symbol
}

func BuildInitTokenArgs(adminAddr xdr.ScAddress, decimals uint32, tokenName string, tokenSymbol string) (xdr.ScVec, error) {

	adminScAddr, err := scval.WrapScAddress(adminAddr)
	if err != nil {
		panic(err)
	}

	decim := xdr.Uint32(decimals)
	scvaltype := xdr.ScValTypeScvU32
	decimSc, err := xdr.NewScVal(scvaltype, decim)
	if err != nil {
		panic(err)
	}

	tokenNameScString := xdr.ScString(tokenName)
	tokenNameXdr := scval.MustWrapScString(tokenNameScString)

	tokenSymbolString := xdr.ScString(tokenSymbol)
	tokenSymbolXdr := scval.MustWrapScString(tokenSymbolString)

	initTokenArgs := xdr.ScVec{
		adminScAddr,
		decimSc,
		tokenNameXdr,
		tokenSymbolXdr,
	}

	return initTokenArgs, nil
}

func InitTokenContract(kp *keypair.Full, contractIDAddress xdr.ScAddress) error {

	cb := NewContractBackendFromKey(kp)

	adminScAddr, err := types.MakeAccountAddress(kp)
	if err != nil {
		panic(err)
	}

	tokenParams := NewTokenParams()
	decimals := tokenParams.GetDecimals()
	name := tokenParams.GetName()
	symbol := tokenParams.GetSymbol()

	initArgs, err := BuildInitTokenArgs(adminScAddr, decimals, name, symbol)
	if err != nil {
		panic(err)
	}

	txMeta, err := cb.InvokeSignedTx("initialize", initArgs, contractIDAddress)
	if err != nil {
		return errors.New("error while invoking and processing host function: initialize" + err.Error())
	}

	_, err = event.DecodeEventsPerun(txMeta)
	if err != nil {
		return err
	}

	return nil
}

func GetTokenName(kp *keypair.Full, contractAddress xdr.ScAddress) error {

	cb := NewContractBackendFromKey(kp)
	TokenNameArgs := xdr.ScVec{}

	_, err := cb.InvokeSignedTx("name", TokenNameArgs, contractAddress)
	if err != nil {
		panic(err)
	}

	return nil
}

func BuildGetTokenBalanceArgs(balanceOf xdr.ScAddress) (xdr.ScVec, error) {

	recScAddr, err := scval.WrapScAddress(balanceOf)
	if err != nil {
		panic(err)
	}

	GetTokenBalanceArgs := xdr.ScVec{
		recScAddr,
	}

	return GetTokenBalanceArgs, nil
}

func BuildTransferTokenArgs(from xdr.ScAddress, to xdr.ScAddress, amount xdr.Int128Parts) (xdr.ScVec, error) {

	fromScAddr, err := scval.WrapScAddress(from)
	if err != nil {
		panic(err)
	}

	toScAddr, err := scval.WrapScAddress(to)
	if err != nil {
		panic(err)
	}

	amountSc, err := scval.WrapInt128Parts(amount)
	if err != nil {
		panic(err)
	}

	GetTokenBalanceArgs := xdr.ScVec{
		fromScAddr,
		toScAddr,
		amountSc,
	}

	return GetTokenBalanceArgs, nil
}

func Deploy(kp *keypair.Full, contractPath string) (xdr.ScAddress, xdr.Hash) {
	deployerCB := NewContractBackendFromKey(kp)
	tr := deployerCB.GetTransactor()
	hzClient := tr.GetHorizonClient()
	deployerAccReq := horizonclient.AccountRequest{AccountID: kp.Address()}
	deployerAcc, err := hzClient.AccountDetail(deployerAccReq)
	if err != nil {
		panic(err)
	}

	installContractOpInstall := channel.AssembleInstallContractCodeOp(kp.Address(), contractPath)
	preFlightOp, _ := client.PreflightHostFunctions(hzClient, &deployerAcc, *installContractOpInstall)

	minFeeInstallCustom := 500000
	txParamsInstall := client.GetBaseTransactionParamsWithFee(&deployerAcc, int64(minFeeInstallCustom), &preFlightOp)
	txSignedInstall, err := client.CreateSignedTransactionWithParams([]*keypair.Full{kp}, txParamsInstall)
	if err != nil {
		panic(err)
	}

	_, err = hzClient.SubmitTransaction(txSignedInstall)
	if err != nil {
		panic(err)
	}

	createContractOp := channel.AssembleCreateContractOp(kp.Address(), contractPath, "a1", client.NETWORK_PASSPHRASE)
	preFlightOpCreate, _ := client.PreflightHostFunctions(hzClient, &deployerAcc, *createContractOp)
	txParamsCreate := client.GetBaseTransactionParamsWithFee(&deployerAcc, int64(minFeeInstallCustom), &preFlightOpCreate)
	txSignedCreate, err := client.CreateSignedTransactionWithParams([]*keypair.Full{kp}, txParamsCreate)
	if err != nil {
		panic(err)
	}

	_, err = hzClient.SubmitTransaction(txSignedCreate)
	if err != nil {
		panic(err)
	}

	contractID := preFlightOpCreate.Ext.SorobanData.Resources.Footprint.ReadWrite[0].MustContractData().Contract.ContractId
	contractHash := preFlightOpCreate.Ext.SorobanData.Resources.Footprint.ReadOnly[0].MustContractCode().Hash
	contractIDAddress := xdr.ScAddress{
		Type:       xdr.ScAddressTypeScAddressTypeContract,
		ContractId: contractID,
	}

	return contractIDAddress, contractHash
}

func MintToken(kp *keypair.Full, contractAddr xdr.ScAddress, amount uint64, recipientAddr xdr.ScAddress) error {
	cb := NewContractBackendFromKey(kp)

	amountTo128Xdr := xdr.Int128Parts{Hi: 0, Lo: xdr.Uint64(amount)}

	amountSc, err := xdr.NewScVal(xdr.ScValTypeScvI128, amountTo128Xdr)
	if err != nil {
		panic(err)
	}
	mintTokenArgs, err := client.BuildMintTokenArgs(recipientAddr, amountSc)
	if err != nil {
		panic(err)
	}
	_, err = cb.InvokeSignedTx("mint", mintTokenArgs, contractAddr)
	if err != nil {
		panic(err)
	}
	return nil
}
