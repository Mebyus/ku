package tsk

// Kind indicates type specifier kind.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	Name
	FullName
	Struct
	Union
	Ref
	Pointer
	AnyRef
	AnyPointer
	Array
	ArrayRef
	ArrayPointer
	Chunk
	Enum
	Bag
	Fun
	Tuple
	Form
	Trivial
	Type
)

var text = [...]string{
	empty: "<nil>",

	Name:         "name",
	FullName:     "name.full",
	Struct:       "struct",
	Union:        "union",
	Ref:          "ref",
	Pointer:      "pointer",
	Array:        "array",
	ArrayRef:     "ref.array",
	ArrayPointer: "pointer.array",
	Chunk:        "chunk",
	Enum:         "enum",
	Bag:          "bag",
	Fun:          "fun",
	Tuple:        "tuple",
	Form:         "form",
	Trivial:      "trivial",
	Type:         "type",

	AnyPointer: "pointer.any",
	AnyRef:     "ref.any",
}

func (k Kind) String() string {
	return text[k]
}
