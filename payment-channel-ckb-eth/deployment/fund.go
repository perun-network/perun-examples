package deployment

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/v2/indexer"
	ckbrpc "github.com/nervosnetwork/ckb-sdk-go/v2/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/v2/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/v2/transaction/signer"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
	"os"
	"path/filepath"
	"perun.network/perun-ckb-backend/wallet/address"
)

// WriteToFile writes the testnet address to a file.
func WriteToFile(testnetAddr string) error {
	dir := filepath.Join("nervos_devnet", "accounts")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(dir, "user.txt"))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("address:\n  testnet: %s\n", testnetAddr))
	return err
}

// FundAccounts funds the given address with the specified amount of CKB.
func FundAccounts(rpcURL string, toAddress address.Participant, amountCKB uint64) error {
	rpcClient, err := ckbrpc.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	indexerClient, err := indexer.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to indexer: %w", err)
	}

	ctx := context.Background()
	keyAlice, err := GetKey("./nervos_devnet/accounts/alice.pk")
	if err != nil {
		return fmt.Errorf("failed to get alice: %w", err)
	}
	aliceParticipant, err := address.NewDefaultParticipant(keyAlice.PubKey())
	if err != nil {
		return fmt.Errorf("failed to make participant from private key: %w", err)
	}
	// 3. Query UTXOs
	capacityNeeded := amountCKB + 100_000 // Add buffer for fee
	liveCells, err := indexerClient.GetCellsCapacity(ctx, &indexer.SearchKey{
		Script:     aliceParticipant.PaymentScript,
		ScriptType: types.ScriptTypeLock,
		Filter:     nil,
	})
	if err != nil {
		return fmt.Errorf("failed to get live cells: %w", err)
	}
	if liveCells.Capacity < capacityNeeded {
		return fmt.Errorf("not enough capacity: %d < %d", liveCells.Capacity, capacityNeeded)
	}

	cells, _, err := collectLiveCells(indexerClient, aliceParticipant.PaymentScript, 100_000)
	if err != nil {
		return fmt.Errorf("failed to collect live cells: %w", err)
	}
	output := &types.CellOutput{
		Capacity: 6100000000,
		Lock:     toAddress.UnlockScript,
	}
	outputs := []*types.CellOutput{output}
	outputsData := [][]byte{[]byte{}}

	// Change output
	totalInputCapacity := uint64(0)
	var inputs []*types.CellInput
	for _, cell := range cells {
		inputs = append(inputs, &types.CellInput{
			Since:          0,
			PreviousOutput: cell.OutPoint,
		})
		totalInputCapacity += cell.Output.Capacity
	}

	changeCapacity := totalInputCapacity - amountCKB - 100_000 // Estimate fee
	if changeCapacity > 0 {
		outputs = append(outputs, &types.CellOutput{
			Capacity: changeCapacity,
			Lock:     aliceParticipant.PaymentScript,
		})
		outputsData = append(outputsData, []byte{})
	}

	tx := &types.Transaction{
		Version: 0,
		CellDeps: []*types.CellDep{
			{
				DepType: types.DepTypeDepGroup,
				OutPoint: &types.OutPoint{
					TxHash: types.HexToHash("0x4299614e5189033916fa258167c537a840f119a6fa63796ec2ccfc0a24619557"),
					Index:  0,
				},
			},
		},
		HeaderDeps:  []types.Hash{},
		Inputs:      inputs,
		Outputs:     outputs,
		OutputsData: outputsData,
		Witnesses:   make([][]byte, len(inputs)),
	}

	// 5. Sign transaction
	txSigner := signer.GetTransactionSignerInstance(types.NetworkTest)
	var inputIndices []uint32
	for i := range tx.Inputs {
		inputIndices = append(inputIndices, uint32(i))
	}
	txWithGroups := &transaction.TransactionWithScriptGroups{
		TxView: tx,
		ScriptGroups: []*transaction.ScriptGroup{
			{
				Script:       aliceParticipant.PaymentScript, // the input lock script (e.g., secp256k1_blake160_sighash_all)
				GroupType:    types.ScriptTypeLock,
				InputIndices: inputIndices,
			},
		},
	}

	if _, err = txSigner.SignTransactionByPrivateKeys(txWithGroups, hex.EncodeToString(keyAlice.Serialize())); err != nil {
		return fmt.Errorf("failed to sign transaction: %d", err)
	}

	// 6. Broadcast
	txHash, err := rpcClient.SendTransaction(ctx, txWithGroups.TxView)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}
	fmt.Println("Transaction hash:", txHash)
	return nil
}

func collectLiveCells(indexerClient indexer.Client, lockScript *types.Script, requiredCapacity uint64) ([]*indexer.LiveCell, uint64, error) {
	searchKey := &indexer.SearchKey{
		Script:     lockScript,
		ScriptType: types.ScriptTypeLock,
	}
	cursor := ""
	liveCells := []*indexer.LiveCell{}
	total := uint64(0)

	for {
		res, err := indexerClient.GetCells(context.Background(), searchKey, indexer.SearchOrderAsc, 100, cursor)
		if err != nil {
			return nil, 0, err
		}

		for _, cell := range res.Objects {
			if cell.Output.Type != nil {
				// Skip type scripts (e.g., SUDT)
				continue
			}
			liveCells = append(liveCells, cell)

			total += cell.Output.Capacity
			if total >= requiredCapacity {
				return liveCells, total, nil
			}
		}

		if res.LastCursor == cursor || len(res.Objects) == 0 {
			break
		}
		cursor = res.LastCursor
	}

	return liveCells, total, nil
}
