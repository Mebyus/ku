package builder

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/parser"
	"github.com/mebyus/ku/goku/compiler/sm"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
	"github.com/mebyus/ku/goku/graphs"
)

// ParserSet collection of parsers for selected unit source texts.
type ParserSet []*parser.Parser

// Bundle contains result of unit walking phase.
type Bundle struct {
	Graph graphs.Graph

	Common stg.Common

	Prelude PreludeUnits

	// List of all program units sorted by import path.
	Units []*stg.Unit

	// Index in this slice corresponds to Unit.DiscoveryIndex.
	// Every parser in this slice has only its header parsed.
	Source []ParserSet

	// Not nil if bundle has main unit inside.
	Main *stg.Unit

	// Contains all source files discovered during uwalk phase.
	Pool *sm.Pool
}

type PreludeUnits struct {
	//	"fmt"
	Format *stg.Unit

	//	"mem"
	Memory *stg.Unit

	// Testing framework unit from standard library.
	//
	//	"stf"
	Test *stg.Unit
}

func (b *Bundle) GetUnitParsers(unit *stg.Unit) ParserSet {
	return b.Source[unit.DiscoveryIndex]
}

func (b *Bundle) makeGraph() *graphs.Cycle {
	b.mapGraphNodes()
	var s graphs.Scout
	return s.RankOrFindCycle(&b.Graph)
}

func CompileBundle(b *Bundle) diag.Error {
	b.Common.Init(b.Pool)
	p := stg.NewPool(&b.Common)

	exportCount := 0
	testCount := 0
	for _, cohort := range b.Graph.Cohorts {
		for _, i := range cohort {
			unit := b.Units[i]
			texts, err := ParseTexts(b.GetUnitParsers(unit))
			if err != nil {
				return err
			}
			typer := p.Get(unit)
			err = typer.Translate(texts)
			if err != nil {
				return err
			}

			exportCount += len(unit.Export)
			testCount += len(unit.Tests)
		}
	}

	if testCount != 0 {
		fmt.Printf("found %d test function(s)\n", testCount)
	}

	if exportCount == 0 {
		fmt.Println("stop: no exported symbols found")
		return nil
	}

	fmt.Printf("found %d exported symbol(s)\n", exportCount)
	exported := make([]*stg.Symbol, 0, exportCount)
	for _, u := range b.Units {
		exported = append(exported, u.Export...)
	}

	stg.Prune(exported)

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

// Fills Unit.Imports.Units according to import paths.
func (b *Bundle) mapGraphNodes() {
	m := make(map[sm.UnitPath]*stg.Unit, len(b.Units))
	b.Graph.Nodes = make([]graphs.Node, len(b.Units))
	b.Graph.Rank = make([]uint32, len(b.Units))
	b.Common.Map = m

	for _, unit := range b.Units {
		m[unit.Path] = unit
	}

	for i, unit := range b.Units {
		// i = unit.Index inside this loop, because we sorted
		// and indexed units beforehand

		b.Graph.Nodes[i].Anc = make([]uint32, 0, len(unit.Imports))
		for _, s := range unit.Imports {
			u, ok := m[s.Path]
			if !ok {
				panic(fmt.Sprintf("imported unit \"%s\" not found", s.Path))
			}
			if u == unit {
				panic("unit imported itself")
			}

			b.Graph.Nodes[i].AddAnc(u.Index)
			b.Graph.Nodes[u.Index].AddDes(uint32(i))
		}

		if len(unit.Imports) == 0 {
			b.Graph.Roots = append(b.Graph.Roots, uint32(i))
		}
	}
}

func convertImportCycle(c *graphs.Cycle, units []*stg.Unit) []sm.ImportSite {
	if len(c.Nodes) < 2 {
		panic("bad cycle data")
	}

	sites := make([]sm.ImportSite, 0, len(c.Nodes))
	for i := 0; i < len(c.Nodes)-1; i += 1 {
		j := c.Nodes[i]
		k := c.Nodes[i+1]

		u := units[j]
		next := units[k]

		s, ok := u.FindImportSite(next.Path)
		if !ok {
			panic(fmt.Sprintf("unable to find \"%s\" import inside \"%s\"", next.Path, u.Path))
		}

		sites = append(sites, s)
	}
	return sites
}
