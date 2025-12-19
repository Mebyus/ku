package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// Exp node that represents an arbitrary expression.
type Exp interface {
	Type() *Type

	Span() sm.Span

	// Use only for debugging.
	String() string
}
