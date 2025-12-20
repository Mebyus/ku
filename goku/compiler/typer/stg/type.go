package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/enums/tpk"
)

type Type struct {
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

	// Zero (default) type value.
	// Zero Exp

	// Byte size of this type's value. May be 0 for some types.
	// More specifically this field equals the stride between two
	// consecutive elements of this type inside an array.
	Size uint32

	// Bit flags with additional type properties. Actual meaning may differ
	// upon Kind.
	Flags TypeFlag

	// Discriminator for type definition category.
	Kind tpk.Kind
}

// String use only for debugging.
func (t *Type) String() string {
	var s string
	switch t.Kind {
	case tpk.Void:
		return "void"
	case tpk.Integer:
		switch t.Size {
		case 0:
			// unsized static
			s = "int"
		case 1:
			s = "8"
		case 2:
			s = "16"
		case 4:
			s = "32"
		case 8:
			s = "64"
		case 16:
			s = "128"
		default:
			panic(fmt.Sprintf("unexpected size (=%d)", t.Size))
		}
		if t.Size != 0 {
			if t.IsSigned() {
				s = "s" + s
			} else {
				s = "u" + s
			}
		}
	case tpk.Boolean:
		s = "bool"
	case tpk.String:
		s = "str"
	case tpk.Float:
		switch t.Size {
		case 0:
			// unsized static
			s = "float"
		case 4:
			s = "f32"
		case 8:
			s = "f64"
		case 16:
			s = "f128"
		default:
			panic(fmt.Sprintf("unexpected size (=%d)", t.Size))
		}
	case tpk.VoidPointer:
		return "*void"
	case tpk.Custom:
		c := t.Def.(*Custom)
		s = c.Symbol.Name
	default:
		return fmt.Sprintf("???(%d)", t.Kind)
	}

	if t.IsStatic() {
		s = "<" + s + ">"
	}
	return s
}

// TypeHash is a pseudo-unique type identifier which depends purely on type definition.
// Value of type hash must not depend on runtime pointer values of specific types,
// symbols, order of type definitions or usages.
//
// TODO: do we really need type hash?

type TypeDef interface {
	Kind() tpk.Kind
}

// TypeFlag bit flags for specifing additional type properties.
type TypeFlag uint16

const (
	// Static variant of the type.
	TypeFlagStatic TypeFlag = 1 << iota

	// Type is a builtin.
	TypeFlagBuiltin

	// Type has recursive definition.
	TypeFlagRecursive

	// Signed integer type.
	TypeFlagSigned
)

func (t *Type) IsStatic() bool {
	return t.Flags&TypeFlagStatic != 0
}

func (t *Type) IsSigned() bool {
	return t.Flags&TypeFlagSigned != 0
}

// Custom defines custom type.
type Custom struct {
	// List of methods which are bound to this custom type.
	Methods []*Symbol

	// Custom type is bound to this symbol.
	Symbol *Symbol

	// Type which was used to define this custom type.
	Type *Type
}

// Explicit interface implementation check.
var _ TypeDef = Custom{}

func (Custom) Kind() tpk.Kind {
	return tpk.Custom
}

type Pointer struct {
	// Type referenced by pointer.
	Type *Type
}

var _ TypeDef = Pointer{}

func (Pointer) Kind() tpk.Kind {
	return tpk.Pointer
}

type ArrayPointer struct {
	// Type referenced by pointer.
	Type *Type
}

var _ TypeDef = ArrayPointer{}

func (ArrayPointer) Kind() tpk.Kind {
	return tpk.ArrayPointer
}

type Ref struct {
	// Type referenced by reference.
	Type *Type
}

var _ TypeDef = Ref{}

func (Ref) Kind() tpk.Kind {
	return tpk.Ref
}

type Span struct {
	// Chunk element type.
	Type *Type
}

// Explicit interface implementation check.
var _ TypeDef = Span{}

func (Span) Kind() tpk.Kind {
	return tpk.Span
}

type Tuple struct {
	// Always not nil.
	Types []*Type
}

// Explicit interface implementation check.
var _ TypeDef = Tuple{}

func (Tuple) Kind() tpk.Kind {
	return tpk.Tuple
}

type Array struct {
	// Array element type.
	Type *Type

	// Number of elements in array.
	Len uint32
}

// Explicit interface implementation check.
var _ TypeDef = Array{}

func (Array) Kind() tpk.Kind {
	return tpk.Array
}

type Field struct {
	Name   string
	Type   *Type
	Offset uint32
}

type Struct struct {
	Fields []Field
}

func (Struct) Kind() tpk.Kind {
	return tpk.Struct
}
