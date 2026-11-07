package token

var kw = map[string]Kind{
	"if":     If,
	"var":    Var,
	"pub":    Pub,
	"fun":    Fun,
	"for":    For,
	"type":   Type,
	"zone":   Zone,
	"else":   Else,
	"const":  Const,
	"import": Import,
	"return": Return,
	"goto":   Goto,
	"struct": Struct,
	"union":  Union,
	"true":   True,
	"false":  False,
	"nil":    Nil,
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
