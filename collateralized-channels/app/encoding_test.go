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

package app

import (
	"bytes"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

const rndSeed = 0
const rndTestCount = 16

var rnd = rand.New(rand.NewSource(rndSeed))

func TestEncodeUInt256(t *testing.T) {
	var (
		assert = assert.New(t)
		err    error
	)

	for i := 0; i < rndTestCount; i++ {
		var buf bytes.Buffer
		a := rndUInt256()
		err = writeUInt256(&buf, a)
		assert.NoError(err, "writing")

		b, err := readUInt256(&buf)
		assert.NoError(err, "reading")
		assert.Equal(a, b, "comparing encoded with decoded value")
	}
}

func rndUInt256() *big.Int {
	return new(big.Int).Rand(rnd, new(big.Int).SetBit(big.NewInt(0), 256, 1))
}

func TestEncodeInt256(t *testing.T) {
	var (
		assert = assert.New(t)
		err    error
	)

	for i := 0; i < rndTestCount; i++ {
		var buf bytes.Buffer
		a := rndInt256()
		err = writeInt256(&buf, a)
		assert.NoError(err, "writing")

		b, err := readInt256(&buf)
		assert.NoError(err, "reading")
		assert.Equal(a, b, "comparing encoded with decoded value")
	}
}

func rndInt256() *big.Int {
	a := rndUInt256()
	return new(big.Int).Sub(a, new(big.Int).SetBit(big.NewInt(0), 255, 1))
}

func TestEncodeInt256ArrayArray(t *testing.T) {
	var (
		assert      = assert.New(t)
		err         error
		lengthRange = 10
	)

	for i := 0; i < rndTestCount; i++ {
		var buf bytes.Buffer
		a := rndInt256ArrayArray(lengthRange)
		err = writeInt256ArrayArray(&buf, a)
		assert.NoError(err, "writing")

		b, err := readInt256ArrayArray(&buf)
		assert.NoError(err, "reading")
		assert.Equal(a, b, "comparing encoded with decoded value")
	}
}

func rndInt256ArrayArray(lengthRange int) [][]*big.Int {
	l1 := 2 + rand.Intn(lengthRange)
	a := make([][]*big.Int, l1)
	for i := range a {
		l2 := 2 + rand.Intn(lengthRange)
		a[i] = make([]*big.Int, l2)
		for j := range a[i] {
			a[i][j] = rndInt256()
		}
	}
	return a
}

func TestEncodeTupleInt256ArrayArray(t *testing.T) {
	var (
		assert      = assert.New(t)
		err         error
		lengthRange = 10
	)

	for i := 0; i < rndTestCount; i++ {
		var buf bytes.Buffer
		a := rndInt256ArrayArray(lengthRange)
		err = writeTupleInt256ArrayArray(&buf, a)
		assert.NoError(err, "writing")

		b, err := readTupleInt256ArrayArray(&buf)
		assert.NoError(err, "reading")
		assert.Equal(a, b, "comparing encoded with decoded value")
	}
}
