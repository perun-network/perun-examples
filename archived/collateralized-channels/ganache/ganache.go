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

package ganache

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os/exec"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type GanacheConfig struct {
	Cmd       string
	Host      string
	Port      uint
	BlockTime time.Duration
	Funding   []struct {
		PrivateKey string
		BalanceEth uint
	}
	StartupTime   time.Duration
	PrintToStdOut bool
}

type Ganache struct {
	Accounts []Account
	Cmd      *exec.Cmd
}

type Account struct {
	PrivateKey *ecdsa.PrivateKey
	Amount     *big.Int
}

func StartGanacheWithPrefundedAccounts(cfg GanacheConfig) (ganache *Ganache, err error) {
	// Create accounts
	accounts := make([]Account, len(cfg.Funding))
	for i, funding := range cfg.Funding {
		accountKey, err := crypto.HexToECDSA(funding.PrivateKey[2:])
		if err != nil {
			return nil, errors.WithMessage(err, "parsing private key")
		}
		accounts[i] = Account{PrivateKey: accountKey, Amount: ethToWei(big.NewFloat(float64(funding.BalanceEth)))}
	}

	// Build ganache command line arguments
	var ganacheArgs []string
	ganacheArgs = append(ganacheArgs, "--port", fmt.Sprint(cfg.Port))
	for _, a := range accounts {
		key := hexutil.Encode(crypto.FromECDSA(a.PrivateKey))
		ganacheArgs = append(ganacheArgs, "--account", fmt.Sprintf("%v,%v", key, a.Amount))
	}
	ganacheArgs = append(ganacheArgs, fmt.Sprintf("--blockTime=%v", int(cfg.BlockTime.Seconds())))

	// Start command
	ganacheCmdTokens := strings.Split(cfg.Cmd, " ")
	cmdName := ganacheCmdTokens[0]
	var cmdArgs []string
	cmdArgs = append(cmdArgs, ganacheCmdTokens[1:]...)
	cmdArgs = append(cmdArgs, ganacheArgs...)
	cmd := exec.Command(cmdName, cmdArgs...)

	// Print ganache command and arguments for debugging.
	fmt.Println("Ganache command and arguments:", cmdName, cmdArgs)

	// Print ganache output while it is running in the background.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.PrintToStdOut {
		go func() {
			rd := bufio.NewReader(stdout)
			for {
				str, err := rd.ReadString('\n')
				if err != nil {
					fmt.Printf("Reading ganache output: %v\n", err)
					return
				}
				fmt.Print("Ganache output:", str)
			}
		}()
	}

	if err := cmd.Start(); err != nil {
		return nil, errors.WithMessage(err, "starting ganache")
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Wait()
	}()
	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(cfg.StartupTime):
	}
	return &Ganache{accounts, cmd}, nil
}

func (g *Ganache) Shutdown() error {
	return g.Cmd.Process.Kill()
}

func ethToWei(eth *big.Float) (wei *big.Int) {
	var weiPerEth = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	var weiPerEthFloat = new(big.Float).SetInt(weiPerEth)
	wei, _ = new(big.Float).Mul(eth, weiPerEthFloat).Int(nil)
	return
}

func (a *Account) Address() common.Address {
	return crypto.PubkeyToAddress(a.PrivateKey.PublicKey)
}
