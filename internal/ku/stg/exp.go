package stg

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
