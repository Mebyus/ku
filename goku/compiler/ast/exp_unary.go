package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/enums/uok"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// Unary represents unary expression.
//
// Formal definition:
//
//	Unary => UnaryOp ( Operand | Unary )
type Unary struct {
	nodeExp

	// Expression to which operator applies.
	//
	// Can be either:
	//	- operand (majority of cases in practice)
	//	- another unary expression
	Exp Exp

	// Operator in unary expression.
	Op UnaryOp
}

// Explicit interface implementation check.
var _ Exp = Unary{}

func (Unary) Kind() exk.Kind {
	return exk.Unary
}

func (u Unary) Span() srcmap.Span {
	return srcmap.Span{Pin: u.Op.Pin}
}

func (u Unary) String() string {
	var g Printer
	g.Unary(u)
	return g.Output()
}

// UnaryOp represents unary operator inside expression.
type UnaryOp struct {
	Pin  srcmap.Pin
	Kind uok.Kind
}
