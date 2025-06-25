package token

var kw = map[string]Kind{
	"import":  Import,
	"test":    Test,
	"else":    Else,
	"if":      If,
	"include": Include,
	"set":     Set,
	"unit":    Unit,
	"main":    Main,
	"module":  Module,
	"link":    Link,
}

const (
	minKeywordLength = 2
	maxKeywordLength = 7
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
