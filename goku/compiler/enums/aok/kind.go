package aok

import "github.com/mebyus/ku/goku/compiler/token"

// Kind indicates assignment operation kind.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	Simple // "="
	Walrus // ":="

	Add // "+="
	Sub // "-="
	Mul // "*="
	Div // "/="
	Rem // "%="
)

var text = [...]string{
	empty: "<nil>",

	Simple: "=",
	Walrus: ":=",

	Add: "+=",
	Sub: "-=",
	Mul: "*=",
	Div: "/=",
	Rem: "%=",
}

func (k Kind) String() string {
	return text[k]
}

func FromToken(kind token.Kind) (Kind, bool) {
	var k Kind

	switch kind {
	case token.Assign:
		k = Simple
	case token.Walrus:
		k = Walrus
	case token.AddAssign:
		k = Add
	case token.SubAssign:
		k = Sub
	case token.MulAssign:
		k = Mul
	case token.DivAssign:
		k = Div
	case token.RemAssign:
		k = Rem
	default:
		return 0, false
	}

	return k, true
}
