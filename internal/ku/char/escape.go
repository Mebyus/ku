package char

import "strings"

func needEscape(c byte) bool {
	switch c {
	case '\\', '\n', '\t', '\r', '"':
		return true
	default:
		return !IsTextByte(c)
	}
}

// Escape accepts arbitrary string and
func Escape(s string) string {
	// first scan the string to check if it needs at least one escape sequence
	i := 0
	for ; i < len(s); i += 1 {
		if needEscape(s[i]) {
			break
		}
	}
	if i >= len(s) {
		// return the string "as is", because it does not need escape sequences
		return s
	}

	// allocate enough space for at least one escape sequence
	var g strings.Builder
	g.Grow(len(s) + 1)
	g.WriteString(s[:i])

	for ; i < len(s); i += 1 {
		switch s[i] {
		case '\\':
			g.WriteByte('\\')
			g.WriteByte('\\')
		case '\n':
			g.WriteByte('\\')
			g.WriteByte('n')
		case '\t':
			g.WriteByte('\\')
			g.WriteByte('t')
		case '\r':
			g.WriteByte('\\')
			g.WriteByte('r')
		case '"':
			g.WriteByte('\\')
			g.WriteByte('"')
		default:
			if IsTextByte(s[i]) {
				g.WriteByte(s[i])
			} else {
				g.WriteString(`\x`)
				hi := s[i] >> 4
				lo := s[i] & 0xF
				g.WriteByte(HexDigit(hi))
				g.WriteByte(HexDigit(lo))
			}
		}
	}

	return g.String()
}

// Unescape data from a given string. Unescape takes a given string
// scans it for escape sequences and forms a new string if necessary.
// In new string all escape sequences are replaced with their raw byte
// representation. Bad or unrecognized escape sequences are skipped.
//
// Returns (string, true) if unescape was successful (no bad escape sequences).
// Otherwise (with at least one bad escape sequence) returns (string, false),
// where returned string is formed by applying valid escape
// sequences and discarding bad ones.
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

	var g strings.Builder
	g.Grow(len(s))
	g.WriteString(s[:i])

	bad := 0 // total number of bad escape sequences
	for ; i < len(s); i++ {
		if s[i] != '\\' {
			g.WriteByte(s[i])
			continue
		}

		i += 1
		if i >= len(s) {
			// we found escape slash at the of the string
			// without actual sequence
			return g.String(), false
		}

		var b byte
		switch s[i] {
		case '\\':
			b = '\\'
		case '0':
			b = 0
		case 'n':
			b = '\n'
		case 't':
			b = '\t'
		case 'r':
			b = '\r'
		case '"':
			b = '"'
		case '\'':
			b = '\''
		case 'x':
			i += 1
			if i+2 > len(s) {
				bad += 1
				return g.String(), false
			}
			if !IsHexDigit(s[i]) || !IsHexDigit(s[i+1]) {
				i += 1
				bad += 1
				continue
			}
			b = byte(ParseHexDigits(s[i : i+2]))
			i += 1
		default:
			bad += 1
			continue
		}
		g.WriteByte(b)
	}

	return g.String(), bad == 0
}
