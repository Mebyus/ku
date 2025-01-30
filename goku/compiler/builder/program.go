package builder

import (
	"github.com/mebyus/ku/goku/compiler/parser"
	"github.com/mebyus/ku/goku/compiler/source/origin"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

// ParserSet collection of parsers for selected unit source texts.
type ParserSet []*parser.Parser

// Bundle contains result of unit walking phase.
type Bundle struct {
	Graph Graph

	// List of all program units sorted by import path.
	Units []*stg.Unit

	// Index in this slice corresponds to Unit.DiscoveryIndex.
	// Every parser in this slice has only its header parsed.
	Source []ParserSet

	Map map[origin.Path]*stg.Unit

	// Not nil if bundle has main unit inside.
	Main *stg.Unit

	Global *stg.Scope
}

func (b *Bundle) makeGraph() *Cycle {
	b.mapGraphNodes()
	var s Scout
	return s.rankOrFindCycle(&b.Graph)
}
