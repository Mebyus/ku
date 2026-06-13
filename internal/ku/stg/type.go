package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/enums/typk"
)

type Type struct {
	symdef // can be used as symbol definition of type symbol

	// Byte size of this type's value. May be 0 for some types.
	// More specifically this field equals the stride between two
	// consecutive elements of this type inside an array.
	Size uint32

	Flags TypeFlag

	Kind typk.Kind
}

// TypeFlag bit flags for specifing additional type properties.
type TypeFlag uint16

const (
	// Static variant of the type.
	TypeStatic TypeFlag = 1 << iota

	// Type is a builtin.
	TypeBuiltin

	// Type has recursive definition.
	TypeRecursive

	// Signed integer type.
	TypeSigned
)

// applicable only for integer types
func (t *Type) IsSigned() bool {
	return t.Flags&TypeSigned != 0
}

func (t *Typer) LookupType(s *Scope, spec ast.TypeSpec) *Type {
	switch p := spec.(type) {
	// case ast.VoidPointer:
	// 	return s.Types.Known.VoidPointer, nil
	// case ast.VoidRef:
	// 	return s.Types.Known.VoidRef, nil
	// case ast.Void:
	// 	return s.Types.Known.Void, nil
	case *ast.TypeName:
		return t.lookupTypeName(s, p)
	default:
		panic(fmt.Sprintf("unexpected %T type specifier", p))
	}
}

func (t *Typer) lookupTypeName(s *Scope, p *ast.TypeName) *Type {
	name := p.Name
	pin := p.Pin

	symbol := s.Lookup(name)
	if symbol == nil {
		t.report(pin, fmt.Sprintf("unknown type \"%s\"", name))
		return t.common.Types.Invalid
	}
	if symbol.Kind != symk.Type {
		t.report(pin, fmt.Sprintf("expected type here, but name \"%s\" refers to %s symbol", name, symbol.Kind))
		return t.common.Types.Invalid
	}

	return symbol.Def.(*Type)
}

type TypeIndex struct {
	Static StaticTypes
	Known  KnownTypes

	Invalid *Type
}

func (x *TypeIndex) init() {
	x.Static.init()
	x.Known.init()
	x.Invalid = &Type{Kind: typk.Invalid}
}

// StaticTypes contains instances of various predefined (builtin) static types.
type StaticTypes struct {
	// For nil literal.
	// Nil *Type

	// Unsized.
	Integer *Type

	// String *Type

	Boolean *Type

	// Rune *Type
}

func (t *StaticTypes) init() {
	// t.Nil = &Type{
	// 	Flags: TypeFlagBuiltin | TypeFlagStatic,
	// 	Kind:  tpk.Nil,
	// }

	t.Integer = &Type{
		Size:  0, // unsized static integer can hold arbitrary large integer number
		Flags: TypeBuiltin | TypeSigned | TypeStatic,
		Kind:  typk.Integer,
	}

	t.Boolean = &Type{
		Flags: TypeBuiltin | TypeStatic,
		Kind:  typk.Boolean,
	}
}

// KnownTypes contains instances of various primitive runtime types and their derivatives
// (spans, pointers, etc.).
//
// Mostly used as shorthand access to types which compiler should be aware of.
type KnownTypes struct {
	// void, empty struct, zero size array
	// Void *Type

	// *void
	// VoidPointer *Type

	// &void
	// VoidRef *Type

	// u8
	U32 *Type

	S32 *Type

	Bool *Type
}

func (t *KnownTypes) init() {
	t.U32 = &Type{
		Size:  4,
		Flags: TypeBuiltin,
		Kind:  typk.Integer,
	}

	t.S32 = &Type{
		Size:  4,
		Flags: TypeBuiltin | TypeSigned,
		Kind:  typk.Integer,
	}

	t.Bool = &Type{
		Size:  1,
		Flags: TypeBuiltin,
		Kind:  typk.Boolean,
	}
}
