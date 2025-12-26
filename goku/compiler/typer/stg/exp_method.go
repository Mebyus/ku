package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// BoundMethod represents usage of method bound to receiver expression.
// BoundMethod can be an operand or a part of call expression.
//
// Type of values for it is function type.
type BoundMethod struct {
	// Always not nil.
	Receiver Exp

	// Method being bound to receiver.
	Symbol *Symbol

	// Type of this expression. Must be a function type.
	typ *Type
}

func (m *BoundMethod) Type() *Type {
	return m.typ
}

func (m *BoundMethod) Span() sm.Span {
	return m.Receiver.Span()
}

func (m *BoundMethod) String() string {
	panic("not implemented")
}
