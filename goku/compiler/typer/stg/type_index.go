package stg

import (
	"encoding/binary"
	"fmt"
	"strings"
	"unsafe"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
)

type TypeIndex struct {
	Static StaticTypes

	Known KnownTypes

	// Maps span element type to the corresponding span type.
	Spans map[ /* span element type */ *Type]*Type

	// Maps type referred by pointer to corresponding pointer type.
	Pointers map[ /* type referred by pointer */ *Type]*Type

	// Maps type referred by array pointer to corresponding array pointer type.
	ArrayPointers map[ /* type referred by pointer */ *Type]*Type

	// Maps type referred by reference to corresponding reference type.
	Refs map[ /* type referred by reference */ *Type]*Type

	// Maps joined field names to a list of struct types with the same field
	// names (including field order).
	Structs map[ /* joined field names */ string][]*Type

	Tuples map[ /* binary encoded string of all types inside tuple */ string]*Type

	// Maps array type definition (element type + size) to the corresponding array type.
	Arrays map[Array]*Type
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

// KnownTypes contains instances of various primitive runtime types and their derivatives
// (spans, pointers, etc.).
//
// Mostly used as shorthand access to types which compiler should be aware of.
type KnownTypes struct {
	// void, empty struct, zero size array
	Void *Type

	// *void
	VoidPointer *Type
}

func (t *KnownTypes) Init() {
	t.Void = &Type{
		Size:  0,
		Flags: TypeFlagBuiltin,
		Kind:  tpk.Trivial,
	}
	t.VoidPointer = &Type{
		Size:  archPointerSize,
		Flags: TypeFlagBuiltin,
		Kind:  tpk.AnyPointer,
	}
}

func (x *TypeIndex) Init() {
	x.Static.Init()
	x.Known.Init()

	x.Spans = make(map[*Type]*Type)
	x.Pointers = make(map[*Type]*Type)
	x.ArrayPointers = make(map[*Type]*Type)
	x.Refs = make(map[*Type]*Type)
	x.Structs = make(map[string][]*Type)
	x.Tuples = make(map[string]*Type)
	x.Arrays = make(map[Array]*Type)
}

func (x *TypeIndex) Lookup(scope *Scope, spec ast.TypeSpec) (*Type, diag.Error) {
	typ, err := x.lookup(scope, spec)
	if err != nil {
		return nil, err
	}
	if typ == nil {
		panic(fmt.Sprintf("%s (=%d) type specifier (%T) produced no type", spec.Kind(), spec.Kind(), spec))
	}
	return typ, nil
}

func (x *TypeIndex) lookup(scope *Scope, spec ast.TypeSpec) (*Type, diag.Error) {
	switch p := spec.(type) {
	case ast.TypeName:
		return x.lookupTypeName(scope, p)
	case ast.VoidPointer:
		return x.Known.VoidPointer, nil
	case ast.Pointer:
		return x.lookupPointer(scope, p)
	case ast.Ref:
		return x.lookupRef(scope, p)
	case ast.Void:
		return x.Known.Void, nil
	case ast.TypeFullName:
	case ast.ArrayPointer:
		return x.lookupArrayPointer(scope, p)
	case ast.Span:
		return x.lookupSpan(scope, p)
	case ast.Array:
		return x.lookupArray(scope, p)
	case ast.Struct:
		return x.lookupStruct(scope, p)
	case ast.Tuple:
		return x.lookupTuple(scope, p)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}
	return nil, nil
}

func (x *TypeIndex) lookupTypeName(scope *Scope, p ast.TypeName) (*Type, diag.Error) {
	name := p.Name.Str
	pin := p.Name.Pin

	s := scope.Lookup(name)
	if s == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("type name \"%s\" refers to undefined symbol", name),
		}
	}
	if s.Kind != smk.Type {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("name \"%s\" refers to %s, not a type", name, s.Kind),
		}
	}

	return s.Def.(*Type), nil
}

func (x *TypeIndex) lookupArray(scope *Scope, a ast.Array) (*Type, diag.Error) {
	sizeExp, err := scope.EvalConstExp(a.Size)
	if err != nil {
		return nil, err
	}
	integer, err := expectInteger(sizeExp)
	if err != nil {
		return nil, err
	}
	if integer.Neg {
		return nil, &diag.SimpleMessageError{
			Pin:  integer.Pin,
			Text: "negative number of elements in array declaration",
		}
	}
	size := integer.Val
	if size == 0 {
		return x.Known.Void, nil
	}

	t, err := x.lookup(scope, a.Type)
	if err != nil {
		return nil, err
	}

	def := Array{
		Type: t,
		Len:  uint32(size),
	}
	typ, ok := x.Arrays[def]
	if ok {
		return typ, nil
	}
	typ = &Type{
		Def:  def,
		Size: t.Size * uint32(size),
		Kind: tpk.Array,
	}
	x.Arrays[def] = typ

	return typ, nil
}

