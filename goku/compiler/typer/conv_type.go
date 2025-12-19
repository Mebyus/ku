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

	typ, err := t.unit.Scope.LookupType(spec)
	if err != nil {
		return err
	}

	if typ.Kind == tpk.Boolean {
		return &diag.SimpleMessageError{
			Pin:  s.Pin,
			Text: "custom boolean types are forbidden",
		}
	}

	if typ.Kind == tpk.Custom {
		// If custom type A is created upon another custom type B
		// then we take base B's base type as A's base type.
		//
		//	A.Base = B.Base (not A.Base = B)
		//
		typ = typ.Def.(*stg.Custom).Type
		if typ.Kind == tpk.Custom {
			panic("custom type inside another custom type")
		}
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
