package stg

import "github.com/mebyus/ku/internal/ku/sx"

func (t *Typer) report(pin sx.Pin, msg string) {
	t.unit.Errors = append(t.unit.Errors, &Error{
		Pin:   pin,
		Short: msg,
	})
}

// InvExp represents invalid expression. It can be syntactically malformed or
// semantically invalid expression.
//
// Can act as expression or operand.
type InvExp struct {
	operand

	pin sx.Pin

	typ *Type
}

// Explicit interface implementation check.
var _ Operand = &InvExp{}

func (e *InvExp) Type() *Type {
	return e.typ
}

func (e *InvExp) Pin() sx.Pin {
	return e.pin
}

func (t *Typer) makeInvExp(pin sx.Pin) *InvExp {
	return &InvExp{
		pin: pin,
		typ: t.common.Types.Invalid,
	}
}
