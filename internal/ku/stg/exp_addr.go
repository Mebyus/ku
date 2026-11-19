package stg

import "github.com/mebyus/ku/internal/ku/sx"

// TakeRef represents expression of taking a reference (address)
// of another expression.
type TakeRef struct {
	operand

	// Always not nil
	Exp Exp

	pin sx.Pin

	typ *Type
}

var _ Operand = &TakeRef{}

func (r *TakeRef) Type() *Type {
	return r.typ
}

func (r *TakeRef) Pin() sx.Pin {
	return r.pin
}

func (r *TakeRef) Target() bool {
	return false
}
