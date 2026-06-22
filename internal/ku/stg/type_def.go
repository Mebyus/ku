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

// Custom defines custom type.
type Custom struct {
	typdef

	// List of methods which are bound to this custom type.
	Methods []*Symbol

	// Custom type is bound to this symbol.
	Symbol *Symbol

	// Type which was used to define this custom type.
	Type *Type

	// Maps method name to its symbol.
	m map[ /* method name */ string]*Symbol
}

type Span struct {
	typdef

	// Span element type.
	Type *Type
}

// Explicit interface implementation check.
var _ TypeDef = &Span{}

type Ref struct {
	typdef

	// Referred type.
	Type *Type
}

// Explicit interface implementation check.
var _ TypeDef = &Ref{}
