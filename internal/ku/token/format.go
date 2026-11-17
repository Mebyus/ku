package token

import (
	"fmt"
	"strconv"

	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/char"
)

// Format token as a single text line (without source position).
func (t *Token) Format(g *fbuf.Buffer) {
	cat := t.Kind.Category()
	g.Puts(cat)
	g.RepSpace(12 - len(cat))

	var flaw string
	if t.Kind == INV || t.Flaw() {
		flaw = strconv.FormatUint(uint64(t.Aux), 10)
	}
	g.Puts(flaw)
	g.RepSpace(12 - len(flaw))

	var text string
	switch t.Kind {
	case EOF, INV:
		// do nothing
	case Word:
		text = t.Data
	case Rune:
		panic("TODO")
	case Integer:
		text = t.IntegerText()
	case String:
		// TODO: refactor into just fbuf usage with Len() storing and Truncate() rollback
		text = "\"" + char.Escape(t.Data) + "\""
	default:
		text = ktext[t.Kind]
		if text == "" {
			panic(fmt.Sprintf("empty token kind (=%d) text", t.Kind))
		}
	}
	g.Puts(text)
}

// IntegerText must only be used if token is an integer.
func (t *Token) IntegerText() string {
	if t.Data != "" {
		return t.Data
	}

	switch t.Style {
	case Dec:
		return strconv.FormatUint(t.Val, 10)
	case Bin:
		return "0b" + strconv.FormatUint(t.Val, 2)
	case Oct:
		return "0o" + strconv.FormatUint(t.Val, 8)
	case Hex:
		return "0x" + strconv.FormatUint(t.Val, 16)
	default:
		panic(fmt.Sprintf("invalid style (=%d)", t.Style))
	}
}

// String returns short textual representation (without position
// or category info) fot this token.
//
// Should be used mainly for error reporting.
func (t *Token) String() string {
	var s string
	var noq bool // true if result does not need quotes

	switch t.Kind {
	case INV, EOF:
		s = t.Kind.Category()
		noq = false
	case Word:
		s = t.Data
	case Rune:
		panic("TODO")
	case Integer:
		s = t.IntegerText()
	case String:
		s = char.Escape(t.Data)
	default:
		s = ktext[t.Kind]
		if s == "" {
			panic(fmt.Sprintf("invalid kind (=%d)", t.Kind))
		}
	}

	if noq {
		return s
	}

	return "\"" + s + "\""
}

// ktext maps Kind to token static ktext or code string.
var ktext = [...]string{
	// Operators/punctuators

	Under:   "_",
	AtSign:  "@",
	Address: ".&",
	Deref:   ".*",
	// DerefSelect:    ".*.",
	Plus:           "+",
	Minus:          "-",
	BoolAnd:        "&&",
	BoolOr:         "||",
	Equal:          "==",
	NotEqual:       "!=",
	LessOrEqual:    "<=",
	GreaterOrEqual: ">=",
	RightArrow:     "=>",
	Walrus:         ":=",
	// AddAssign:      "+=",
	// SubAssign:      "-=",
	// MulAssign:      "*=",
	// DivAssign:      "/=",
	// RemAssign:      "%=",
	// AndAssign:      "&=",
	// OrAssign:       "|=",
	Pipe: "|",
	// Caret:          "^",
	// LeftShift:      "<<",
	// RightShift:     ">>",
	// BitAndNot:      "&^",
	Assign: "=",
	Colon:  ":",
	// Ellipsis:       "...",
	Semicolon: ";",
	Asterisk:  "*",
	// Quest:          "?",
	Ampersand: "&",
	Exclam:    "!",
	Slash:     "/",
	// Percent:        "%",
	Period:      ".",
	Comma:       ",",
	LeftAngle:   "<",
	RightAngle:  ">",
	LeftCurly:   "{",
	RightCurly:  "}",
	LeftSquare:  "[",
	RightSquare: "]",
	LeftParen:   "(",
	RightParen:  ")",
	// HashCurly:      "#{",
	HashSquare: "#[",
	Alter:      ".{",
	AlterZero:  ".{}",
	SliceAll:   "[:]",
	// DerefIndex:     ".[",
	// BagSelect:      ".(",
	// Insist:         ".!",
	// Chain:          ".?",
	// PairSquare:     "[]",
	// AutoLen:        "[_]",
	ArrayPointer: "[*]",
	// ArrayRef:       "[&]",
	// CapBuf:         "[^]",

	// Keywords

	Import: "import",
	Fun:    "fun",
	Goto:   "goto",
	Return: "return",
	For:    "for",
	Else:   "else",
	If:     "if",
	Var:    "var",
	Type:   "type",
	Const:  "const",
	Union:  "union",
	Struct: "struct",
	Pub:    "pub",
	Zone:   "zone",

	// Special literals

	Nil:   "nil",
	True:  "true",
	False: "false",
}
