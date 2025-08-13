package tokens

type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Special tokens
	EOF

	Semicolon // ;
	Period    // .
	Colon     // :
	Comma     // ,

	LeftCurly   // {
	RightCurly  // }
	LeftSquare  // [
	RightSquare // ]

	// Keywords

	Fun
	Entry
	Data

	staticLiteralEnd
)

const (
	// Common tokens with baselex package.
	// Order of common tokens must the same as in baselex.

	Illegal Kind = staticLiteralEnd + iota // any byte sequence unknown to lexer

	// Identifiers and basic type literals
	Word       // myvar, main, Line, print
	BinInteger // 0b1101100001
	OctInteger // 0o43671
	DecInteger // 5367, 43432, 1000097
	HexInteger // 0x43da1
	DecFloat   // 123.45
	Rune       // 'a', '\t', 'p'
	String     // "abc", "", "\t\n  42Hello\n"

	// Custom tokens

	Label // @.my_label

	// Comments
	LineComment  // Line comment starts with //
	BlockComment // Comment inside /* comment */ block

	maxKind
)

// FromBaseKind transforms baselex token kind into Kind.
func FromBaseKind(k uint32) Kind {
	return Kind(k) + staticLiteralEnd
}
