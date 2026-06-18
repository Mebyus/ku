package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/enums/typk"
)

type Type struct {
	symdef // can be used as symbol definition of type symbol

	// For some types this field is nil, since all necessary properties
	// are stored in other fields.
	//
	//	- void
	//	- integers
	//	- floats
	//	- strings
	//	- boolean
	//	- rune
	//	- *void
	Def TypeDef

	// Byte size of this type's value. May be 0 for some types.
	// More specifically this field equals the stride between two
	// consecutive elements of this type inside an array.
	Size uint32

	Flags TypeFlag

	Kind typk.Kind
}

func (t *Type) IsInvalid() bool {
	return t.Kind == typk.Invalid
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

func (t *Type) IsStatic() bool {
	return t.Flags&TypeStatic != 0
}

func (t *Type) String() string {
	switch t.Kind {
	case typk.Span:
		return "[]" + t.Def.(*Span).Type.String()
	case typk.String:
		if t.IsStatic() {
			return "<str>"
		}
		return "str"
	case typk.Integer:
		if t.Size == 0 {
			return "<int>"
		}
		if t.IsSigned() {
			switch t.Size {
			case 1:
				return "s8"
			case 2:
				return "s16"
			case 4:
				return "s32"
			case 8:
				return "s64"
			default:
				panic(fmt.Sprintf("unexpected integer size %d", t.Size))
			}
		}
		switch t.Size {
		case 1:
			return "u8"
		case 2:
			return "u16"
		case 4:
			return "u32"
		case 8:
			return "u64"
		default:
			panic(fmt.Sprintf("unexpected integer size %d", t.Size))
		}
	case typk.Boolean:
		if t.IsStatic() {
			return "<bool>"
		}
		return "bool"
	default:
		panic(fmt.Sprintf("unexpected %d kind", t.Kind))
	}
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
	case *ast.Span:
		return t.lookupSpan(s, p)
	case *ast.InvType:
		return t.common.Types.Invalid
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

func (t *Typer) lookupSpan(s *Scope, p *ast.Span) *Type {
	typ := t.LookupType(s, p.Type)
	if typ.IsInvalid() {
		return typ
	}

	return t.common.Types.getSpan(typ)
}

type TypeIndex struct {
	Static StaticTypes
	Known  KnownTypes

	// Maps span element type to the corresponding span type.
	Spans map[ /* span element type */ *Type]*Type

	Invalid *Type

	ArchPointerSize uint32
}

func (x *TypeIndex) init(archPointerSize uint32) {
	x.ArchPointerSize = archPointerSize

	x.Static.init()
	x.Known.init(archPointerSize)

	x.Invalid = &Type{Kind: typk.Invalid}

	x.Spans = make(map[*Type]*Type)
}

// StaticTypes contains instances of various predefined (builtin) static types.
type StaticTypes struct {
	// For nil literal.
	// Nil *Type

	// Unsized.
	Integer *Type

	Boolean *Type

	String *Type

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

	t.String = &Type{
		Flags: TypeBuiltin | TypeStatic,
		Kind:  typk.String,
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

	Uint *Type

	U8  *Type
	U32 *Type

	S32 *Type

	Bool *Type

	Str *Type
}

func (t *KnownTypes) init(archPointerSize uint32) {
	t.Uint = &Type{
		Size:  archPointerSize,
		Flags: TypeBuiltin,
		Kind:  typk.Integer,
	}

	t.U8 = &Type{
		Size:  1,
		Flags: TypeBuiltin,
		Kind:  typk.Integer,
	}

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

	t.Str = &Type{
		Size:  2 * archPointerSize,
		Flags: TypeBuiltin,
		Kind:  typk.String,
	}
}

func (x *TypeIndex) getSpan(t *Type) *Type {
	typ, ok := x.Spans[t]
	if ok {
		return typ
	}
	typ = &Type{
		Def:  &Span{Type: t},
		Size: 2 * x.ArchPointerSize,
		Kind: typk.Span,
	}
	x.Spans[t] = typ
	return typ
}
