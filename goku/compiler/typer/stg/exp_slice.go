package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// SpanSlice represents expression of slice from span expression.
type SpanSlice struct {
	// Always has span type.
	Exp Exp

	// Can be nil if omitted. Must have integer type.
	Start Exp

	// Can be nil if omitted. Must have integer type.
	End Exp
}

func (s *SpanSlice) Type() *Type {
	return s.Exp.Type()
}

func (s *SpanSlice) Span() sm.Span {
	if s.Start != nil {
		return s.Start.Span()
	}
	if s.End != nil {
		return s.End.Span()
	}
	return s.Exp.Span()
}

func (s *SpanSlice) String() string {
	panic("not implemented")
}
