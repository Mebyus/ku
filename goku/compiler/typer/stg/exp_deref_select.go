package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// DerefSelectField represents expression which selects a struct field
// via pointer dereference.
type DerefSelectField struct {
	// Select target. Always contains a pointer value.
	Exp Exp

	Pin sm.Pin

	Field *Field
}

func (s *DerefSelectField) Type() *Type {
	return s.Field.Type
}

func (s *DerefSelectField) Span() sm.Span {
	return sm.Span{Pin: s.Pin}
}

func (s *DerefSelectField) String() string {
	return s.Exp.String() + ".*." + s.Field.Name
}
