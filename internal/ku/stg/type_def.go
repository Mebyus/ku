package stg

type TypeDef interface {
	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_typdef()
}

type typdef struct{}

func (typdef) _typdef() {}

type Span struct {
	typdef

	// Span element type.
	Type *Type
}

// Explicit interface implementation check.
var _ TypeDef = &Span{}
