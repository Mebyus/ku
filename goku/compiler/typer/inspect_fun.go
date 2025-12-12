package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) inspectFunSymbol(s *stg.Symbol) diag.Error {
	fun := t.box.Fun(s.Aux)

	err := t.inspectSignature(fun.Signature)
	if err != nil {
		return err
	}

	return t.inspectBody(fun.Body)
}

func (t *Typer) inspectMethodSymbol(s *stg.Symbol) diag.Error {
	m := t.box.Method(s.Aux)

	err := t.inspectSignature(m.Signature)
	if err != nil {
		return err
	}

	return nil
}

func (t *Typer) inspectSignature(sig ast.Signature) diag.Error {
	err := t.inspectResultType(sig.Result)
	if err != nil {
		return err
	}

	return t.inspectParams(sig.Params)
}

func (t *Typer) inspectResultType(spec ast.TypeSpec) diag.Error {
	switch p := spec.(type) {
	case nil:
		// function returns nothing or never returns
		return nil
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}
}

func (t *Typer) inspectParams(params []ast.Param) diag.Error {
	for _, p := range params {
		err := t.inspectParam(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) inspectParam(param ast.Param) diag.Error {
	switch p := param.Type.(type) {
	case ast.TypeName:
		return t.linkTypeName(p)
	case ast.TypeFullName:
		return t.inspectTypeFullName(p)
	case ast.Ref:
		return t.linkRef(p)
	case ast.Pointer:
		return t.linkPointer(p)
	case ast.ArrayPointer:
		return t.linkArrayPointer(p)
	case ast.AnyPointer:
		return nil
	case ast.Array:
		return t.linkArray(p)
	case ast.Chunk:
		return t.linkChunk(p)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}
}

func (t *Typer) inspectBody(body ast.Block) diag.Error {
	for _, s := range body.Nodes {
		err := t.inspectStatement(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) inspectStatement(stm ast.Statement) diag.Error {
	return nil
}
