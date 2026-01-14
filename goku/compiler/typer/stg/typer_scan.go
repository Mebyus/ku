package stg

import (
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
)

// DepSet is a hepler object for collecting dependencies of unit-level symbols
// on other unit-level symbols.
type DepSet struct {
	// Contains set of symbols referenced by the symbol being scanned.
	set map[*Symbol]struct{}
}

func (g *DepSet) init() {
	g.set = make(map[*Symbol]struct{})
}

func (g *DepSet) reset() {
	clear(g.set)
}

// Take returns list of all dependencies found during scan and resets
// internal state afterward. After calling this method DepSet can be
// used again for scanning another symbol.
func (g *DepSet) take() []*Symbol {
	if len(g.set) == 0 {
		return nil
	}

	list := make([]*Symbol, 0, len(g.set))
	for s := range g.set {
		list = append(list, s)
	}

	g.reset()
	return list
}

// Add dependency between the given symbol and the one currently under inspection.
func (g *DepSet) add(s *Symbol) {
	if s.Scope.Kind != sck.Unit {
		// do not link global symbols, they are defined implicitly
		// before everything else
		return
	}
	if s.Kind == smk.Import {
		// no need to keep track of links to import symbol
		// since we handle unit ranking during uwalk phase
		return
	}

	g.set[s] = struct{}{}
}

func (t *Typer) scanConsts() {
	for i, s := range t.consts {
		err := t.scanConst(i, s)
		if err != nil {
			t.report(err)
		}
	}
}

func (t *Typer) scanVars() {
	for i, s := range t.vars {
		err := t.scanVar(i, s)
		if err != nil {
			t.report(err)
		}
	}
}

func (t *Typer) scanTypes() {
	for i, s := range t.types {
		err := t.scanType(i, s)
		if err != nil {
			t.report(err)
		}
	}
}

func (t *Typer) scanConst(i int, s *Symbol) diag.Error {
	c := t.box.consts[i]
	_ = c
	return nil
}

func (t *Typer) scanVar(i int, s *Symbol) diag.Error {
	v := t.box.vars[i]
	_ = v
	return nil
}

func (t *Typer) scanType(i int, s *Symbol) diag.Error {
	typ := t.box.types[i]
	_ = typ
	return nil
}
