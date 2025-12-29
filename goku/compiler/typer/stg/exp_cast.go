package stg

import "github.com/mebyus/ku/goku/compiler/sm"

type Cast struct {
	Exp Exp

	Pin sm.Pin

	// Type to which expression is being cast.
	typ *Type
}

func (c *Cast) Type() *Type {
	return c.typ
}

func (c *Cast) Span() sm.Span {
	return sm.Span{Pin: c.Pin}
}

func (c *Cast) String() string {
	panic("not implemented")
}
