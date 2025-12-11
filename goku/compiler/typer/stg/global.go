package stg

import (
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
)

func (s *Scope) InitGlobal() {
	s.Init(sck.Global, nil)
}

func addBuiltinTypes(c *Context) {
	addUnsignedIntegerType(c, "u8", 1)
	addUnsignedIntegerType(c, "u16", 2)
	addUnsignedIntegerType(c, "u32", 4)
	addUnsignedIntegerType(c, "u64", 8)
	addUnsignedIntegerType(c, "u128", 16)

	addSignedIntegerType(c, "s8", 1)
	addSignedIntegerType(c, "s16", 2)
	addSignedIntegerType(c, "s32", 4)
	addSignedIntegerType(c, "s64", 8)
	addSignedIntegerType(c, "s128", 16)

	addBoolType(c)
}

func addUnsignedIntegerType(c *Context, name string, size uint32) {
	s := &Symbol{
		Name:  name,
		Kind:  smk.Type,
		Flags: SymbolBuiltin,
	}
	t := &Type{
		Size:  size,
		Flags: TypeFlagBuiltin,
		Kind:  tpk.Integer,
	}
	c.addType(s, t)
}

func addSignedIntegerType(c *Context, name string, size uint32) {
	s := &Symbol{
		Name:  name,
		Kind:  smk.Type,
		Flags: SymbolBuiltin,
	}
	t := &Type{
		Size:  size,
		Flags: TypeFlagBuiltin | TypeFlagSigned,
		Kind:  tpk.Integer,
	}
	c.addType(s, t)
}

func addBoolType(c *Context) {
	s := &Symbol{
		Name:  "bool",
		Kind:  smk.Type,
		Flags: SymbolBuiltin,
	}
	t := &Type{
		Size:  1,
		Flags: TypeFlagBuiltin,
		Kind:  tpk.Boolean,
	}
	c.addType(s, t)
}

func (c *Context) addType(s *Symbol, t *Type) {
	s.Def = t
	c.Global.Bind(s)
}
