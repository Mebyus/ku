package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// SymExp represents expression formed from symbol usage (as operand or part of chain)
// inside an expression.
type SymExp struct {
	Pin sm.Pin

	Symbol *Symbol
}

func (v *SymExp) Type() *Type {
	return v.Symbol.Type
}

func (v *SymExp) Span() sm.Span {
	return sm.Span{Pin: v.Pin}
}

func (v *SymExp) String() string {
	return v.Symbol.Name
}
