package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

type Rune struct {
	operand

	// Rune literal value represented by token.
	Val uint64

	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Operand = &Rune{}

func (r *Rune) Format(buf *fbuf.Buffer) {
	if r.Val > 127 {
		panic("stub")
	}

	buf.Putb('\'')
	buf.Putb(byte(r.Val))
	buf.Putb('\'')
}
