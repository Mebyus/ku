package token

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mebyus/ku/goku/compiler/char"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

// literal maps Kind to token static literal or code string.
var literal = [...]string{
	empty: "<nil>",

	EOF: "EOF",

	// Operators/punctuators

	Equal:          "==",
	NotEqual:       "!=",
	LessOrEqual:    "<=",
	GreaterOrEqual: ">=",
	Not:            "!",
	Semicolon:      ";",
	Period:         ".",
	Assign:         "=",
	LeftAngle:      "<",
	RightAngle:     ">",
	LeftParen:      "(",
	RightParen:     ")",
	LeftCurly:      "{",
	RightCurly:     "}",

	// Keywords

	Import:  "import",
	Include: "include",
	Test:    "test",
	Set:     "set",
	Module:  "module",
	Unit:    "unit",
	Main:    "main",
	Link:    "link",
	If:      "if",
	Else:    "else",

	// Non static literals

	Illegal:    "ILG",
	Word:       "WORD",
	String:     "STR",
	DecInteger: "INT.DEC",
	HexInteger: "INT.HEX",

	// Comments

	LineComment:  "COM.LINE",
	BlockComment: "COM.BLOCK",
}

func (k Kind) String() string {
	return literal[k]
}

func (k Kind) hasStaticLiteral() bool {
	return k < staticLiteralEnd
}

func (t Token) String() string {
	if t.Kind.hasStaticLiteral() {
		return t.Kind.String()
	}

	switch t.Kind {
	case Word:
		return t.Data
	case Illegal:
		if t.Val == LengthOverflow {
			return "Illegal(token length overflow)"
		}
		return "[[" + t.Data + "]]"
	case DecInteger:
		return strconv.FormatUint(t.Val, 10)
	case HexInteger:
		return "0x" + strconv.FormatUint(t.Val, 16)
	case String:
		return "\"" + char.Escape(t.Data) + "\""
	default:
		panic(fmt.Sprintf("unexpected \"%s\" token kind (=%d)", t.Kind, t.Kind))
	}
}

func FormatTokenLine(m srcmap.PinMap, t Token) string {
	pos, err := m.DecodePin(t.Pin)
	if err != nil {
		panic(err)
	}

	buf := strings.Builder{}
	buf.Grow(64)

	buf.WriteString(pos.Pos.String())
	for range 16 - buf.Len() {
		buf.WriteByte(' ')
	}

	if t.Kind == EOF {
		buf.WriteString(t.Kind.String())
		return buf.String()
	}

	if !t.Kind.hasStaticLiteral() {
		buf.WriteString(t.Kind.String())
	}
	for range 32 - buf.Len() {
		buf.WriteByte(' ')
	}

	buf.WriteString(t.String())

	return buf.String()
}
