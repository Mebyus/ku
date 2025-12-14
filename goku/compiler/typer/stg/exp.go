package stg

import "github.com/mebyus/ku/goku/compiler/srcmap"

// Exp node that represents an arbitrary expression.
type Exp interface {
	Type() *Type

	Span() srcmap.Span

	// Use only for debugging.
	String() string
}
