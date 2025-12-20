package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// Nil represents usage of "nil" literal as operand or expression.
type Nil struct {
	Pin sm.Pin

	typ *Type
}

func (n *Nil) Type() *Type {
	return n.typ
}

func (n *Nil) Span() sm.Span {
	return sm.Span{Pin: n.Pin, Len: 3}
}

func (*Nil) String() string {
	return "nil"
}

func (x *TypeIndex) MakeNil(pin sm.Pin) *Nil {
	return &Nil{Pin: pin, typ: x.Static.Nil}
}
