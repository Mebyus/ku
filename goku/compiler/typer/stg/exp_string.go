package stg

import "github.com/mebyus/ku/goku/compiler/srcmap"

// String represents a string constant (directly from source or evaluated)
// which value is known at compile time.
type String struct {
	Pin srcmap.Pin

	Val string

	typ *Type
}

func (s String) Type() *Type {
	return s.typ
}

func (s String) Span() srcmap.Span {
	return srcmap.Span{Pin: s.Pin}
}

func (s String) String() string {
	if s.Val == "" {
		return `""`
	}

	return "\"" + s.Val + "\""
}

// Explicit interface implementation check.
var _ Exp = String{}

// MakeInteger create non-negative static unsized integer.
func (x *TypeIndex) MakeString(pin srcmap.Pin, v string) String {
	return String{
		Pin: pin,
		Val: v,
		typ: x.Static.String,
	}
}
