package ast

import (
	"github.com/mebyus/ku/internal/ku/enums/bop"
	"github.com/mebyus/ku/internal/ku/enums/uop"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Exp node that represents an arbitrary (but not empty) expression.
type Exp interface {
	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_exp()
}

// Embed this to quickly implement _exp() discriminator from Exp interface.
// Do not use it for anything else.
type exp struct{}

// Explicit interface implementation check.
var _ Exp = exp{}

func (exp) _exp() {}

// Operand node that represents an expression which can be used as operand
// inside another expressions.
//
// Each operand can be used as standalone expression, but not every expression
// is an operand.
type Operand interface {
	Exp

	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_operand()
}

// Embed this to quickly implement _operand() discriminator from Operand interface.
// Do not use it for anything else.
type operand struct{}

// Explicit interface implementation check.
var _ Operand = operand{}

func (operand) _exp()     {}
func (operand) _operand() {}

// Unary represents unary expression.
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
