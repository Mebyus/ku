package stg

import (
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
)

func (s *Scope) InitGlobal(types *TypeIndex) {
	s.Init(sck.Global, nil)
	s.Types = types
}

const archPointerSize = 8

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

	addBuiltinType(c, "uint", c.Types.Known.Uint)
	addSignedIntegerType(c, "sint", archPointerSize)

	addBoolType(c)
	addUnsignedIntegerType(c, "rune", 4)
	addUnsignedIntegerType(c, "error_id", archPointerSize)

	addFloatType(c, "f32", 4)
	addFloatType(c, "f64", 8)
	addFloatType(c, "f128", 16)

	addStringType(c)
}

func addBuiltinType(c *Context, name string, t *Type) {
	s := &Symbol{
		Name:  name,
		Kind:  smk.Type,
		Flags: SymbolBuiltin,
	}
	c.addType(s, t)
}

func addStringType(c *Context) {
	s := &Symbol{
		Name:  "str",
		Kind:  smk.Type,
		Flags: SymbolBuiltin,
	}
	t := &Type{
		Size:  2 * archPointerSize,
		Flags: TypeFlagBuiltin,
		Kind:  tpk.String,
	}
	c.addType(s, t)
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
	c.addType(s, c.Types.Known.Bool)
}

func addFloatType(c *Context, name string, size uint32) {
	s := &Symbol{
		Name:  name,
		Kind:  smk.Type,
		Flags: SymbolBuiltin,
	}
	t := &Type{
		Size:  size,
		Flags: TypeFlagBuiltin,
		Kind:  tpk.Float,
	}
	c.addType(s, t)
}

func (c *Context) addType(s *Symbol, t *Type) {
	s.Def = SymDefType{Type: t}
	c.Global.Bind(s)
}
