package stg

import (
	"github.com/mebyus/ku/goku/compiler/enums/bgk"
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
)

func (s *Scope) InitGlobal(types *TypeIndex, gens *GenIndex) {
	s.init(sck.Global)
	s.Types = types
	s.Gens = gens
}

const archPointerSize = 8

func addBuiltinTypes(c *Context) {
	addBuiltinType(c, "u8", c.Types.Known.U8)

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

	addBuiltinType(c, "bool", c.Types.Known.Bool)
	addUnsignedIntegerType(c, "rune", 4)

	addBuiltinType(c, "errid", c.Types.Known.ErrId)
	addBuiltinType(c, "error", c.Types.Known.Error)

	addFloatType(c, "f32", 4)
	addFloatType(c, "f64", 8)
	addFloatType(c, "f128", 16)

	addBuiltinType(c, "str", c.Types.Known.Str)
}

func addBuiltinType(c *Context, name string, t *Type) {
	s := &Symbol{
		Name:  name,
		Kind:  smk.Type,
		Flags: SymbolBuiltin,
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

func (c *Context) addGenFun(name string, kind bgk.Kind) {
	s := &Symbol{
		Name:  name,
		Kind:  smk.BgenFun,
		Flags: SymbolBuiltin,
		Aux:   uint32(kind),
	}
	c.Global.Bind(s)
}

func addBuiltinGens(c *Context) {
	c.addGenFun("min", bgk.Min)
	c.addGenFun("max", bgk.Max)
	c.addGenFun("copy", bgk.Copy)
	c.addGenFun("clear", bgk.Clear)
}
