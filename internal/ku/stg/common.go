package stg

import (
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Common contains things which (by their nature) exist as a single instance
// during whole program compilation.
//
// Contains global scope, containers for types and generics.
// Holds build conditions under which unit or program compilation is performed.
type Common struct {
	Types TypeIndex

	Global Scope

	Pool *sx.Pool
}

func (c *Common) Init(pool *sx.Pool) {
	c.Pool = pool
	c.Types.init()
	c.Global.InitGlobal()
	c.bindBuiltinSymbols()
}

func (c *Common) bindBuiltinSymbols() {
	c.bindBuiltinTypeSymbol("u32", c.Types.Known.U32)
	c.bindBuiltinTypeSymbol("bool", c.Types.Known.Bool)
}

func (c *Common) bindBuiltinTypeSymbol(name string, t *Type) {
	s := &Symbol{
		Name: name,
		Kind: symk.Type,
		// Flags: SymbolBuiltin,
	}
	c.bindType(s, t)
}

func (c *Common) bindType(s *Symbol, t *Type) {
	s.Def = t
	c.Global.Bind(s)
}
