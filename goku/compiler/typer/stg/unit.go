package stg

import (
	"sort"

	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
)

type Unit struct {
	nodeSymDef // TODO: move this to separate wrapper struct

	// Unit top level scope.
	Scope Scope

	// Scope that holds all unit test symbols from all unit texts.
	//
	// This field is always not nil and Scope.Kind is always equal to sck.Test.
	TestScope Scope

	// Unit path of this unit.
	Path origin.Path

	// All imports inside this unit.
	Imports []srcmap.ImportSite

	// Unit index assigned by order in which units are discovered
	// during unit discovery phase (uwalk).
	DiscoveryIndex uint32

	// Unit index assigned after path sorting.
	Index uint32
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

func SortImports(ss []srcmap.ImportSite) {
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
	// TODO: we can add this check in symbol indexing phase and store result in flags
	s := u.Scope.Get("main")
	if s == nil {
		return false
	}
	return s.Kind == smk.Fun
}

func (u *Unit) FindImportSite(path origin.Path) (srcmap.ImportSite, bool) {
	for _, s := range u.Imports {
		if s.Path == path {
			return s, true
		}
	}
	return srcmap.ImportSite{}, false
}

func (u *Unit) InitScopes(global *Scope) {
	u.Scope.Init(sck.Unit, global)
	u.TestScope.Init(sck.Test, &u.Scope)
}
