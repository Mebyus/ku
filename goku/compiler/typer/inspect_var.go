package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) inspectVarSymbol(s *stg.Symbol) diag.Error {
	v := t.box.Const(s.Aux)

	err := t.inspectType(v.Type)
	if err != nil {
		return err
	}

	err = t.inspectExp(v.Exp)
	if err != nil {
		return err
	}

	_, ok := t.ins.links[s]
	if ok {
		return &diag.SimpleMessageError{
			Pin:  s.Pin,
			Text: fmt.Sprintf("variable \"%s\" definition references itself", s.Name),
		}
	}

	return nil
}
