package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// MakeSpan represents expression of making a span from array pointer
// or array ref.
type MakeSpan struct {
	// Always not nil. Must have array ref or array pointer type.
	Exp Exp

	// Could be nil, if start is omitted.
	Start Exp

	// Always not nil.
	End Exp

	typ *Type
}

func (s *MakeSpan) Type() *Type {
	return s.typ
}

func (s *MakeSpan) Span() sm.Span {
	return s.End.Span()
}

func (s *MakeSpan) String() string {
	panic("not implemented")
}
