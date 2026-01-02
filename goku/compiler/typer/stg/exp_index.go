package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// SpanIndex represents index expression on expression of span type.
type SpanIndex struct {
	// Expression being indexed. Always has span type.
	Exp Exp

	// Always has integer type.
	Index Exp

	typ *Type
}

func (x *SpanIndex) Type() *Type {
	return x.typ
}

func (x *SpanIndex) Span() sm.Span {
	return x.Index.Span()
}

func (x *SpanIndex) String() string {
	return x.Exp.String() + "[" + x.Index.String() + "]"
}
