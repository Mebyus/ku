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

// BoolExp represents a static boolean value with runtime expression inside.
//
// Expression needs to be evaluated first before boolean value is used.
// Note that runtime expression evaluation result if discarded.
//
// Although description of this node may sound strange it has its place
// in real programs. Consider this piece of code for example:
//
//	x := ...
//	...
//	const a := 1;
//	b := check(x) || a == 1;
//
// Here we can easily deduce that initial expression for variable b always
// evaluates to true at compile-time. But we must preserve check(x) call
// regardless for it may contain side-effects. Thus the whole expression
// becomes BoolExp with true value baked-in, but it also contains runtime
// expression check(x) inside.
//
// Type of this expression is runtime bool. We cannot make this type static
// in order to avoid possible compile-time simplifications in subsequent
// expressions which will discard internal expression completely.
type BoolExp struct {
	// This expression is evaluated in runtime, but its result is discarded.
	Exp Exp

	Pin sm.Pin

	// Always a runtime bool.
	typ *Type

	// Perceived value of this expression.
	Val bool
}

func (e *BoolExp) Type() *Type {
	return e.typ
}

func (e *BoolExp) Span() sm.Span {
	return sm.Span{Pin: e.Pin}
}

func (e *BoolExp) String() string {
	panic("not implemented")
}

// MakeBoolean create static boolean value.
func (x *TypeIndex) MakeBoolExp(pin sm.Pin, exp Exp, v bool) *BoolExp {
	return &BoolExp{
		Exp: exp,
		Pin: pin,
		Val: v,
		typ: x.Known.Bool,
	}
}
