package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// Boolean represents a boolean constant true/false (directly from source or evaluated)
// which value is known at compile time.
type Boolean struct {
	Pin sm.Pin

	typ *Type

	Val bool
}

func (b *Boolean) Type() *Type {
	return b.typ
}

func (b *Boolean) Span() sm.Span {
	return sm.Span{Pin: b.Pin}
}

func (b *Boolean) String() string {
	if b.Val {
		return "true"
	}
	return "false"
}

// Explicit interface implementation check.
var _ Exp = &Boolean{}

// MakeBoolean create static boolean value.
func (x *TypeIndex) MakeBoolean(pin sm.Pin, v bool) *Boolean {
	return &Boolean{
		Pin: pin,
		Val: v,
		typ: x.Static.Boolean,
	}
}
