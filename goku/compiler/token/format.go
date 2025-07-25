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

	Underscore:     "_",
	Address:        ".&",
	Deref:          ".*",
	DerefSelect:    ".*.",
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
	AndAssign:      "&=",
	OrAssign:       "|=",
	Pipe:           "|",
	Caret:          "^",
	LeftShift:      "<<",
	RightShift:     ">>",
	BitAndNot:      "&^",
	Assign:         "=",
	Colon:          ":",
	Ellipsis:       "...",
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
	HashCurly:      "#{",
	HashSquare:     "#[",
	Tweak:          ".{",
	DerefIndex:     ".[",
	BagSelect:      ".(",
	Insist:         ".!",
	Chain:          ".?",
	Chunk:          "[]",
	AutoLen:        "[_]",
	ArrayPointer:   "[*]",
	ArrayRef:       "[&]",
	CapBuffer:      "[^]",
	Nillable:       "?|",
	NillableChunk:  "[?]",

	LeftShiftAssign:  "<<=",
	RightShiftAssign: ">>=",

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
	Const:  "const",
	Test:   "test",
	Union:  "union",
	Struct: "struct",
	Pub:    "pub",
	Unsafe: "unsafe",
	Let:    "let",
	Gen:    "gen",
	Must:   "must",
	Panic:  "panic",
	Cast:   "cast",
	Tint:   "tint",

	// Special literals

	Nil:   "nil",
	True:  "true",
	False: "false",
	Any:   "any",

	StaticMust: "#must",

	Debug:  "#debug",
	Build:  "#build",
	Stub:   "#stub",
	Never:  "#never",
	Size:   "#size",
	Lookup: "#lookup",

	TypeId:  "#typeid",
	ErrorId: "#error",
	Enum:    "#enum",

	Check: "#check",
	Len:   "#len",

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
