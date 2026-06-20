package token

var kw = map[string]Kind{
	"if":     If,
	"fun":    Fun,
	"type":   Type,
	"else":   Else,
	"const":  Const,
	"import": Import,
	"return": Return,
	"struct": Struct,
	"true":   True,
	"false":  False,
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
