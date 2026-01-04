package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// SelectField represents expression of selecting (via offset, without indirection)
// a field from struct or union.
type SelectField struct {
	// Expression being selected. Always struct or union.
	Exp Exp

	Pin sm.Pin

	Field *Field
}

func (s *SelectField) Type() *Type {
	return s.Field.Type
}

func (s *SelectField) Span() sm.Span {
	return sm.Span{Pin: s.Pin}
}

func (s *SelectField) String() string {
	panic("not implemented")
}
