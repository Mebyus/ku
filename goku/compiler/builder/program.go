package builder

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/parser"
	"github.com/mebyus/ku/goku/compiler/typer"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

// ParserSet collection of parsers for selected unit source texts.
type ParserSet []*parser.Parser

// Bundle contains result of unit walking phase.
type Bundle struct {
	Graph Graph

	Context stg.Context

	// List of all program units sorted by import path.
	Units []*stg.Unit

	// Index in this slice corresponds to Unit.DiscoveryIndex.
	// Every parser in this slice has only its header parsed.
	Source []ParserSet

	// Not nil if bundle has main unit inside.
	Main *stg.Unit
}

func (b *Bundle) GetUnitParsers(unit *stg.Unit) ParserSet {
	return b.Source[unit.DiscoveryIndex]
}

func (b *Bundle) makeGraph() *Cycle {
	b.mapGraphNodes()
	var s Scout
	return s.rankOrFindCycle(&b.Graph)
}

func CompileBundle(b *Bundle) diag.Error {
	b.Context.Init()

	for _, cohort := range b.Graph.Cohorts {
		for _, i := range cohort {
			unit := b.Units[i]
			texts, err := ParseTexts(b.GetUnitParsers(unit))
			if err != nil {
				return err
			}
			err = typer.Compile(&b.Context, unit, texts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ParseTexts(s ParserSet) ([]*ast.Text, diag.Error) {
	texts := make([]*ast.Text, 0, len(s))
	for _, p := range s {
		t, err := p.Nodes()
		if err != nil {
			return nil, err
		}
		texts = append(texts, t)
	}
	return texts, nil
}
