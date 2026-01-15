// Copyright 2025 PolyCrypt GmbH
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

package deployment

import (
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// GetKey reads a private key from a file.
func GetKey(path string) (*secp256k1.PrivateKey, error) {

	keyFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()

	rawBytes, err := io.ReadAll(keyFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(rawBytes), "\n")
	if len(lines) == 2 {
		x := strings.Trim(lines[0], " \n")
		xBytes, err := hex.DecodeString(x)
		if err != nil {
			return nil, err
		}
		return secp256k1.PrivKeyFromBytes(xBytes), nil
	} else {
		x := lines[0]
		xBytes, err := hex.DecodeString(x)
		if err != nil {
			return nil, err
		}
		return secp256k1.PrivKeyFromBytes(xBytes), nil
	}
}
