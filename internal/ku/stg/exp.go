package stg

import (
	"github.com/mebyus/ku/internal/ku/enums/bop"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Exp node that represents an arbitrary (but not empty) expression.
type Exp interface {
	Pin() sx.Pin

	// Must return value's type of evaluated expression.
	//
	// Can return nil only for function calls that return nothing.
	Type() *Type

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

func (operand) _exp()     {}
func (operand) _operand() {}

// Integer represents an integer constant (directly from source or evaluated)
// which value is known at compile time.
type Integer struct {
	operand

	pin sx.Pin

	Val uint64

	typ *Type

	Neg bool
}

// Explicit interface implementation check.
var _ Exp = &Integer{}

func (n *Integer) Type() *Type {
	return n.typ
}

func (n *Integer) Pin() sx.Pin {
	return n.pin
}

// create static integer value.
func (t *Typer) makeInteger(pin sx.Pin, v uint64) *Integer {
	return &Integer{
		pin: pin,
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

	pin sx.Pin

	typ *Type
}

// Explicit interface implementation check.
var _ Exp = &BinExp{}

func (b *BinExp) Type() *Type {
	return b.typ
}

func (b *BinExp) Pin() sx.Pin {
	return b.pin
}

// Boolean represents a boolean constant true/false (directly from source or evaluated)
// which value is known at compile time.
type Boolean struct {
	operand

	pin sx.Pin

	typ *Type

	Val bool
}

// Explicit interface implementation check.
var _ Operand = &Boolean{}

func (b *Boolean) Type() *Type {
	return b.typ
}

func (b *Boolean) Pin() sx.Pin {
	return b.pin
}

// create static boolean value.
func (t *Typer) makeBoolean(pin sx.Pin, v bool) *Boolean {
	return &Boolean{
		pin: pin,
		Val: v,
		typ: t.common.Types.Static.Boolean,
	}
}

// SymExp represents expression formed from symbol usage (as operand or part of chain)
// inside an expression.
type SymExp struct {
	operand

	pin sx.Pin

	// Type of expression. It may differ from type declared for symbol.
	typ *Type

	Symbol *Symbol
}

// Explicit interface implementation check.
var _ Operand = &SymExp{}

func (s *SymExp) Type() *Type {
	return s.typ
}

func (s *SymExp) Pin() sx.Pin {
	return s.pin
}
