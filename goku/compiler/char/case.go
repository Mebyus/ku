package char

import "strings"

func SnakeCase(s string) string {
	var g strings.Builder
	for {
		var word string
		word, s = cutLatinWord(s)
		if word == "" {
			return g.String()
		}

		if g.Len() != 0 {
			g.WriteByte('_')
		}
		g.Grow(len(word))
		for i := range len(word) {
			g.WriteByte(LowerLatinLetter(word[i]))
		}
	}
}

func cutLatinWord(s string) (word, rest string) {
	if s == "" {
		return "", ""
	}

	var i int
	for i < len(s) && !IsLatinLetter(s[i]) {
		i += 1
	}
	start := i

	lower := 0 // number of encoutered lower letters
	for i < len(s) && IsLatinLetter(s[i]) {
		c := s[i]
		if c == LowerLatinLetter(c) {
			lower += 1
		} else if lower != 0 {
			break
		}

		i += 1
	}
	end := i
	return s[start:end], s[end:]
}
