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
	defer log.Println("PollBalances: stopped")
	pollingInterval := time.Second
	searchKey := &indexer.SearchKey{
		Script:           address.AsParticipant(p.Account.Address()).PaymentScript,
		ScriptType:       types.ScriptTypeLock,
		ScriptSearchMode: types.ScriptSearchModeExact,
		Filter:           nil,
		WithData:         true,
	}
	log.Println("PollBalances")
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
			//ckbBal := p.balance.Int64()

			// Update sudt balance.
			p.sudtBalance = sudtBalance

			p.balanceMutex.Unlock()
			//p.NotifyAllBalance(ckbBal) // TODO: Update demo tui to allow for big.Int balances
		} else {
			p.balanceMutex.Unlock()
		}
	}
	updateBalance()
	//return "CKbytes: " + p.balance.String() + ", SUDT: " + p.sudtBalance.String()
	/*
		// Poll the balance every 5 seconds.
		for {
			updateBalance()
			time.Sleep(pollingInterval)
		}
	*/
}

func FormatBalance(ckbBal, sudtBal *big.Int) string {
	log.Printf("balances: ckb = %s || sudt = %s", ckbBal.String(), sudtBal.String())
	balCKByte, _ := ShannonToCKByte(ckbBal).Float64()
	return fmt.Sprintf("[green]%s\t[yellow]%s[white]",
		strconv.FormatFloat(balCKByte, 'f', 2, 64)+" CKByte",
		fmt.Sprintf("%v", sudtBal.Int64())+" SUDT")
}
