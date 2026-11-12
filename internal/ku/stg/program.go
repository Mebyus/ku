package stg

import (
	"strconv"

	"github.com/mebyus/ku/internal/ku/enums/scok"
	"github.com/mebyus/ku/internal/ku/sx"
)

type Program struct {
	Common

	Units []*Unit
}

type Unit struct {
	Scope Scope

	Funs []*Symbol

	Consts []*Symbol

	// List of function symbols that should be exported in produced object file.
	Exports []*Symbol

	// Unit path.
	Path sx.Path

	// System path of directory with unit source files.
	Dir string

	Name string

	// Used to generate unique names (during codegen phase) for symbols from this unit.
	LinkName string

	// Function (cannot be method) symbol marked as program entrypoint.
	// Each unit can have at most one such function.
	//
	// In practice most units do not have entrypoint symbol.
	//
	// When unit does not contain entrypoint this field equals nil.
	Entry *Symbol
}

type Import struct {
	symdef

	Unit *Unit
}

func (u *Unit) init(global *Scope) {
	u.Scope.Init(scok.Unit, global)
}

// AssignLinkNames generates and assigns unique link names for each unit in
// the given list. Generated link names are based on original unit names.
func AssignLinkNames(units []*Unit) {
	// maps unit name to number units with that name already encountered
	// during link name generation
	m := make(map[ /* unit name */ string]uint64, len(units))

	for _, u := range units {
		if u.Name == "" {
			u.LinkName = "ku"
			continue
		}

		k := m[u.Name]
		if k == 0 {
			u.LinkName = u.Name
		} else {
			u.LinkName = u.Name + strconv.FormatUint(k, 10)
			m[u.Name] = k + 1
		}
	}
}
