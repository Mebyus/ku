package bop

import (
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/ku/token"
)

// Op represents binary operator inside expression.
type Op struct {
	Pin  sx.Pin
	Kind Kind
}

// Kind indicates binary operator kind.
type Kind uint32

const (
	Equal    Kind = iota // ==
	NotEqual             // !=

	Less           // <
	Greater        // >
	LessOrEqual    // <=
	GreaterOrEqual // >=

	BoolAnd // &&
	BoolOr  // ||
	In      // in

	Add // +
	Sub // -
	Mul // *
	Div // /
	Mod // %

	Xor        // ^
	BitAnd     // &
	BitOr      // |
	BitAndNot  // &^
	LeftShift  // <<
	RightShift // >>

	// total number of different operators
	num
)

var text = [...]string{
	Equal:    "==",
	NotEqual: "!=",

	Less:           "<",
	Greater:        ">",
	LessOrEqual:    "<=",
	GreaterOrEqual: ">=",

	BoolAnd: "&&",
	BoolOr:  "||",
	In:      "in",

	Add: "+",
	Sub: "-",
	Mul: "*",
	Div: "/",
	Mod: "%",

	Xor:        "^",
	BitAnd:     "&",
	BitOr:      "|",
	BitAndNot:  "&^",
	LeftShift:  "<<",
	RightShift: ">>",
}

func (k Kind) String() string {
	return text[k]
}

var precedence = [...]int{
	Mul:        1,
	Div:        1,
	Mod:        1,
	LeftShift:  1,
	RightShift: 1,
	BitAnd:     1,
	BitAndNot:  1,

	Add:   2,
	Sub:   2,
	Xor:   2,
	BitOr: 2,

	Equal:          3,
	NotEqual:       3,
	Less:           3,
	Greater:        3,
	LessOrEqual:    3,
	GreaterOrEqual: 3,

	In: 3,

	BoolAnd: 4,

	BoolOr: 5,
}

// Precedence gives binary operator precedence.
// Greater values mean later binding
//
//	1 :  *  /  %  <<  >>  &  &^
//	2 :  +  -  ^  |
//	3 :  ==  !=  <  <=  >  >=  in
//	4 :  &&
//	5 :  ||
func (k Kind) Precedence() int {
	return precedence[k]
}

// Power should be used for (recursive) pratt parsing of binary expressions.
// Starting value of power is 0 in such case.
func (k Kind) Power() int {
	return 6 - k.Precedence()
}

func FromToken(t token.Kind) (Kind, bool) {
	var k Kind
	switch t {
	case token.Equal:
		k = Equal
	case token.NotEqual:
		k = NotEqual
	case token.LeftAngle:
		k = Less
	case token.RightAngle:
		k = Greater
	case token.LessOrEqual:
		k = LessOrEqual
	case token.GreaterOrEqual:
		k = GreaterOrEqual
	case token.BoolAnd:
		k = BoolAnd
	case token.BoolOr:
		k = BoolOr
	// case token.In:
	// 	k = In
	case token.Plus:
		k = Add
	case token.Minus:
		k = Sub
	case token.Asterisk:
		k = Mul
	case token.Slash:
		k = Div
	// case token.Percent:
	// 	k = Mod
	// case token.Caret:
	// 	k = Xor
	case token.Ampersand:
		k = BitAnd
	case token.Pipe:
		k = BitOr
	// case token.BitAndNot:
	// 	k = BitAndNot
	// case token.LeftShift:
	// 	k = LeftShift
	// case token.RightShift:
	// 	k = RightShift
	default:
		return 0, false
	}
	return k, true
}
