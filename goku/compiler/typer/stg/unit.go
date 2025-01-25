package stg

import (
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
