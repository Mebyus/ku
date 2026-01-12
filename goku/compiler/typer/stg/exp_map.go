package stg

import "github.com/mebyus/ku/goku/compiler/sm"

type DerefSelectMapNum struct {
	// Always map pointer (or ref) type.
	Exp Exp

	Pin sm.Pin

	// Always uint type.
	typ *Type
}

func (s *DerefSelectMapNum) Type() *Type {
	return s.typ
}

func (s *DerefSelectMapNum) Span() sm.Span {
	return sm.Span{Pin: s.Pin}
}

func (s *DerefSelectMapNum) String() string {
	return s.Exp.String() + ".*.num"
}
