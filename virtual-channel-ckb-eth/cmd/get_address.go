// Copyright 2025 PolyCrypt GmbH
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

package main

import (
	"encoding/hex"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
	"log"
	"os"
	"path/filepath"
	"perun.network/perun-ckb-backend/channel/test"
	"perun.network/perun-ckb-backend/wallet/address"
)

const (
	accountDir      = "accounts"
	ethAuthCodeHash = "0x9c6933d977360f115a3e9cd5a2e0e475853681b80d775d93ad0f8969da343e56" //"0xf329effd1c475a2978453c8600e1eaf0bc2087ee093c3ee64cc96ec6847752cb"
	ckbNetwork      = types.NetworkTest
)

func main() {
	names := []string{"alice", "bob"}

	hash := types.HexToHash(ethAuthCodeHash)

	for _, name := range names {
		pkFile := filepath.Join(accountDir, name+".pk")
		outFile := filepath.Join(accountDir, name+".txt")
		defaultFile := filepath.Join(accountDir, name+"_default.txt")

		if err := processAccount(name, pkFile, outFile, defaultFile, hash); err != nil {
			log.Fatalf("Failed processing %s: %v", name, err)
		}
	}
}

func processAccount(name, pkPath, outPath string, outPathDef string, ethAuthCodeHash types.Hash) error {
	privKey, err := test.GetKey(pkPath)
	if err != nil {
		return err
	}

	participant, ethAddress, err := address.NewEthereumParticipantFromPublicKey(privKey.PubKey(), ethAuthCodeHash)
	if err != nil {
		return fmt.Errorf("create participant: %w", err)
	}
	ckbAddr, _ := participant.ToCKBAddress(ckbNetwork).EncodeFullBech32m()

	// Write result
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "eth_address: 0x%x\n", ethAddress)
	fmt.Fprintf(f, "lock_arg: 0x%x\n", participant.PaymentScript.Args)
	fmt.Fprintf(f, "ckb_address: %s\n", ckbAddr)
	log.Println("Public key:", hex.EncodeToString(privKey.PubKey().SerializeCompressed()))
	defpart, err := address.NewDefaultParticipant(privKey.PubKey())
	if err != nil {
		return fmt.Errorf("create participant: %w", err)
	}
	defCKBAddr, _ := defpart.ToCKBAddress(ckbNetwork).EncodeFullBech32m()
	log.Printf("✅ %s: wrote %s", name, outPath)
	// Write result
	f, err = os.Create(outPathDef)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "eth_address: 0x%x\n", ethAddress)
	fmt.Fprintf(f, "lock_arg: 0x%x\n", defpart.PaymentScript.Args)
	fmt.Fprintf(f, "ckb_address: %s\n", defCKBAddr)
	return nil
}
