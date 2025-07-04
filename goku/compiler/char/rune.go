package char

func EscapeRune(r rune) string {
	switch r {
	case '\\':
		return `\\`
	case '\'':
		return `\'`
	case '\n':
		return `\n`
	case '\t':
		return `\t`
	case '\r':
		return `\r`
	default:
		return string([]rune{r})
	}
}
