package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) convTypeSymbol(s *stg.Symbol) diag.Error {
	spec := t.box.Type(s.Aux).Spec
	methods := t.methodsByReceiver[s]

	p, ok := spec.(ast.Struct)
	if ok {
		err := t.checkStructCollisions(s, p.Fields, methods)
		if err != nil {
			return err
		}
	}

	typ, err := t.ctx.Types.Lookup(&t.unit.Scope, spec)
	if err != nil {
		return err
	}

	custom := &stg.Custom{
		Symbol:  s,
		Methods: methods,
		Type:    typ,
	}
	def := &stg.Type{
		Def:  custom,
		Kind: tpk.Custom,
	}

	s.Def = def
	return nil
}

func (t *Typer) checkStructCollisions(s *stg.Symbol, fields []ast.Field, methods []*stg.Symbol) diag.Error {
	d := t.fields
	clear(d)

	for _, f := range fields {
		name := f.Name.Str
		pin := f.Name.Pin
		_, ok := d[name]
		if ok {
			return &diag.SimpleMessageError{
				Pin:  pin,
				Text: fmt.Sprintf("struct \"%s\" already has \"%s\" field", s.Name, name),
			}
		}

		d[name] = pin
	}

	for _, m := range methods {
		name := m.GetMethodName()
		pin := m.Pin

		_, ok := d[name]
		if ok {
			return &diag.SimpleMessageError{
				Pin:  pin,
				Text: fmt.Sprintf("struct \"%s\" already has \"%s\" field", s.Name, name),
			}
		}

		d[name] = pin
	}

	return nil
}
