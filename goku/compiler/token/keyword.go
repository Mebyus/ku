package token

var kw = map[string]Kind{
	"import": Import,
	"fun":    Fun,
	"test":   Test,
	"jump":   Jump,
	"ret":    Ret,
	"for":    For,
	"else":   Else,
	"if":     If,
	"defer":  Defer,
	"bag":    Bag,
	"in":     In,
	"var":    Var,
	"let":    Let,
	"type":   Type,
	"const":  Const,
	"union":  Union,
	"struct": Struct,
	"pub":    Pub,
	"gen":    Gen,
	"unsafe": Unsafe,

	"must":  Must,
	"panic": Panic,
	"cast":  Cast,
	"tint":  Tint,

	"nil":   Nil,
	"true":  True,
	"false": False,

	"any": Any,
}

const (
	minKeywordLength = 2
	maxKeywordLength = 6
)

// Keyword returns keyword kind if a given word is keyword.
//
// Returns (kind, true) if word is keyword.
// Returns (0, false) otherwise.
func Keyword(word string) (Kind, bool) {
	if len(word) < minKeywordLength || len(word) > maxKeywordLength {
		return 0, false
	}
	k, ok := kw[word]
	return k, ok
}
