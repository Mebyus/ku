package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) checkAndConvertAST() diag.Error {
	graph := t.graph
	for _, i := range graph.Isolated {
		err := t.convSymbol(graph.Nodes[i].Symbol)
		if err != nil {
			return err
		}
	}

	for _, comp := range graph.Comps {
		for _, c := range comp.Cohorts {
			for _, k := range c {
				err := t.convSymbol(graph.Nodes[comp.V[k].Index].Symbol)
				if err != nil {
					return err
				}
			}
		}
	}

	return t.translateSymbols()
}

func (t *Typer) convSymbol(s *stg.Symbol) diag.Error {
	switch s.Kind {
	case smk.Const:
		return t.convConstSymbol(s)
	case smk.Fun:
		return t.convFunSymbol(s)
	case smk.Method:
		return t.convMethodSymbol(s)
	case smk.Type:
		return t.convTypeSymbol(s)
	case smk.Var:
		return t.convVarSymbol(s)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) symbol (%s)", s.Kind, s.Kind, s.Name))
	}
}
