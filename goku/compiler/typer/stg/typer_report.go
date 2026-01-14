package stg

import "github.com/mebyus/ku/goku/compiler/diag"

func (t *Typer) report(msg diag.Error) {
	t.errors = append(t.errors, msg)
}
