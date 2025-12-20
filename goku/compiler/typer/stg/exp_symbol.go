package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// VarExp represents expression formed from variable usage (as operand)
// inside an expression.
type VarExp struct {
	Pin sm.Pin

	Symbol *Symbol
}

func (v *VarExp) Type() *Type {
	return v.Symbol.Type
}

func (v *VarExp) Span() sm.Span {
	return sm.Span{Pin: v.Pin}
}

func (v *VarExp) String() string {
	return v.Symbol.Name
}
