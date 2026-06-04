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

	Pin sx.Pin
}
