package app

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
)

// TicTacToeAppData is the app data struct.
// Grid:
// 0 1 2
// 3 4 5
// 6 7 8
type TicTacToeAppData struct {
	NextActor uint8
	Grid      [9]FieldValue
}

func (d *TicTacToeAppData) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%v|%v|%v\n", d.Grid[0], d.Grid[1], d.Grid[2])
	fmt.Fprintf(&b, "%v|%v|%v\n", d.Grid[3], d.Grid[4], d.Grid[5])
	fmt.Fprintf(&b, "%v|%v|%v\n", d.Grid[6], d.Grid[7], d.Grid[8])
	fmt.Fprintf(&b, "Next actor: %v\n", d.NextActor)
	return b.String()
}

// Encode encodes app data onto an io.Writer.
func (d *TicTacToeAppData) Encode(w io.Writer) error {
	err := writeUInt8(w, d.NextActor)
	if err != nil {
		return errors.WithMessage(err, "writing actor")
	}

	err = writeUInt8Array(w, makeUInt8Array(d.Grid[:]))
	return errors.WithMessage(err, "writing grid")
}

// Clone returns a deep copy of the app data.
func (d *TicTacToeAppData) Clone() channel.Data {
	_d := *d
	return &_d
}

func (d *TicTacToeAppData) Set(x, y int, actorIdx channel.Index) {
	if d.NextActor != uint8safe(uint16(actorIdx)) {
		panic("invalid actor")
	}
	v := makeFieldValueFromPlayerIdx(actorIdx)
	d.Grid[y*3+x] = v
	d.NextActor = calcNextActor(d.NextActor)
}

func calcNextActor(actor uint8) uint8 {
	return (actor + 1) % numParts
}
