package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) inspectConstSymbol(s *stg.Symbol) diag.Error {
	c := t.box.Const(s.Aux)

	err := t.inspectConstType(c.Type)
	if err != nil {
		return err
	}

	return t.inspectConstExp(c.Exp)
}

func (t *Typer) inspectConstType(spec ast.TypeSpec) diag.Error {
	switch p := spec.(type) {
	case nil:
		// untyped constant
		return nil
	case ast.TypeName:
		return t.linkTypeName(p)
	case ast.TypeFullName:
		return t.inspectTypeFullName(p)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}
}
