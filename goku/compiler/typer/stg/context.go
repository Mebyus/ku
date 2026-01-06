package stg

import "github.com/mebyus/ku/goku/compiler/sm"

// Context contains type and symbol information about imported units.
// It also holds build conditions under which unit or program compilation is performed.
type Context struct {
	Global Scope
	Types  TypeIndex
	Gens   GenIndex
	Map    map[sm.UnitPath]*Unit
}

func (c *Context) Init() {
	c.Types.Init()
	c.Gens.Init()
	c.Global.InitGlobal(&c.Types, &c.Gens)

	if c.Map == nil {
		c.Map = make(map[sm.UnitPath]*Unit)
	}

	addBuiltinTypes(c)
	addBuiltinGens(c)
}
