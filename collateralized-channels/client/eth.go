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

package client

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) transactAndConfirm(transact func(tr *bind.TransactOpts) (*types.Transaction, error)) (err error) {
	ctx := c.defaultContext()

	tr, err := c.perunClient.Wallet.NewTransactor(c.perunClient.Account.Account)
	tr.Context = ctx
	if err != nil {
		return errors.WithMessage(err, "creating transactor")
	}

	tx, err := transact(tr)
	if err != nil {
		return errors.WithMessage(err, "submitting transaction")
	}

	receipt, err := bind.WaitMined(ctx, c.perunClient.ContractBackend, tx)
	if err != nil {
		return errors.WithMessage(err, "awaiting transaction confirmation")
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return errors.New("transaction failed")
	}

	return nil
}
