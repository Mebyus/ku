package tokens

import (
	"github.com/mebyus/ku/goku/compiler/baselex"
	"github.com/mebyus/ku/goku/compiler/sm"
)

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

	Pin sm.Pin

	Kind Kind

	Flags uint32
}

const (
	FlagKeyword = 1 << iota
)

func (t *Token) IsKeyword() bool {
	return t.Flags&FlagKeyword != 0
}

func FromBaseToken(tok baselex.Token) Token {
	return Token{
		Data:  tok.Data,
		Val:   tok.Val,
		Pin:   tok.Pin,
		Kind:  FromBaseKind(tok.Kind),
		Flags: tok.Flags,
	}
}

func (t *Token) SetIllegalError(code uint64) {
	t.Kind = Illegal
	t.Val = code
}
