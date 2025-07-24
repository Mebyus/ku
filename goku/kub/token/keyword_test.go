package token

import "testing"

func TestKeyword(t *testing.T) {
	minLen := 1 << 10 // arbitrary large number
	maxLen := 0

	for word, kind := range kw {
		if len(word) > maxLen {
			maxLen = len(word)
		}
		if len(word) < minLen {
			minLen = len(word)
		}

		lit := literal[kind]
		if lit != word {
			t.Errorf("keyword \"%s\" has inconsistent literal \"%s\"", word, lit)
		}
	}

	if minLen != minKeywordLength {
		t.Errorf("min keyword length should be %d, not %d", minLen, minKeywordLength)
	}
	if maxLen != maxKeywordLength {
		t.Errorf("max keyword length should be %d, not %d", maxLen, maxKeywordLength)
	}
}
