package ast

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
