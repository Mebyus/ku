package stg

import (
	"sort"

	"github.com/mebyus/ku/goku/compiler/source"
	"github.com/mebyus/ku/goku/compiler/source/origin"
)

type Unit struct {
	// Unit top level scope.
	Scope Scope

	// Unit path of this unit.
	Path origin.Path

	// All imports inside this unit.
	Imports []ImportSite

	// Unit index assigned by order in which units are discovered
	// during unit discovery phase (uwalk).
	DiscoveryIndex uint32

	// Unit index assigned after path sorting.
	Index uint32
}

// ImportSite represents a single unit import inside an import block.
type ImportSite struct {
	// Unit path of imported unit.
	Path origin.Path

	// Unit is imported under this name.
	Name string

	// Place where import occurs in source code.
	Pin source.Pin
}

func SortAndOrderUnits(units []*Unit) {
	if len(units) == 0 {
		panic("invalid argument: <nil>")
	}

	if len(units) == 1 {
		return
	}

	u := units
	sort.Slice(u, func(i, j int) bool {
		a := u[i]
		b := u[j]
		return origin.Less(a.Path, b.Path)
	})

	for i := range len(u) {
		u[i].Index = uint32(i)
	}
}

func SortImports(ss []ImportSite) {
	if len(ss) < 2 {
		return
	}

	sort.Slice(ss, func(i, j int) bool {
		a := ss[i]
		b := ss[j]
		return origin.Less(a.Path, b.Path)
	})
}

func (u *Unit) HasMain() bool {
	return false
}
