package token

import "github.com/mebyus/ku/internal/ku/sx"

// Kind the most significant token property which determines
// its lexical meaning and serves as discriminator for all
// other token properties.
type Kind uint8

// Token represents a minimal unit of language lexicon.
// Meaning of each token is determined by its kind.
//
// Ultimately when it comes to malformed or invalid tokens
// this struct can represent 3 distinct token states:
//
//   - normal  (well formed token, with no "flaws")
//   - flawed  (recognized by lexer, but contains "flaws")
//   - invalid (piece of text unrecognizable by lexer)
//
// It worth noting that a program cannot be fully compiled
// if it has at least one flawed or invalid token. This is
// also true if program does not have any other parsing or
// typechecking errors.
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
	//	INV:		error code (always not 0)
	Val uint64

	Pin sx.Pin

	// Auxiliary value. Often used for specifying flaw or error type.
	Aux uint32

	// Used for some types of tokens:
	//
	//	Integer: specifies literal style (decimal, binary, etc.)
	Style Style

	// Carry additional token properties which can vary
	// for the same token type.
	//
	// Bitflags field.
	Flags Flag

	Kind Kind
}

type Flag uint8

const (
	// Token has a flaw if flag is set.
	FlagFlaw Flag = 1 << iota

	// Token is a keyword if flag is set.
	FlagKey

	// Token is a compiler directive if flag is set.
	FlagDir
)

type Style uint8

const (
	Dec Style = iota
	Bin
	Oct
	Hex
)

// Flaw returns true if token has a flaw.
func (t *Token) Flaw() bool {
	return t.Flags&FlagFlaw != 0
}

const (
	// Invalid token.
	INV Kind = iota

	// End of source text.
	EOF

	Semicolon // ;

	Colon  // :
	Period // .
	Comma  // ,
	AtSign // @

	RightArrow // ->

	Plus      // +
	Minus     // -
	Asterisk  // *
	Slash     // /
	Exclam    // !
	Ampersand // &
	Pipe      // |

	Under  // _
	Assign // =
	Walrus // :=

	Equal          // ==
	NotEqual       // !=
	LessOrEqual    // <=
	GreaterOrEqual // >=

	BoolOr  // ||
	BoolAnd // &&

	LeftAngle   // <
	RightAngle  // >
	LeftCurly   // {
	RightCurly  // }
	LeftParen   // (
	RightParen  // )
	LeftSquare  // [
	RightSquare // ]

	Address   // .&
	Deref     // .*
	Alter     // .{
	AlterZero // .{}

	SliceAll // [:]

	HashSquare // #[

	// start of keyword Kinds
	// used by some formatting functions
	_keystart

	Pub
	Fun
	Var
	Type
	Zone
	Const
	Import

	Struct
	Union

	If
	Else
	Return
	Goto
	For

	True
	False
	Nil

	_keyend

	// Identifiers and basic type literals.
	Word
	Integer
	String
	Rune
)

// Category returns token kind category in string format.
// Category indicates broad role of token in the language.
//
// It is mostly used for token listings and parser error reports.
func (k Kind) Category() string {
	switch k {
	case INV:
		return "INV"
	case EOF:
		return "EOF"
	case Word:
		return "WORD"
	case Integer:
		return "INTEGER"
	case String:
		return "STRING"
	case Rune:
		return "RUNE"
	}

	if k > _keystart {
		return "KEYWORD"
	}
	return "OTHER"
}
