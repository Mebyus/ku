package stg

import (
	"crypto/sha1"
	"fmt"
	"hash"

	"github.com/mebyus/ku/goku/compiler/enums/tpk"
)

type Type struct {
	nodeSymDef

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
	case tpk.Trivial:
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
	case tpk.AnyPointer:
		return "*void"
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
type TypeHash struct {
	h [20]byte
	f TypeFlag
	k tpk.Kind
}

func (h TypeHash) String() string {
	return fmt.Sprintf("%02X:%04X:%X", h.k, h.f, h.h)
}

func (t *Type) Hash() TypeHash {
	digest := sha1.New()
	t.hash(digest)
	h := TypeHash{
		f: t.Flags,
		k: t.Kind,
	}
	digest.Sum(h.h[:])
	return h
}

func (t *Type) hash(digest hash.Hash) {
	switch def := t.Def.(type) {
	case nil:
		panic("nil def")
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) def (%T)", def.Kind(), def.Kind(), def))
	}
}

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

type Chunk struct {
	// Chunk element type.
	Type *Type
}

// Explicit interface implementation check.
var _ TypeDef = Chunk{}

func (Chunk) Kind() tpk.Kind {
	return tpk.Chunk
}

type Array struct {
	// Array element type.
	Type *Type

	// Number of elements in array.
	Size uint64
}

// Explicit interface implementation check.
var _ TypeDef = Array{}

func (Array) Kind() tpk.Kind {
	return tpk.Array
}

type TypeIndex struct {
	Static StaticTypes

	// Maps span element type to the corresponding span type.
	Spans map[ /* span element type */ *Type]*Type

	// Maps array type definition to the corresponding array type.
	// Arrays map[Array]*Type
}

// StaticTypes contains instances of various predefined (builtin) static types.
type StaticTypes struct {
	// Unsized.
	Integer *Type

	String *Type
}

func (t *StaticTypes) Init() {
	t.Integer = &Type{
		Size:  0, // unsized static integer can hold arbitrary large integer number
		Flags: TypeFlagBuiltin | TypeFlagSigned | TypeFlagStatic,
		Kind:  tpk.Integer,
	}

	t.String = &Type{
		Size:  0,
		Flags: TypeFlagBuiltin | TypeFlagStatic,
		Kind:  tpk.String,
	}
}

func (x *TypeIndex) Init() {
	x.Static.Init()

	x.Spans = make(map[*Type]*Type)
}
