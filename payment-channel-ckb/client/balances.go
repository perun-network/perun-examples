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

package client

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/nervosnetwork/ckb-sdk-go/v2/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
	"perun.network/perun-ckb-backend/channel"
	"perun.network/perun-ckb-backend/wallet/address"
)

type BalanceExtractor func(*indexer.LiveCell) *big.Int

func ckbBalanceExtractor(cell *indexer.LiveCell) *big.Int {
	return new(big.Int).SetUint64(cell.Output.Capacity)
}

func sudtBalanceExtractor(cell *indexer.LiveCell) *big.Int {
	if len(cell.OutputData) != 16 {
		return big.NewInt(0)
	}
	return new(big.Int).SetUint64(binary.LittleEndian.Uint64(cell.OutputData))
}

func (p *PaymentClient) PollBalances() {
	pollingInterval := time.Second
	searchKey := &indexer.SearchKey{
		Script:           address.AsParticipant(p.WalletAddress()[channel.CKBBackendID]).PaymentScript,
		ScriptType:       types.ScriptTypeLock,
		ScriptSearchMode: types.ScriptSearchModeExact,
		Filter:           nil,
		WithData:         true,
	}
	updateBalance := func() {
		ctx, _ := context.WithTimeout(context.Background(), pollingInterval)

		cells, err := p.rpcClient.GetCells(ctx, searchKey, indexer.SearchOrderDesc, math.MaxUint32, "")
		if err != nil {
			log.Println("balance poll error: ", err)
			return
		}
		ckbBalance := big.NewInt(0)
		sudtBalance := big.NewInt(0)
		for _, cell := range cells.Objects {
			ckbBalance = new(big.Int).Add(ckbBalance, ckbBalanceExtractor(cell))
			sudtBalance = new(big.Int).Add(sudtBalance, sudtBalanceExtractor(cell))
		}

		p.balanceMutex.Lock()
		if ckbBalance.Cmp(p.balance) != 0 || sudtBalance.Cmp(p.sudtBalance) != 0 {
			// Update ckb balance.
			p.balance = ckbBalance

			// Update sudt balance.
			p.sudtBalance = sudtBalance

			p.balanceMutex.Unlock()
		} else {
			p.balanceMutex.Unlock()
		}
	}
	updateBalance()

}

func FormatBalance(ckbBal, sudtBal *big.Int) string {
	balCKByte, _ := ShannonToCKByte(ckbBal).Float64()
	return fmt.Sprintf("[green]%s\t[yellow]%s[white]",
		strconv.FormatFloat(balCKByte, 'f', 2, 64)+" CKByte",
		fmt.Sprintf("%v", sudtBal.Int64())+" SUDT")
}
