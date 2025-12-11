package stg

import "github.com/mebyus/ku/goku/compiler/srcmap/origin"

// Context contains type and symbol information about imported units.
// It also holds build conditions under which unit or program compilation is performed.
type Context struct {
	Global Scope
	Types  TypeIndex
	Map    map[origin.Path]*Unit
}

func (c *Context) Init() {
	c.Types.Init()
	c.Global.InitGlobal()

	if c.Map == nil {
		c.Map = make(map[origin.Path]*Unit)
	}

	addBuiltinTypes(c)
}
