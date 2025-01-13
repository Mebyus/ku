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
	Pointer
	AnyPointer
	Array
	ArrayPointer
	Chunk
	Enum
	Bag
	Fun
	Tuple
)

var text = [...]string{
	empty: "<nil>",

	Name:         "name",
	FullName:     "name.full",
	Struct:       "struct",
	Pointer:      "pointer",
	Array:        "array",
	ArrayPointer: "pointer.array",
	Chunk:        "chunk",
	Enum:         "enum",
	Bag:          "bag",
	Fun:          "fun",
	Tuple:        "tuple",

	AnyPointer: "pointer.any",
}

func (k Kind) String() string {
	return text[k]
}
