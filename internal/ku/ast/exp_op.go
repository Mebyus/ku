package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/enums/bop"
	"github.com/mebyus/ku/internal/ku/enums/uop"
)

// UnExp represents unary expression.
//
// Formal definition:
//
//	Unary -> UnaryOp ( Operand | Unary )
type UnExp struct {
	operand

	// Expression to which operator applies.
	//
	// Can be either:
	//	- operand (majority of cases in practice)
	//	- another unary expression
	A Exp

	// Operator in unary expression.
	Op uop.Op
}

// Explicit interface implementation check.
var _ Operand = &UnExp{}

func (e *UnExp) Format(buf *fbuf.Buffer) {
	buf.Puts(e.Op.Kind.String())
	e.A.Format(buf)
}

// BinExp represents binary expression.
type BinExp struct {
	exp

	// Left side of binary expression.
	A Exp

	// Right side of binary expression.
	B Exp

	Op bop.Op
}

// Explicit interface implementation check.
var _ Exp = &BinExp{}

func (e *BinExp) Format(buf *fbuf.Buffer) {
	e.A.Format(buf)
	buf.Space()
	buf.Puts(e.Op.Kind.String())
	buf.Space()
	e.B.Format(buf)
}
