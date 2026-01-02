package stg

import (
	"strconv"

	"github.com/mebyus/ku/goku/compiler/sm"
)

// Integer represents an integer constant (directly from source or evaluated)
// which value is known at compile time.
type Integer struct {
	Pin sm.Pin

	Val uint64

	typ *Type

	Neg bool
}

func (n *Integer) Type() *Type {
	return n.typ
}

func (n *Integer) Span() sm.Span {
	return sm.Span{Pin: n.Pin}
}

func (n *Integer) String() string {
	if n.Val == 0 {
		return "0"
	}

	if n.Val < 16 {
		s := strconv.FormatUint(n.Val, 10)
		if n.Neg {
			return "-" + s
		}
		return s
	}

	s := strconv.FormatUint(n.Val, 16)
	if n.Neg {
		return "-0x" + s
	}
	return "0x" + s
}

// Explicit interface implementation check.
var _ Exp = &Integer{}

// WithPin clones an Integer with specified pin position.
func (n *Integer) WithPin(pin sm.Pin) *Integer {
	v := *n
	v.Pin = pin
	return &v
}

// Add creates a new Integer of the same type with value n + v.
// Pin is not set for it.
func (n *Integer) Add(v uint64) *Integer {
	if !n.Neg {
		return &Integer{
			Val: n.Val + v,
			typ: n.typ,
		}
	}

	if n.Val > v {
		return &Integer{
			Neg: true,
			Val: n.Val - v,
			typ: n.typ,
		}
	}

	val := v - n.Val
	if val == 0 {
		return &Integer{typ: n.typ}
	}

	return &Integer{
		Val: val,
		typ: n.typ,
	}
}

// MakeInteger create non-negative static unsized integer.
func (x *TypeIndex) MakeInteger(pin sm.Pin, v uint64) *Integer {
	return &Integer{
		Pin: pin,
		Val: v,
		typ: x.Static.Integer,
	}
}

func (x *TypeIndex) MakeNegInteger(pin sm.Pin, v uint64) *Integer {
	n := x.MakeInteger(pin, v)
	if v != 0 {
		n.Neg = true
	}
	return n
}
