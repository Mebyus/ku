package stg

import (
	"fmt"
	"strings"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/enums/typk"
)

type Type struct {
	symdef // can be used as symbol definition of type symbol

	// Local (inside the unit) name of the type. Mainly used for
	// custom and builtin types.
	Name string

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

	// Points to unit where this type was first instantiated.
	//
	// Can be nil for builtin and compiler known types.
	Unit *Unit

	// Byte size of this type's value. May be 0 for some types.
	// More specifically this field equals the stride between two
	// consecutive elements of this type inside an array.
	Size uint32

	// Natural alignment of this type.
	// May be 0 for some types.
	Align uint32

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

	// Target architecture dependant integer.
	TypeArchInt

	// Signed integer type.
	TypeSigned

	// Type has recursive definition.
	TypeRecursive
)

// applicable only for integer types
func (t *Type) IsSigned() bool {
	return t.Flags&TypeSigned != 0
}

func (t *Type) IsStatic() bool {
	return t.Flags&TypeStatic != 0
}

func (t *Type) IsStaticInteger() bool {
	return t.IsStatic() && t.Kind == typk.Integer
}

func (t *Type) String() string {
	if t.Name != "" {
		if t.Unit == nil {
			return t.Name
		}
		return t.Unit.Name + "." + t.Name
	}

	switch t.Kind {
	case typk.Span:
		return "[]" + t.Def.(*Span).Type.String()
	case typk.Ref:
		return "&" + t.Def.(*Ref).Type.String()
	default:
		panic(fmt.Sprintf("unexpected kind (=%d)", t.Kind))
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
	case *ast.Ref:
		return t.lookupRef(s, p)
	case *ast.ArrayPointer:
		return t.lookupArrayPointer(s, p)
	case *ast.Array:
		return t.lookupArray(s, p)
	case *ast.Struct:
		return t.lookupStruct(s, p)
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
		t.report(pin, 0, fmt.Sprintf("unknown type \"%s\"", name))
		return t.common.Types.Invalid
	}
	if symbol.Kind != symk.Type {
		t.report(pin, 0, fmt.Sprintf("expected type here, but name \"%s\" refers to %s symbol", name, symbol.Kind))
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

func (t *Typer) lookupRef(s *Scope, r *ast.Ref) *Type {
	typ := t.LookupType(s, r.Type)
	if typ.IsInvalid() {
		return typ
	}

	return t.common.Types.getRef(typ)
}

func (t *Typer) lookupArrayPointer(s *Scope, p *ast.ArrayPointer) *Type {
	typ := t.LookupType(s, p.Type)
	if typ.IsInvalid() {
		return typ
	}

	return t.common.Types.getArrayPointer(typ)
}

func (t *Typer) lookupArray(s *Scope, a *ast.Array) *Type {
	c := context{
		scope:  s,
		static: true,
	}
	exp := t.convertExp(&c, a.Exp)
	typ := t.LookupType(s, a.Type)
	if typ.IsInvalid() || exp.Type().IsInvalid() {
		return t.common.Types.Invalid
	}

	num, errtext := asArrayNum(exp)
	if errtext != "" {
		t.report(exp.Pin(), 0, errtext)
		return t.common.Types.Invalid
	}
	if num == 0 {
		panic("stub: array of zero size")
	}

	return t.common.Types.getArray(typ, num)
}

func (t *Typer) lookupStruct(s *Scope, p *ast.Struct) *Type {
	if len(p.Fields) == 0 {
		// TODO: return void type
		panic("TODO")
	}

	key, fields := t.lookupFields(s, p.Fields)

	list := t.common.Types.Structs[key]
	for _, typ := range list {
		if equalFieldTypes(fields, typ.Def.(*Struct).Fields) {
			return typ
		}
	}

	def := &Struct{Fields: fields}
	def.layout()
	typ := &Type{
		Def:   def,
		Size:  def.Size,
		Align: def.Align,
		Unit:  t.unit,
		Kind:  typk.Struct,
	}
	t.common.Types.Structs[key] = append(t.common.Types.Structs[key], typ)
	return typ
}

func (t *Typer) lookupFields(s *Scope, list []ast.Field) (string, []Field) {
	n := 0 // total length of joined field names
	fields := make([]Field, 0, len(list))
	for _, f := range list {
		t := t.LookupType(s, f.Type)

		name := f.Name
		fields = append(fields, Field{
			Name: name,
			Type: t,
		})
		n += len(name)
	}
	n += len(fields) - 1 // for field separators

	// builds joined field names in the form:
	//	"field_1|field_2|field_3"
	var buf strings.Builder
	buf.Grow(n)

	buf.WriteString(fields[0].Name)
	for _, f := range fields[1:] {
		buf.WriteByte('|')
		buf.WriteString(f.Name)
	}
	key := buf.String()

	return key, fields
}

// Try to convert static expression into array
// element number at compile time
//
// Returns (num, "") when operation was successful.
// Return (0, "<...>") with error description otherwise.
func asArrayNum(exp Exp) (uint, string) {
	typ := exp.Type()
	if !typ.IsStaticInteger() {
		return 0, fmt.Sprintf("cannot use %s expression as number of array elements", typ)
	}

	n := exp.(*Integer)
	if n.Val == 0 {
		return 0, ""
	}

	if n.Neg {
		return 0, fmt.Sprintf("negative number (-%d) of array elements", n.Val)
	}

	return uint(n.Val), ""
}

type TypeIndex struct {
	Static StaticTypes
	Known  KnownTypes

	// Maps type referred by reference to corresponding
	// reference type.
	Refs map[ /* type referred by reference */ *Type]*Type

	// Maps type referred by array pointer to corresponding array pointer type.
	ArrayPointers map[ /* type referred by pointer */ *Type]*Type

	// Maps span element type to the corresponding span type.
	Spans map[ /* span element type */ *Type]*Type

	// Maps array type definition (element type + size)
	// to the corresponding array type.
	Arrays map[Array]*Type

	// Maps joined field names to a list of struct types with the same field
	// names (including field order).
	Structs map[ /* joined field names */ string][]*Type

	Invalid *Type

	ArchPointerSize uint32
}

func (x *TypeIndex) init(archPointerSize uint32) {
	x.ArchPointerSize = archPointerSize

	x.Static.init()
	x.Known.init(archPointerSize)

	x.Invalid = &Type{Kind: typk.Invalid}

	x.Refs = make(map[*Type]*Type)
	x.ArrayPointers = make(map[*Type]*Type)
	x.Spans = make(map[*Type]*Type)
	x.Arrays = make(map[Array]*Type)
	x.Structs = make(map[string][]*Type)
}

// StaticTypes contains instances of various predefined (builtin) static types.
type StaticTypes struct {
	// For nil literal.
	// Nil *Type

	// Unsized and signed.
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
		Name:  "<int>",
		Size:  0, // unsized static integer can hold arbitrary large integer number
		Flags: TypeBuiltin | TypeSigned | TypeStatic,
		Kind:  typk.Integer,
	}

	t.Boolean = &Type{
		Name:  "<bool>",
		Flags: TypeBuiltin | TypeStatic,
		Kind:  typk.Boolean,
	}

	t.String = &Type{
		Name:  "<str>",
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

	Int  *Type
	Uint *Type
	Sint *Type

	U8  *Type
	U32 *Type

	S32 *Type

	Bool *Type

	Str *Type
}

func (t *KnownTypes) init(archPointerSize uint32) {
	t.Int = &Type{
		Size:  archPointerSize,
		Align: archPointerSize,
		Flags: TypeBuiltin | TypeArchInt | TypeSigned,
		Kind:  typk.Integer,
	}
	t.Uint = &Type{
		Size:  archPointerSize,
		Align: archPointerSize,
		Flags: TypeBuiltin | TypeArchInt,
		Kind:  typk.Integer,
	}
	t.Sint = &Type{
		Size:  archPointerSize,
		Align: archPointerSize,
		Flags: TypeBuiltin | TypeArchInt | TypeSigned,
		Kind:  typk.Integer,
	}

	t.U8 = &Type{
		Size:  1,
		Align: 1,
		Flags: TypeBuiltin,
		Kind:  typk.Integer,
	}

	t.U32 = &Type{
		Size:  4,
		Align: 4,
		Flags: TypeBuiltin,
		Kind:  typk.Integer,
	}

	t.S32 = &Type{
		Size:  4,
		Align: 4,
		Flags: TypeBuiltin | TypeSigned,
		Kind:  typk.Integer,
	}

	t.Bool = &Type{
		Size:  1,
		Align: 1,
		Flags: TypeBuiltin,
		Kind:  typk.Boolean,
	}

	t.Str = &Type{
		Size:  2 * archPointerSize,
		Align: archPointerSize,
		Flags: TypeBuiltin,
		Kind:  typk.String,
	}
}

func (x *TypeIndex) getRef(t *Type) *Type {
	typ, ok := x.Refs[t]
	if ok {
		return typ
	}
	typ = &Type{
		Def:   &Ref{Type: t},
		Size:  x.ArchPointerSize,
		Align: x.ArchPointerSize,
		Kind:  typk.Ref,
	}
	x.Refs[t] = typ
	return typ
}

func (x *TypeIndex) getArrayPointer(t *Type) *Type {
	typ, ok := x.ArrayPointers[t]
	if ok {
		return typ
	}
	typ = &Type{
		Def:   &Ref{Type: t},
		Size:  x.ArchPointerSize,
		Align: x.ArchPointerSize,
		Kind:  typk.ArrayPointer,
	}
	x.Refs[t] = typ
	return typ
}

func (x *TypeIndex) getSpan(t *Type) *Type {
	typ, ok := x.Spans[t]
	if ok {
		return typ
	}
	typ = &Type{
		Def:   &Span{Type: t},
		Size:  2 * x.ArchPointerSize,
		Align: x.ArchPointerSize,
		Kind:  typk.Span,
	}
	x.Spans[t] = typ
	return typ
}

func (x *TypeIndex) getArray(elem *Type, num uint) *Type {
	a := Array{Type: elem, Num: num}
	typ, ok := x.Arrays[a]
	if ok {
		return typ
	}
	typ = &Type{
		Def:   &a,
		Size:  uint32(num) * elem.Size,
		Align: elem.Align,
		Kind:  typk.Array,
	}
	x.Arrays[a] = typ
	return typ
}
