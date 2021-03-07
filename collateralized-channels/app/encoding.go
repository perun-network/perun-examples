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
	"io"
	"math/big"

	"github.com/pkg/errors"
)

func readInt256ArrayArray(r io.Reader) (a [][]*big.Int, err error) {
	// Read number of sub-arrays.
	l, err := readUInt256(r)
	if err != nil {
		return nil, errors.WithMessage(err, "reading uint256")
	} else if !l.IsInt64() {
		return nil, errors.Errorf("length out of range: %v", l)
	}
	a = make([][]*big.Int, int(l.Int64()))

	// Read and discard offsets.
	for range a {
		_, err := readUInt256(r)
		if err != nil {
			return nil, errors.WithMessage(err, "reading uint256")
		}
	}

	// Read sub-arrays.
	a = make([][]*big.Int, int(l.Int64()))
	for i := range a {
		v, err := readInt256Array(r)
		if err != nil {
			return nil, errors.WithMessage(err, "reading int256[]")
		}
		a[i] = v
	}
	return a, nil
}

func readInt256Array(r io.Reader) (a []*big.Int, err error) {
	l, err := readUInt256(r)
	if err != nil {
		return nil, errors.WithMessage(err, "reading uint256")
	} else if !l.IsInt64() {
		return nil, errors.Errorf("length out of range: %v", l)
	}
	a = make([]*big.Int, int(l.Int64()))
	for i := range a {
		v, err := readInt256(r)
		if err != nil {
			return nil, errors.WithMessage(err, "reading int256")
		}
		a[i] = v
	}
	return a, nil
}

func readUInt256(r io.Reader) (*big.Int, error) {
	buf := make([]byte, 32)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, errors.WithMessage(err, "reading into byte buffer")
	}

	return new(big.Int).SetBytes(buf), nil
}

func readInt256(r io.Reader) (*big.Int, error) {
	buf := make([]byte, 32)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, errors.WithMessage(err, "reading into byte buffer")
	}

	v := new(big.Int).SetBytes(buf)
	if v.Bit(255) == 1 {
		maxVal := new(big.Int).SetBit(big.NewInt(0), 256, 1)
		v.Sub(v, maxVal)
	}

	return v, nil
}

func writeInt256ArrayArray(w io.Writer, a [][]*big.Int) (err error) {
	l := int64(len(a))

	// Write the number of elements.
	err = writeUInt256(w, big.NewInt(l))
	if err != nil {
		return errors.WithMessage(err, "writing uint256")
	}

	// Encode the subarrays separately because we need to know their length.
	subArrayEncodings := make([][]byte, l)
	for i, v := range a {
		var buf bytes.Buffer
		err = writeInt256Array(&buf, v)
		if err != nil {
			return errors.WithMessage(err, "writing int256[] to buffer")
		}
		subArrayEncodings[i] = buf.Bytes()
	}

	// Write the subarray offsets.
	offset := 32 * l
	for _, v := range subArrayEncodings {
		err = writeUInt256(w, big.NewInt(offset))
		if err != nil {
			return errors.WithMessage(err, "writing uint256")
		}
		offset += int64(len(v))
	}

	// Write subarrays.
	for _, v := range subArrayEncodings {
		_, err = w.Write(v)
		if err != nil {
			return errors.WithMessage(err, "writing byte[]")
		}
	}
	return nil
}

func writeUInt256(w io.Writer, a *big.Int) error {
	buf := a.Bytes()
	_, err := w.Write(make([]byte, 32-len(buf)))
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func writeInt256Array(w io.Writer, a []*big.Int) error {
	l := int64(len(a))
	err := writeUInt256(w, big.NewInt(l))
	if err != nil {
		return errors.WithMessage(err, "writing uint256")
	}

	for _, v := range a {
		err := writeInt256(w, v)
		if err != nil {
			return errors.WithMessage(err, "writing int256")
		}
	}
	return nil
}

func writeInt256(w io.Writer, a *big.Int) error {
	v := new(big.Int).Set(a)
	if a.Sign() < 0 {
		maxVal := new(big.Int).SetBit(big.NewInt(0), 256, 1)
		v.Add(v, maxVal)
	}
	return writeUInt256(w, v)
}

func writeTupleInt256ArrayArray(w io.Writer, a [][]*big.Int) (err error) {
	// Write tuple head: single-element offset.
	offset := 32 // size of uint256
	err = writeUInt256(w, big.NewInt(int64(offset)))
	if err != nil {
		return errors.WithMessage(err, "writing tuple head")
	}

	// Write tuple body.
	err = writeInt256ArrayArray(w, a)
	return errors.WithMessage(err, "writing int256[][]")
}

func readTupleInt256ArrayArray(r io.Reader) (a [][]*big.Int, err error) {
	// Read tuple head.
	_, err = readUInt256(r)
	if err != nil {
		return nil, errors.WithMessage(err, "reading tuple head")
	}

	// Read tuple body.
	a, err = readInt256ArrayArray(r)
	return a, errors.WithMessage(err, "reading int256[][]")
}
