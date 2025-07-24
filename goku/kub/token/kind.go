package token

type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Special tokens
	EOF

	// Operators and/or punctuators

	Assign    // =
	Semicolon // ;
	Period    // .

	Equal          // ==
	NotEqual       // !=
	LessOrEqual    // <=
	GreaterOrEqual // >=
	LeftAngle      // <
	RightAngle     // >
	Not            // !

	// Brackets

	LeftParen  // (
	RightParen // )
	LeftCurly  // {
	RightCurly // }

	// Keywords

	Import
	Test
	Exe
	Include
	Set
	Unit
	Main
	Module
	Link
	If
	Else

	staticLiteralEnd

	Illegal // any byte sequence unknown to lexer

	// Identifiers and basic type literals
	Word       // myvar, main, Line, print
	DecInteger // 5367, 43432, 1000097
	HexInteger // 0x43da1
	String     // "abc", "", "\t\n  42Hello\n"

	// Comments
	LineComment  // Line comment starts with //
	BlockComment // Comment inside /* comment */ block

	maxKind
)
