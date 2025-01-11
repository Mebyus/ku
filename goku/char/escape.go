package char

func NeedEscape(c byte) bool {
	switch c {
	case '\\', '\n', '\t', '\r', '"':
		return true
	default:
		return false
	}
}

func Escape(s string) string {
	// first scan the string to check if it needs at least one escape sequence
	i := 0
	for ; i < len(s); i += 1 {
		if NeedEscape(s[i]) {
			break
		}
	}
	if i >= len(s) {
		// return the string "as is", because it does not need escape sequences
		return s
	}

	// allocate enough space for at least one escape sequence
	buf := make([]byte, i, len(s)+1)
	copy(buf[:i], s[:i])

	for ; i < len(s); i += 1 {
		switch s[i] {
		case '\\':
			buf = append(buf, '\\', '\\')
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\t':
			buf = append(buf, '\\', 't')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '"':
			buf = append(buf, '\\', '"')
		default:
			buf = append(buf, s[i])
		}
	}

	return string(buf)
}
