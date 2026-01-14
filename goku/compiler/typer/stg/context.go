package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// Common contains things which (by their nature) exist as a single instance
// during whole program compilation.
//
// Contains global scope, containers for types and generics.
// Holds build conditions under which unit or program compilation is performed.
type Common struct {
	Global Scope
	Types  TypeIndex
	Gens   GenIndex
	Map    map[sm.UnitPath]*Unit

	pool *sm.Pool
}

func (c *Common) Init(pool *sm.Pool) {
	c.pool = pool

	c.Types.Init()
	c.Gens.Init()
	c.Global.InitGlobal(&c.Types, &c.Gens)

	if c.Map == nil {
		c.Map = make(map[sm.UnitPath]*Unit)
	}

	addBuiltinTypes(c)
	addBuiltinGens(c)
}
