package stg

import (
	"fmt"
	"strings"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/scok"
	"github.com/mebyus/ku/internal/ku/enums/symk"
)

// Implements hoisting for unit-level constants and custom types.
//
// Constants are hoisted first. Then resulting constant values
// can be used (if they are needed) for expanding and hoisting
// custom types.
//
// Constant hoisting algorithm contains roughly 4 steps:
//
//  1. gather dependencies between constants
//  2. construct dependency graph for constants
//  3. analyze graph structure to rank constant definition order
//  4. evaluate and assign types for constants on each rank
//
// Custom type hoisting algorithm:
//
//  1. resolve and expand constant usage in type definitions (for arrays)
//  2. gather dependencies between custom types
//  3. construct dependency graph for custom types
//  4. analyze graph structure to detect components and clusters
//  5. check each cluster for recursive validity
//  6. rank custom types
//  7. resolve type specifiers for each rank
func (t *Typer) expand() {
	t.gatherConstDeps()
	t.hoistConsts()
}

// Inspect AST and gather dependencies between unit-level constants.
func (t *Typer) gatherConstDeps() {
	t.graph.Reset(len(t.unit.Consts))
	for _, s := range t.unit.Consts {
		def := &t.box.consts[s.Aux]
		t.scanConstDeps(s, def)
	}
}

// Hoist and rank usage of unit-level constants based on dependencies
// between them.
func (t *Typer) hoistConsts() {
	t.graph.Build()

	const debug = false

	isolated := t.graph.Isolated()
	components := t.graph.Components()

	if debug {
		symbols := t.graph.Symbols()
		if len(symbols) != 0 {
			fmt.Printf("[chst] constant symbols (n=%d): %s\n", len(symbols), strings.Join(getSymbolNames(symbols), ", "))
		}

		if len(isolated) != 0 {
			fmt.Printf("[chst] found %d isolated constants:\n", len(isolated))
			fmt.Printf("     | %s\n", strings.Join(getSymbolNames(isolated), ", "))
		}

		if len(components) != 0 {
			fmt.Printf("[chst] found %d components:\n", len(components))
			for i := range len(components) {
				c := &components[i]
				fmt.Printf("     | %s\n", strings.Join(getSymbolNames(c.Symbols()), ", "))
			}
		}
	}

	for i := range len(components) {
		c := &components[i]
		if c.HasClusters() {
			// TODO: report with pin
			t.report(0, "constants form a cluster")
			continue
		}

		t.rank.Reset(c)
		tower := t.rank.Rank()

		if debug {
			fmt.Printf("[chst] tower %d:\n", i)

			for _, c := range tower.Cohorts {
				fmt.Printf("     | %s\n", strings.Join(getSymbolNames(c.Symbols), ", "))
			}
		}
	}
}

// Perform deep scan of unit-level dependencies inside
// constant defenition expression.
func (t *Typer) scanConstDeps(s *Symbol, def *ast.Const) {
	t.deps.Reset(s.RawSym())

	if def.Type != nil {
		panic("stub")
	}
	t.scanConstExpDeps(def.Exp)

	t.graph.Take(&t.deps)
}

func (t *Typer) scanConstExpDeps(exp ast.Exp) {
	switch e := exp.(type) {
	case *ast.InvExp:
		// already reported during parsing
	case *ast.Integer, *ast.False, *ast.True,
		*ast.String, *ast.Rune:
		// literals do not cause symbol dependencies
	case *ast.SymExp:
		t.scanConstSymExpDeps(e)
	case *ast.ParenExp:
		t.scanConstExpDeps(e.Exp)
	case *ast.UnExp:
		t.scanConstExpDeps(e.A)
	case *ast.BinExp:
		t.scanConstExpDeps(e.A)
		t.scanConstExpDeps(e.B)
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}

func (t *Typer) scanConstSymExpDeps(exp *ast.SymExp) {
	name := exp.Name
	pin := exp.Pin

	s := t.unit.Scope.Lookup(name)
	if s == nil {
		t.report(pin, fmt.Sprintf("unresolved symbol \"%s\" used as expression", name))
		return
	}
	if s.Kind != symk.Const {
		t.report(pin, fmt.Sprintf("symbol \"%s\" is not a constant, but used in constant expression", name))
		return
	}

	if s.Scope.Kind != scok.Unit {
		// we are gathering dependencies only between unit-level
		// symbols in current unit
		return
	}

	loop := t.deps.Add(s.RawSym())
	if loop {
		t.report(pin, fmt.Sprintf("definition for constant \"%s\" references itself", name))
	}
}