func (x *TypeIndex) lookupPointer(scope *Scope, p ast.Pointer) (*Type, diag.Error) {
	t, err := x.lookup(scope, p.Type)
	if err != nil {
		return nil, err
	}

	typ, ok := x.Pointers[t]
	if ok {
		return typ, nil
	}
	typ = &Type{
		Def:  Pointer{Type: t},
		Size: archPointerSize,
		Kind: tpk.Pointer,
	}
	x.Pointers[t] = typ

	return typ, nil
}

func (x *TypeIndex) lookupArrayPointer(scope *Scope, p ast.ArrayPointer) (*Type, diag.Error) {
	t, err := x.lookup(scope, p.Type)
	if err != nil {
		return nil, err
	}

	typ, ok := x.ArrayPointers[t]
	if ok {
		return typ, nil
	}
	typ = &Type{
		Def:  ArrayPointer{Type: t},
		Size: archPointerSize,
		Kind: tpk.ArrayPointer,
	}
	x.ArrayPointers[t] = typ

	return typ, nil
}

func (x *TypeIndex) lookupRef(scope *Scope, p ast.Ref) (*Type, diag.Error) {
	t, err := x.lookup(scope, p.Type)
	if err != nil {
		return nil, err
	}

	typ, ok := x.Refs[t]
	if ok {
		return typ, nil
	}
	typ = &Type{
		Def:  Ref{Type: t},
		Size: archPointerSize,
		Kind: tpk.Ref,
	}
	x.Refs[t] = typ

	return typ, nil
}

func (x *TypeIndex) lookupSpan(scope *Scope, p ast.Span) (*Type, diag.Error) {
	t, err := x.lookup(scope, p.Type)
	if err != nil {
		return nil, err
	}

	typ, ok := x.Spans[t]
	if ok {
		return typ, nil
	}
	typ = &Type{
		Def:  Chunk{Type: t},
		Size: 2 * archPointerSize,
		Kind: tpk.Chunk,
	}
	x.Spans[t] = typ

	return typ, nil
}

func (x *TypeIndex) lookupStruct(scope *Scope, p ast.Struct) (*Type, diag.Error) {
	if len(p.Fields) == 0 {
		panic("struct with no fields")
	}

	n := 0 // total length of joined field names
	fields := make([]Field, 0, len(p.Fields))
	for _, f := range p.Fields {
		t, err := x.lookup(scope, f.Type)
		if err != nil {
			return nil, err
		}

		name := f.Name.Str
		fields = append(fields, Field{
			Name: name,
			Type: t,
			// TODO: calculate offset
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
	joined := buf.String()

	list := x.Structs[joined]
	for _, t := range list {
		if equalFieldTypes(fields, t.Def.(Struct).Fields) {
			return t, nil
		}
	}
	typ := &Type{
		// TODO: calculate size
		Def:  Struct{Fields: fields},
		Kind: tpk.Struct,
	}
	x.Structs[joined] = append(x.Structs[joined], typ)
	return typ, nil
}

func (x *TypeIndex) lookupTuple(scope *Scope, tuple ast.Tuple) (*Type, diag.Error) {
	if len(tuple.Types) == 0 {
		panic("empty tuple")
	}

	types := make([]*Type, 0, len(tuple.Types))
	for _, p := range tuple.Types {
		typ, err := x.lookup(scope, p)
		if err != nil {
			return nil, err
		}
		types = append(types, typ)
	}

	key := encodeTypesAsKey(types)
	typ, ok := x.Tuples[key]
	if ok {
		return typ, nil
	}
	typ = &Type{
		// TODO: calculate size
		Def:  Tuple{Types: types},
		Kind: tpk.Tuple,
	}
	x.Tuples[key] = typ

	return typ, nil
}

// Checks that all corresponding fields in two lists have the same types:
//
//	a[0].Type == b[0].Type
//	a[1].Type == b[1].Type
//	...
//
// Both slices are assumed to have the same length.
func equalFieldTypes(a, b []Field) bool {
	for i := range len(a) {
		if a[i].Type != b[i].Type {
			return false
		}
	}
	return true
}

// HACK: this function produces a non-utf8 string which is essentially a slice of
// binary encoded pointer values of each *Type in list.
//
// We use this hack to uniquely identify tuple types with map[string]*Type.
func encodeTypesAsKey(list []*Type) string {
	var buf strings.Builder
	buf.Grow(8 * len(list))

	for _, t := range list {
		p := uint64(uintptr(unsafe.Pointer(t)))

		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], p)
		buf.Write(b[:])
	}

	return buf.String()
}
