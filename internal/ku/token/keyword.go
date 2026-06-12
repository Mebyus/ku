package token

var kw = map[string]Kind{
	"fun":    Fun,
	"const":  Const,
	"return": Return,
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
