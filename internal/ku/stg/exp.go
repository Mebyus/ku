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

	Val uint64

	pin sx.Pin

	typ *Type

	Neg bool
}

// Explicit interface implementation check.
var _ Operand = &Integer{}

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

// Integer represents a string constant (directly from source or evaluated)
// which value is known at compile time.
type String struct {
	operand

	Val string

	pin sx.Pin

	typ *Type
}

var _ Operand = &String{}

func (s *String) Type() *Type {
	return s.typ
}

func (s *String) Pin() sx.Pin {
	return s.pin
}

// create static integer value.
func (t *Typer) makeString(pin sx.Pin, s string) *String {
	return &String{
		pin: pin,
		Val: s,
		typ: t.common.Types.Static.String,
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

// SpanNum represents expression of selecting span elements number via ".num".
type SpanNum struct {
	operand

	// Span expression being selected.
	Exp Exp

	// pin of select place (not its target)
	pin sx.Pin

	// always uint type
	typ *Type
}

var _ Operand = &SpanNum{}

func (s *SpanNum) Type() *Type {
	return s.typ
}

func (s *SpanNum) Pin() sx.Pin {
	return s.pin
}
