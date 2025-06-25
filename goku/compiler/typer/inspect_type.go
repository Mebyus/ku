package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) inspectTypeSymbol(s *stg.Symbol) diag.Error {
	spec := t.box.Type(s.Aux).Spec

	var err diag.Error
	switch p := spec.(type) {
	case ast.Struct:
		fmt.Printf("struct %s %d %T\n", s.Name, s.Aux, spec)
		err = t.inspectStructFields(p.Fields)
	case ast.Bag:
		fmt.Printf("WARN: bag type specifier not implemented (%s %d %T)\n", s.Name, s.Aux, spec)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", spec.Kind(), spec.Kind(), spec))
	}
	return err
}

func (t *Typer) inspectStructFields(fields []ast.Field) diag.Error {
	return nil
}
