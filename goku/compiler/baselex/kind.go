package baselex

// Token kinds common for several lexers inside the project.
const (
	Illegal = iota // any byte sequence unknown to lexer

	// Identifiers and basic type literals
	Word       // myvar, main, Line, print
	BinInteger // 0b1101100001
	OctInteger // 0o43671
	DecInteger // 5367, 43432, 1000097
	HexInteger // 0x43da1
	DecFloat   // 123.45
	Rune       // 'a', '\t', 'p'
	String     // "abc", "", "\t\n  42Hello\n"
)
