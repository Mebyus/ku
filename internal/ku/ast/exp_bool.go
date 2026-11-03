package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

type True struct {
	operand

	Pin sx.Pin
}

type False struct {
	operand

	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Operand = &True{}
var _ Operand = &False{}

func (*True) Format(buf *fbuf.Buffer) {
	buf.Puts("true")
}

func (*False) Format(buf *fbuf.Buffer) {
	buf.Puts("false")
}
