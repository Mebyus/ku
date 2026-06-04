package stg

import (
	"github.com/mebyus/ku/internal/ku/enums/bop"
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

// Integer represents an integer constant (directly from source or evaluated)
// which value is known at compile time.
type Integer struct {
	operand

	Pin sx.Pin

	Val uint64

	typ *Type

	Neg bool
}

// create static integer value.
func (t *Typer) makeInteger(pin sx.Pin, v uint64) *Integer {
	return &Integer{
		Pin: pin,
		Val: v,
		typ: t.common.Types.Static.Integer,
	}
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

// Boolean represents a boolean constant true/false (directly from source or evaluated)
// which value is known at compile time.
type Boolean struct {
	operand

	Pin sx.Pin

	typ *Type

	Val bool
}

// Explicit interface implementation check.
var _ Operand = &Boolean{}

// create static boolean value.
func (t *Typer) makeBoolean(pin sx.Pin, v bool) *Boolean {
	return &Boolean{
		Pin: pin,
		Val: v,
		typ: t.common.Types.Static.Boolean,
	}
}
