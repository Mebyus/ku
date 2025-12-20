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
	// For nil literal.
	Nil *Type

	// Unsized.
	Integer *Type

	String *Type
}

func (t *StaticTypes) Init() {
	t.Nil = &Type{
		Size:  0,
		Flags: TypeFlagBuiltin | TypeFlagStatic,
		Kind:  tpk.Nil,
	}

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

	// &void
	VoidRef *Type

	// bool
	Bool *Type
}

func (t *KnownTypes) Init() {
	t.Void = &Type{
		Size:  0,
		Flags: TypeFlagBuiltin,
		Kind:  tpk.Void,
	}
	t.VoidPointer = &Type{
		Size:  archPointerSize,
		Flags: TypeFlagBuiltin,
		Kind:  tpk.VoidPointer,
	}
	t.VoidRef = &Type{
		Size:  archPointerSize,
		Flags: TypeFlagBuiltin,
		Kind:  tpk.VoidRef,
	}
	t.Bool = &Type{
		Size:  1,
		Flags: TypeFlagBuiltin,
		Kind:  tpk.Boolean,
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

func (s *Scope) LookupType(spec ast.TypeSpec) (*Type, diag.Error) {
	switch p := spec.(type) {
	case ast.VoidPointer:
		return s.Types.Known.VoidPointer, nil
	case ast.VoidRef:
		return s.Types.Known.VoidRef, nil
	case ast.Void:
		return s.Types.Known.Void, nil
	case ast.TypeName:
		return s.lookupTypeName(p)
	case ast.Pointer:
		return s.lookupPointer(p)
	case ast.Ref:
		return s.lookupRef(p)
	case ast.TypeFullName:
		return s.lookupTypeFullName(p)
	case ast.ArrayPointer:
		return s.lookupArrayPointer(p)
	case ast.Span:
		return s.lookupSpan(p)
	case ast.Array:
		return s.lookupArray(p)
	case ast.Struct:
		return s.lookupStruct(p)
	case ast.Tuple:
		return s.lookupTuple(p)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", p.Kind(), p.Kind(), p))
	}
}

func (s *Scope) lookupTypeName(p ast.TypeName) (*Type, diag.Error) {
	name := p.Name.Str
	pin := p.Name.Pin

	symbol := s.Lookup(name)
	if symbol == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("type name \"%s\" refers to undefined symbol", name),
		}
	}
	if symbol.Kind != smk.Type {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("name \"%s\" refers to %s, not a type", name, s.Kind),
		}
	}

	return symbol.Def.(SymDefType).Type, nil
}

func (s *Scope) lookupTypeFullName(p ast.TypeFullName) (*Type, diag.Error) {
	iname := p.Import.Str
	m := s.Lookup(iname)
	if m == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", iname),
		}
	}
	if m.Kind != smk.Import {
		return nil, &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to %s, not an import", iname, m.Kind),
		}
	}

	unit := m.Def.(SymDefUnit).Unit
	name := p.Name.Str
	symbol := unit.Scope.Lookup(name)
	if symbol == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", name), // TODO: error text with unit name
		}
	}
	if symbol.Kind != smk.Type {
		return nil, &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("name \"%s\" refers to %s, not a type", name, s.Kind),
		}
	}
	if !symbol.IsPublic() {
		return nil, &diag.SimpleMessageError{
			Pin:  p.Name.Pin,
			Text: fmt.Sprintf("type \"%s\" is not public", name),
		}
	}

	return symbol.Def.(SymDefType).Type, nil
}

func (s *Scope) lookupArray(a ast.Array) (*Type, diag.Error) {
	sizeExp, err := s.EvalConstExp(a.Size)
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
		return s.Types.Known.Void, nil
	}

	t, err := s.LookupType(a.Type)
	if err != nil {
		return nil, err
	}

	def := Array{
		Type: t,
		Len:  uint32(size),
	}
	typ, ok := s.Types.Arrays[def]
	if ok {
		return typ, nil
	}
	typ = &Type{
		Def:  def,
		Size: t.Size * uint32(size),
		Kind: tpk.Array,
	}
	s.Types.Arrays[def] = typ

	return typ, nil
}

func (s *Scope) lookupPointer(p ast.Pointer) (*Type, diag.Error) {
	t, err := s.LookupType(p.Type)
	if err != nil {
		return nil, err
	}

	typ, ok := s.Types.Pointers[t]
	if ok {
		return typ, nil
	}
	typ = &Type{
		Def:  Pointer{Type: t},
		Size: archPointerSize,
		Kind: tpk.Pointer,
	}
	s.Types.Pointers[t] = typ

	return typ, nil
}

func (s *Scope) lookupArrayPointer(p ast.ArrayPointer) (*Type, diag.Error) {
	t, err := s.LookupType(p.Type)
	if err != nil {
		return nil, err
	}

	typ, ok := s.Types.ArrayPointers[t]
	if ok {
		return typ, nil
	}
	typ = &Type{
		Def:  ArrayPointer{Type: t},
		Size: archPointerSize,
		Kind: tpk.ArrayPointer,
	}
	s.Types.ArrayPointers[t] = typ

	return typ, nil
}

func (s *Scope) lookupRef(p ast.Ref) (*Type, diag.Error) {
	t, err := s.LookupType(p.Type)
	if err != nil {
		return nil, err
	}

	typ, ok := s.Types.Refs[t]
	if ok {
		return typ, nil
	}
	typ = &Type{
		Def:  Ref{Type: t},
		Size: archPointerSize,
		Kind: tpk.Ref,
	}
	s.Types.Refs[t] = typ

	return typ, nil
}

func (s *Scope) lookupSpan(p ast.Span) (*Type, diag.Error) {
	t, err := s.LookupType(p.Type)
	if err != nil {
		return nil, err
	}

	typ, ok := s.Types.Spans[t]
	if ok {
		return typ, nil
	}
	typ = &Type{
		Def:  Span{Type: t},
		Size: 2 * archPointerSize,
		Kind: tpk.Span,
	}
	s.Types.Spans[t] = typ

	return typ, nil
}

func (s *Scope) lookupStruct(p ast.Struct) (*Type, diag.Error) {
	if len(p.Fields) == 0 {
		panic("struct with no fields")
	}

	n := 0 // total length of joined field names
	fields := make([]Field, 0, len(p.Fields))
	for _, f := range p.Fields {
		t, err := s.LookupType(f.Type)
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

	list := s.Types.Structs[joined]
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
	s.Types.Structs[joined] = append(s.Types.Structs[joined], typ)
	return typ, nil
}

func (s *Scope) lookupTuple(tuple ast.Tuple) (*Type, diag.Error) {
	if len(tuple.Types) == 0 {
		return s.Types.Known.Void, nil
	}
	if len(tuple.Types) == 1 {
		typ, err := s.LookupType(tuple.Types[0])
		if err != nil {
			return nil, err
		}
		return typ, nil
	}

	types := make([]*Type, 0, len(tuple.Types))
	for _, p := range tuple.Types {
		typ, err := s.LookupType(p)
		if err != nil {
			return nil, err
		}
		types = append(types, typ)
	}

	return s.Types.getTuple(types), nil
}

func (x *TypeIndex) getTuple(types []*Type) *Type {
	key := encodeTypesAsKey(types)
	typ, ok := x.Tuples[key]
	if ok {
		return typ
	}
	typ = &Type{
		// TODO: calculate size
		Def:  Tuple{Types: types},
		Kind: tpk.Tuple,
	}
	x.Tuples[key] = typ
	return typ
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
