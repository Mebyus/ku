package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/bok"
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/source"
)

type Binary struct {
	Op BinOp

	// Left side of binary expression.
	A Exp

	// Right side of binary expression.
	B Exp
}

var _ Exp = Binary{}

func (Binary) Kind() exk.Kind {
	return exk.Binary
}

func (b Binary) Span() source.Span {
	return b.A.Span()
}

func (b Binary) String() string {
	var g Printer
	g.Binary(b)
	return g.Str()
}

// BinOp represents binary operator inside expression.
type BinOp struct {
	Pin  source.Pin
	Kind bok.Kind
}

func (o BinOp) String() string {
	return o.Kind.String()
}
