package hg

// DepSet is a hepler object for collecting dependencies between symbols.
type DepSet struct {
	// Contains set of symbols referenced by the symbol being scanned.
	m map[Symbol]struct{}

	// Symbol id which is currently being scanned.
	sym Symbol

	// Symbol has self-loop i.e. it directly depends on itself.
	loop bool
}

func (g *DepSet) Reset(sym Symbol) {
	if g.m == nil {
		g.m = make(map[Symbol]struct{})
	} else {
		clear(g.m)
	}

	g.loop = false
	g.sym = sym
}

// returns list of all stored dependencies (except self-loop).
func (g *DepSet) list() []Symbol {
	if len(g.m) == 0 {
		return nil
	}

	list := make([]Symbol, 0, len(g.m))
	for s := range g.m {
		list = append(list, s)
	}
	return list
}

// Add dependency between the given symbol and the one currently under inspection.
//
// Returns true if given symbol forms a self-loop.
func (g *DepSet) Add(sym Symbol) bool {
	if sym == g.sym {
		g.loop = true
		return true
	}

	g.m[sym] = struct{}{}
	return false
}
