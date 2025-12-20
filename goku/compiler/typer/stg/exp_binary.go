package stg

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/sm"
)

type BinOp = ast.BinOp

// Binary represents a binary expression.
type Binary struct {
	Op BinOp

	A Exp

	B Exp

	typ *Type
}

func (b *Binary) Type() *Type {
	return b.typ
}

func (b *Binary) Span() sm.Span {
	return sm.Span{Pin: b.A.Span().Pin}
}

func (b *Binary) String() string {
	return b.A.String() + " " + b.Op.String() + " " + b.B.String()
}
