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

// Unescape data from a given string.
//
// Returns (string, true) if unescape was successful.
// Otherwise (on bad escape sequence) returns ("", false).
func Unescape(s string) (string, bool) {
	// first scan the string to check if it has at least one escape sequence
	i := 0
	for ; i < len(s); i++ {
		if s[i] == '\\' {
			// escape sequence found
			break
		}
	}
	if i >= len(s) {
		// return the string "as is", because it does not have escape sequences
		return s, true
	}

	buf := make([]byte, i, len(s))
	copy(buf[:i], s[:i])

	for ; i < len(s); i++ {
		if s[i] == '\\' {
			i += 1
			var b byte
			switch s[i] {
			case '\\':
				b = '\\'
			case 'n':
				b = '\n'
			case 't':
				b = '\t'
			case 'r':
				b = '\r'
			case '"':
				b = '"'
			default:
				return "", false
			}
			buf = append(buf, b)
		} else {
			buf = append(buf, s[i])
		}
	}

	return string(buf), true
}
