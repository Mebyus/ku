package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

type LinkKind uint8

const (
	LinkDirect LinkKind = iota
	LinkIndirect
)

type Link struct {
	Symbol *stg.Symbol
	Kind   LinkKind
}

type Inspector struct {
	// contains symbols referenced by the symbol being inspected
	links map[*stg.Symbol]LinkKind

	// keeps track of current link type kind when descending/ascending
	// nested type specifiers
	//
	// only used during type inspection
	kind LinkKind
}

func (p *Inspector) Init() {
	p.links = make(map[*stg.Symbol]LinkKind)
}

func (p *Inspector) Reset() {
	clear(p.links)
	p.kind = LinkDirect
}

func (p *Inspector) Links() []Link {
	if len(p.links) == 0 {
		return nil
	}

	links := make([]Link, 0, len(p.links))
	for s, k := range p.links {
		links = append(links, Link{
			Symbol: s,
			Kind:   k,
		})
	}
	// TODO: sort links
	return links
}

// Add link between the given symbol and the one currently under inspection.
func (p *Inspector) link(s *stg.Symbol) {
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

	k, ok := p.links[s]
	if ok && k == LinkDirect {
		// direct link is stronger than indirect one
		return
	}

	p.links[s] = p.kind
}

// Change link nature to indirect and return previous state.
func (p *Inspector) indirect() LinkKind {
	k := p.kind
	p.kind = LinkIndirect
	return k
}

func (p *Inspector) restore(k LinkKind) {
	p.kind = k
}

func (t *Typer) inspectSymbol(s *stg.Symbol) diag.Error {
	// TODO: current fnction body inspect algorithm cannot properly
	// detect usage of method symbols

	k := s.Kind
	switch k {
	case smk.Fun:
		return t.inspectFunSymbol(s)
	case smk.Type:
		return t.inspectTypeSymbol(s)
	case smk.Alias:
		return t.inspectAliasSymbol(s)
	case smk.Method:
		return t.inspectMethodSymbol(s)
	case smk.Gen:
		return t.inspectGenSymbol(s)
	case smk.Const:
		return t.inspectConstSymbol(s)
	case smk.Var:
		return t.inspectVarSymbol(s)
	case smk.Test:
		return nil
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) symbol kind", k, k))
	}
}
