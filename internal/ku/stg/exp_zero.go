package stg

import "github.com/mebyus/ku/internal/ku/sx"

// Zero represents zero (default) value of a type.
type Zero struct {
	exp

	pin sx.Pin

	typ *Type
}

var _ Exp = &Zero{}

func (z *Zero) Type() *Type {
	return z.typ
}

func (z *Zero) Pin() sx.Pin {
	return z.pin
}

func (*Zero) Target() bool {
	return false
}
