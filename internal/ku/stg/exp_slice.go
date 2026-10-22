package stg

import "github.com/mebyus/ku/internal/ku/sx"

// SliceArray represents expression which slices an array.
type SliceArray struct {
	exp

	// Array expression being sliced.
	Exp Exp

	// Start index.
	A Exp

	// End index.
	B Exp

	pin sx.Pin

	// type of resulting span
	typ *Type

	// if true then there is no need for runtime index check
	safe bool
}

var _ Exp = &SliceArray{}

func (s *SliceArray) Pin() sx.Pin {
	return s.pin
}

func (s *SliceArray) Type() *Type {
	return s.typ
}

func (*SliceArray) Target() bool {
	return false
}
