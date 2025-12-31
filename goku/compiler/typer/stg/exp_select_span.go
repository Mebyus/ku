package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// SelectSpanLen represents expression of selecting span length.
type SelectSpanLen struct {
	// Span expression being selected.
	Exp Exp

	Pin sm.Pin

	typ *Type
}

func (s *SelectSpanLen) Type() *Type {
	return s.typ
}

func (s *SelectSpanLen) Span() sm.Span {
	return sm.Span{Pin: s.Pin}
}

func (s *SelectSpanLen) String() string {
	return s.Exp.String() + ".len"
}

// SelectSpanPtr represents expression of selecting span pointer.
type SelectSpanPtr struct {
	// Span expression being selected.
	Exp Exp

	Pin sm.Pin

	typ *Type
}

func (s *SelectSpanPtr) Type() *Type {
	return s.typ
}

func (s *SelectSpanPtr) Span() sm.Span {
	return sm.Span{Pin: s.Pin}
}

func (s *SelectSpanPtr) String() string {
	return s.Exp.String() + ".ptr"
}
