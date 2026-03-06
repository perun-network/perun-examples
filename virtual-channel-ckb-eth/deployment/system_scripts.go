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
	"encoding/json"
	"io"
	"os"
	"path"
	"perun.network/perun-ckb-backend/channel/test"
)

const systemScriptName = "default_scripts.json"

// GetSystemScripts reads the system scripts from a file.
func GetSystemScripts(systemScriptDir string) (test.SystemScripts, error) {
	var ss test.SystemScripts
	err := readJSON(systemScriptDir, &ss)
	if err != nil {
		return test.SystemScripts{}, err
	}
	return ss, nil
}

func readJSON(systemScriptDir string, systemScripts *test.SystemScripts) error {
	systemScriptFile, err := os.Open(path.Join(systemScriptDir, systemScriptName))
	defer func() { _ = systemScriptFile.Close() }()
	if err != nil {
		return err
	}
	systemScriptContent, err := io.ReadAll(systemScriptFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(systemScriptContent, systemScripts)
}
