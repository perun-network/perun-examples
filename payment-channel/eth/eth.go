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

package eth

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"perun.network/go-perun/backend/ethereum/bindings/adjudicator"
	"perun.network/go-perun/backend/ethereum/bindings/assetholdereth"
)

type Client struct {
	ethClient      *ethclient.Client
	key            *ecdsa.PrivateKey
	chainID        *big.Int
	contextTimeout time.Duration
	nonce          int64
}

func NewEthClient(nodeURL string, key *ecdsa.PrivateKey, chainID *big.Int, contextTimeout time.Duration) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()
	client, err := ethclient.DialContext(ctx, nodeURL)
	if err != nil {
		return nil, err
	}
	return &Client{client, key, chainID, contextTimeout, 0}, nil
}

func (c *Client) DeployAdjudicator() (addr common.Address, tx *types.Transaction, err error) {
	return c.deployContract(func(to *bind.TransactOpts, c *ethclient.Client) (addr common.Address, tx *types.Transaction, err error) {
		addr, tx, _, err = adjudicator.DeployAdjudicator(to, c)
		return
	}, false)
}

func (c *Client) DeployAssetHolderETH(adjudicatorAddr common.Address) (addr common.Address, tx *types.Transaction, err error) {
	return c.deployContract(func(to *bind.TransactOpts, c *ethclient.Client) (addr common.Address, tx *types.Transaction, err error) {
		addr, tx, _, err = assetholdereth.DeployAssetHolderETH(to, c, adjudicatorAddr)
		return
	}, false)
}

func (c *Client) deployContract(
	deployContract func(*bind.TransactOpts, *ethclient.Client) (common.Address, *types.Transaction, error),
	waitConfirmation bool,
) (common.Address, *types.Transaction, error) {
	ctx, cancel := c.defaultContext()
	defer cancel()
	ethClient := c.ethClient
	tr, err := c.newTransactor(ctx)
	if err != nil {
		return common.Address{}, nil, err
	}
	addr, tx, err := deployContract(tr, ethClient)
	if err != nil {
		return common.Address{}, nil, errors.WithMessage(err, "sending deployment transaction")
	}

	if waitConfirmation {
		addr, err = bind.WaitDeployed(ctx, ethClient, tx)
		if err != nil {
			return common.Address{}, nil, errors.WithMessage(err, "waiting for the deployment transaction to be mined")
		}
	}
	return addr, tx, nil
}

func (c *Client) WaitDeployment(txs ...*types.Transaction) (err error) {
	ctx, cancel := c.defaultContext()
	defer cancel()
	for _, tx := range txs {
		_, err = bind.WaitDeployed(ctx, c.ethClient, tx)
		if err != nil {
			return errors.WithMessagef(err, "waiting for deployment: %v", tx)
		}
	}
	return nil
}

func (c *Client) defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.contextTimeout)
}

func (c *Client) newTransactor(ctx context.Context) (*bind.TransactOpts, error) {
	tr, err := bind.NewKeyedTransactorWithChainID(c.key, c.chainID)
	if err != nil {
		return nil, err
	}
	tr.Context = ctx
	tr.Nonce = big.NewInt(c.nonce)
	c.nonce++
	// tr.GasPrice = big.NewInt(20000000000)
	// tr.GasLimit = 6721975
	return tr, nil
}

func (c *Client) AccountBalance(a common.Address) (b *big.Int, err error) {
	return c.ethClient.BalanceAt(context.Background(), a, nil)
}

func WeiToEth(weiAmount *big.Int) (ethAmount *big.Float) {
	weiPerEth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	weiPerEthFloat := new(big.Float).SetInt(weiPerEth)
	weiAmountFloat := new(big.Float).SetInt(weiAmount)
	return new(big.Float).Quo(weiAmountFloat, weiPerEthFloat)
}

func EthToWei(ethAmount *big.Float) (weiAmount *big.Int) {
	weiPerEth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	weiPerEthFloat := new(big.Float).SetInt(weiPerEth)
	weiAmountFloat := new(big.Float).Mul(ethAmount, weiPerEthFloat)
	weiAmount, _ = weiAmountFloat.Int(nil)
	return weiAmount
}
