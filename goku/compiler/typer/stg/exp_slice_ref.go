package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// SliceArrayRef repsents expression which slices expression of array ref type
// to obtain new array ref shifted by a specific value from original.
// This operation is similar to pointer arithmetic in C.
type SliceArrayRef struct {
	// Expression being sliced. Always has array ref type.
	Exp Exp

	// Always has integer type.
	Index Exp
}

func (r *SliceArrayRef) Type() *Type {
	return r.Exp.Type()
}

func (r *SliceArrayRef) Span() sm.Span {
	return r.Index.Span()
}

func (r *SliceArrayRef) String() string {
	panic("not implemented")
}
