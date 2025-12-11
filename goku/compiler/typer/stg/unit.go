package stg

import (
	"sort"

	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
)

type Unit struct {
	nodeSymDef

	// Unit top level scope.
	Scope Scope

	// Scope that holds all unit test symbols from all unit texts.
	//
	// This field is always not nil and Scope.Kind is always equal to sck.Test.
	TestScope Scope

	// Scope that holds all unit unsafe symbols from all unit texts.
	//
	// This field is always not nil and Scope.Kind is always equal to sck.Unsafe.
	UnsafeScope Scope

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
	Pin srcmap.Pin
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

func (u *Unit) FindImportSite(path origin.Path) (ImportSite, bool) {
	for _, s := range u.Imports {
		if s.Path == path {
			return s, true
		}
	}
	return ImportSite{}, false
}

func (u *Unit) InitScopes(global *Scope) {
	u.Scope.Init(sck.Unit, global)
	u.TestScope.Init(sck.Test, &u.Scope)
	u.UnsafeScope.Init(sck.Unsafe, &u.Scope)
}
