package stg

import "github.com/mebyus/ku/internal/ku/sx"

// IndexSpan represents expression using span
// element via index.
type IndexSpan struct {
	operand

	// Span expression being indexed.
	//
	// Always not nil.
	Exp Exp

	// Always not nil.
	Index Exp

	pin sx.Pin

	// type of resulting span element
	typ *Type
}

var _ Operand = &IndexSpan{}

func (x *IndexSpan) Pin() sx.Pin {
	return x.pin
}

func (x *IndexSpan) Type() *Type {
	return x.typ
}

func (*IndexSpan) Target() bool {
	return true
}

// IndexArray represents expression of using array
// element via index.
type IndexArray struct {
	operand

	// Array expression being indexed.
	//
	// Always not nil.
	Exp Exp

	// Always not nil.
	Index Exp

	pin sx.Pin

	// type of resulting array element
	typ *Type
}

var _ Operand = &IndexArray{}

func (x *IndexArray) Pin() sx.Pin {
	return x.pin
}

func (x *IndexArray) Type() *Type {
	return x.typ
}

func (*IndexArray) Target() bool {
	return true
}
