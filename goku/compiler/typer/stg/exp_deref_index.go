package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// DerefIndex represents deref index expression on array pointer or ref pointer value.
type DerefIndex struct {
	// Expression being indexed. Always not nil.
	Exp Exp

	// Always not nil.
	Index Exp

	typ *Type
}

func (n *DerefIndex) Type() *Type {
	return n.typ
}

func (n *DerefIndex) Span() sm.Span {
	return n.Index.Span()
}

func (n *DerefIndex) String() string {
	return n.Exp.String() + ".[" + n.Index.String() + "]"
}
