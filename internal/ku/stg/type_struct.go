package stg

import "github.com/mebyus/ku/internal/ku/sx"

// Field represents struct or union field.
type Field struct {
	Name string

	// name pin in struct definition
	Pin sx.Pin

	Type *Type

	Offset uint32

	// Number of bytes field uses inside the struct,
	// including padding.
	Size uint32

	Align uint32

	// size of padding inserted before this field in bytes
	Pad uint32
}

// Struct represents struct type definition.
type Struct struct {
	typdef

	Fields []Field

	Size  uint32
	Align uint32

	// Size of tail padding.
	Tail uint32

	// Total number of bytes "wasted" in paddings.
	Waste uint32

	// Total number of non-zero paddings used.
	// Includes tail padding if present.
	Npads uint32
}

func (s *Struct) getField(name string) *Field {
	for i := range s.Fields {
		f := &s.Fields[i]
		if f.Name == name {
			return f
		}
	}
	return nil
}

// fills layout of struct fields
func (s *Struct) layout() {
	var offset uint32
	var align uint32 // struct alignment
	var n uint32     // number of used paddings
	var w uint32     // struct waste

	for i := range len(s.Fields) {
		f := &s.Fields[i]

		a := f.Type.Align
		s := f.Type.Size

		var pad uint32
		rem := offset % a
		if rem != 0 {
			pad = a - rem
			offset += pad
			w += pad
			n += 1
		}

		f.Offset = offset
		f.Size = s + pad
		f.Align = a
		f.Pad = pad

		offset += s
		if a > align {
			align = a
		}
	}

	// add tail padding if necessary
	if offset != 0 && align != 0 {
		rem := offset % align
		if rem != 0 {
			pad := align - rem
			offset += pad
			s.Tail = pad
			w += pad
			n += 1
		}
	}

	s.Align = align
	s.Size = offset
	s.Waste = w
	s.Npads = n
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
