package token

import "github.com/mebyus/ku/internal/ku/sx"

type Kind uint32

type Token struct {
	// Not empty only for tokens which do not have static literal.
	//
	// Examples: identifiers, strings, illegal tokens
	//
	// For tokens obtained from regular string literals (as in "some string")
	// this field contains unescaped string value.
	Data string

	// Meaning of this value is dependant on token Kind
	//
	//	Integer:	parsed integer value (if it fits into 64 bits)
	//	Rune:		integer value of code point
	//	EOF:		error code (can be 0, in case end of text was reached without error)
	//	Illegal:	error code (always not 0)
	Val uint64

	Pin sx.Pin

	Kind Kind

	Flags uint32
}

const (
	// Invalid token.
	INV Kind = iota

	// End of source text.
	EOF

	Semicolon // ;

	Colon // :
	Comma // ,

	RightArrow // ->

	Plus     // +
	Minus    // -
	Asterisk // *
	Slash    // /

	LeftCurly  // {
	RightCurly // }
	LeftParen  // (
	RightParen // )

	// Keywords.
	Fun
	Return

	True
	False

	// Identifiers and basic type literals.
	Word
	Integer
)

const (
	DecInt = iota
	BinInt
	OctInt
	HexInt
)
