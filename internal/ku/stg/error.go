package stg

import "github.com/mebyus/ku/internal/ku/sx"

func (t *Typer) report(pin sx.Pin, msg string) {
	t.unit.Errors = append(t.unit.Errors, &Error{
		Pin:   pin,
		Short: msg,
	})
}
