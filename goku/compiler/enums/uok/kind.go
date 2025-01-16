package uok

import "github.com/mebyus/ku/goku/compiler/token"

// Kind indicates unary operator kind.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	Not // !

	Plus  // +
	Minus // -

	BitNot // ^
)

var text = [...]string{
	empty: "<nil>",

	Not:    "!",
	Plus:   "+",
	Minus:  "-",
	BitNot: "^",
}

func (k Kind) String() string {
	return text[k]
}

func FromToken(kind token.Kind) (Kind, bool) {
	var k Kind

	switch kind {
	case token.Not:
		k = Not
	case token.Plus:
		k = Plus
	case token.Minus:
		k = Minus
	case token.Caret:
		k = BitNot
	default:
		return 0, false
	}

	return k, true
}
