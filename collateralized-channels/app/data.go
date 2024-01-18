package app

import (
	"bytes"
	"io"
	"math/big"

	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
)

// CollateralAppData is the app data struct.
type CollateralAppData struct {
	balances [][]*big.Int
}

func (d *CollateralAppData) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	err := d.Encode(&b)
	return b.Bytes(), err
}

func (d *CollateralAppData) UnmarshalBinary(data []byte) error {
	var err error
	buffer := bytes.NewBuffer(data)
	balances, err := readTupleInt256ArrayArray(buffer)
	if err != nil {
		return errors.WithMessage(err, "reading (int256[][])")
	}
	d.balances = balances
	return nil
}

// Encode encodes app data onto an io.Writer.
func (d CollateralAppData) Encode(w io.Writer) (err error) {
	err = writeTupleInt256ArrayArray(w, d.balances)
	return errors.WithMessage(err, "writing (int256[][])")
}

// Clone returns a deep copy of the app data.
func (d CollateralAppData) Clone() channel.Data {
	balances := make([][]*big.Int, len(d.balances))
	for i := range balances {
		balances[i] = make([]*big.Int, len(d.balances[i]))
		for j := range balances[i] {
			balances[i][j] = new(big.Int).Set(d.balances[i][j])
		}
	}
	return &CollateralAppData{balances: balances}
}
