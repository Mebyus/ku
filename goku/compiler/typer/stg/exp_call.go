package stg

import (
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Call expression which directly (not via function pointer) calls
// a specific symbol (function or method).
type Call struct {
	// For methods first arg is always receiver value.
	// The only exception is custom void type receiver in which case
	// it is skipped from list of arguments.
	Args []Exp

	Pin sm.Pin

	Symbol *Symbol

	// Call result type.
	typ *Type
}

func (c *Call) Type() *Type {
	return c.typ
}

func (c *Call) Span() sm.Span {
	return sm.Span{Pin: c.Pin}
}

func (c *Call) String() string {
	panic("not implemented")
	// if c.Symbol.Kind == smk.Fun {
	// 	return c.Symbol.Name + ""
	// }
}
