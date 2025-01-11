package char

// utf-8 code point starts either from [0xxx xxxx] or [11xx xxxx].
// Non-starting byte of non-ascii code point has from [10xx xxxx].
// Thus we need to check that higher bits of a given byte are not 10.
func IsCodePointStart(c byte) bool {
	return c>>6 != 0b10
}

// InsideCodePoint returns true if byte can represent non-starting byte
// of utf-8 non-ascii code point.
func InsideCodePoint(c byte) bool {
	return c>>6 == 0b10
}
