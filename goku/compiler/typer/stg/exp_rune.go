package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// Rune represents a rune (integer) constant (directly from source or evaluated)
// which value is known at compile time.
type Rune struct {
	Pin sm.Pin

	typ *Type

	Val uint32
}

func (r *Rune) Type() *Type {
	return r.typ
}

func (r *Rune) Span() sm.Span {
	return sm.Span{Pin: r.Pin}
}

func (r *Rune) String() string {
	return "'" + string([]rune{rune(r.Val)}) + "'"
}

func (x *TypeIndex) MakeRune(pin sm.Pin, v uint32) *Rune {
	return &Rune{
		typ: x.Static.Rune,
		Pin: pin,
		Val: v,
	}
}
