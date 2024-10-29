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
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
	"perun.network/perun-stellar-backend/channel"
	"perun.network/perun-stellar-backend/client"
)

type TokenParams struct {
	decimals uint32
	name     string
	symbol   string
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
	preFlightOp, minFeeInstall, err := client.PreflightHostFunctions(hzClient, &deployerAcc, *installContractOpInstall)
	if err != nil {
		panic(err)
	}

	txParamsInstall := client.GetBaseTransactionParamsWithFee(&deployerAcc, int64(100)+minFeeInstall, &preFlightOp)
	txSignedInstall, err := client.CreateSignedTransactionWithParams([]*keypair.Full{kp}, txParamsInstall)
	if err != nil {
		panic(err)
	}

	_, err = hzClient.SubmitTransaction(txSignedInstall)
	if err != nil {
		panic(err)
	}

	createContractOp := channel.AssembleCreateContractOp(kp.Address(), contractPath, "a1", client.NETWORK_PASSPHRASE)
	preFlightOpCreate, minFeeDeploy, err := client.PreflightHostFunctions(hzClient, &deployerAcc, *createContractOp)
	if err != nil {
		panic(err)
	}
	txParamsCreate := client.GetBaseTransactionParamsWithFee(&deployerAcc, int64(100)+minFeeDeploy, &preFlightOpCreate)
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
