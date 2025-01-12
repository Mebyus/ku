package token

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mebyus/ku/goku/compiler/char"
	"github.com/mebyus/ku/goku/compiler/source"
)

// literal maps Kind to token static literal or code string.
var literal = [...]string{
	empty: "<nil>",

	EOF: "EOF",

	// Operators/punctuators

	Underscore:     "_",
	Address:        ".&",
	Deref:          ".@",
	Plus:           "+",
	Minus:          "-",
	And:            "&&",
	Or:             "||",
	Equal:          "==",
	NotEqual:       "!=",
	LessOrEqual:    "<=",
	GreaterOrEqual: ">=",
	LeftArrow:      "<-",
	RightArrow:     "=>",
	Walrus:         ":=",
	AddAssign:      "+=",
	SubAssign:      "-=",
	MulAssign:      "*=",
	DivAssign:      "/=",
	RemAssign:      "%=",
	Pipe:           "|",
	Caret:          "^",
	LeftShift:      "<<",
	RightShift:     ">>",
	BitAndNot:      "&^",
	Assign:         "=",
	Colon:          ":",
	DoubleColon:    "::",
	Semicolon:      ";",
	Asterisk:       "*",
	Quest:          "?",
	Ampersand:      "&",
	Not:            "!",
	Slash:          "/",
	Percent:        "%",
	Period:         ".",
	Comma:          ",",
	LeftAngle:      "<",
	RightAngle:     ">",
	LeftCurly:      "{",
	RightCurly:     "}",
	LeftSquare:     "[",
	RightSquare:    "]",
	LeftParen:      "(",
	RightParen:     ")",
	PropStart:      "#[",
	Compound:       ".{",
	DerefIndex:     ".[",
	BagSelect:      ".(",
	Insist:         ".!",
	Chain:          ".?",
	Chunk:          "[]",
	AutoLen:        "[_]",
	ArrayPointer:   "[*]",
	CapBuffer:      "[^]",
	Nillable:       "?|",
	NillableChunk:  "[?]",
	DerefSelect:    ".@.",

	// Keywords

	Import: "import",
	Fun:    "fun",
	Jump:   "jump",
	Ret:    "ret",
	For:    "for",
	Else:   "else",
	If:     "if",
	Defer:  "defer",
	Bag:    "bag",
	In:     "in",
	Var:    "var",
	Type:   "type",
	Test:   "test",
	Enum:   "enum",
	Struct: "struct",
	Pub:    "pub",
	Unit:   "unit",
	Let:    "let",

	// Special literals

	Never: "never",
	Stub:  "stub",
	Dirty: "dirty",
	Nil:   "nil",
	True:  "true",
	False: "false",

	Cast: "#cast",
	Tint: "#tint",
	Size: "#size",

	Any: "any",

	LabelNext: "@.next",
	LabelOut:  "@.out",

	DirName:    "#name",
	DirInclude: "#include",
	DirDefine:  "#define",
	DirLink:    "#link",
	DirIf:      "#if",

	// Non static literals

	Illegal:    "ILG",
	Word:       "WORD",
	String:     "STR",
	RawString:  "STR.RAW",
	FillString: "STR.FILL",
	Rune:       "RUNE",
	BinInteger: "INT.BIN",
	OctInteger: "INT.OCT",
	DecInteger: "INT.DEC",
	HexInteger: "INT.HEX",
	DecFloat:   "FLT.DEC",
	Macro:      "MACRO",
	Env:        "ENV",

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
	case BinInteger:
		return "0b" + strconv.FormatUint(t.Val, 2)
	case OctInteger:
		return "0o" + strconv.FormatUint(t.Val, 8)
	case DecInteger:
		return strconv.FormatUint(t.Val, 10)
	case HexInteger:
		return "0x" + strconv.FormatUint(t.Val, 16)
	case DecFloat:
		return t.Data
	case Rune:
		if t.Data != "" {
			return "'" + t.Data + "'"
		}
		switch t.Val {
		case '\\':
			return `'\\'`
		case '\'':
			return `'\''`
		case '\n':
			return `'\n'`
		case '\t':
			return `'\t'`
		case '\r':
			return `'\r'`
		}
		return "'" + string([]rune{rune(t.Val)}) + "'"
	case String:
		return "\"" + char.Escape(t.Data) + "\""
	default:
		panic(fmt.Sprintf("unexpected \"%s\" token kind (=%d)", t.Kind, t.Kind))
	}
}

func FormatTokenLine(m source.PinMap, t Token) string {
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
