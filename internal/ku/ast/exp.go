package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
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

	// Format expression as text that should be very close
	// or identical to what it was parsed from.
	//
	// Used mainly for tests and debugging.
	Format(*fbuf.Buffer)
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

// ParenExp represents an expression which consists of another expression and
// parenthesis surrounding it.
//
// Formal definition:
//
//	ParenExp -> "(" Exp ")"
type ParenExp struct {
	operand

	// Expression surrounded by parenthesis.
	Exp Exp

	// Pin of opening (left) parenthesis.
	Pin sx.Pin
}

// Explicit interface implementation check.
var _ Operand = &ParenExp{}

func (p *ParenExp) Format(buf *fbuf.Buffer) {
	buf.Putb('(')
	p.Exp.Format(buf)
	buf.Putb(')')
}
