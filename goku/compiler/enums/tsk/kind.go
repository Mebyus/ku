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
	Map
	Ref
	Pointer
	VoidRef
	VoidPointer
	Array
	ArrayRef
	ArrayPointer
	Span
	CapBuf
	Enum
	Bag
	Fun
	Tuple
	Form
	Void
	Type
)

var text = [...]string{
	empty: "<nil>",

	Name:         "name",
	FullName:     "full name",
	Struct:       "struct",
	Union:        "union",
	Map:          "map",
	Ref:          "ref",
	Pointer:      "pointer",
	Array:        "array",
	ArrayRef:     "array ref",
	ArrayPointer: "array pointer",
	Span:         "span",
	CapBuf:       "capbuf",
	Enum:         "enum",
	Bag:          "bag",
	Fun:          "fun",
	Tuple:        "tuple",
	Form:         "form",
	Void:         "void",
	Type:         "type",

	VoidPointer: "void pointer",
	VoidRef:     "void ref",
}

func (k Kind) String() string {
	return text[k]
}
