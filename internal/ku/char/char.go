package char

// IsCodePointStart returns true if given byte starts a utf-8 code point.
//
// utf-8 code point starts either from [0xxx xxxx] or [11xx xxxx].
// Non-starting byte of non-ascii code point has form [10xx xxxx].
// Thus we need to check that higher bits of a given byte are not 10.
func IsCodePointStart(c byte) bool {
	return c>>6 != 0b10
}

func IsSimpleWhitespace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}

// IsTextByte returns true if given byte represents ASCII printable or
// control character which holds meaning for modern text files and is
// compatible with utf-8.
func IsTextByte(b byte) bool {
	return (' ' <= b && b <= '~') || IsSimpleWhitespace(b)
}

const capitalLatinLetterMask = 0xDF

// CapitalLatinLetter transforms ASCII latin letter character to its upper
// (capital) form.
func CapitalLatinLetter(c byte) byte {
	return c & capitalLatinLetterMask
}

const lowerLatinLetterMask = 0x20

func LowerLatinLetter(c byte) byte {
	return c | lowerLatinLetterMask
}

func IsLatinLetter(c byte) bool {
	c = CapitalLatinLetter(c)
	return 'A' <= c && c <= 'Z'
}

func IsLatinLetterOrUnderscore(c byte) bool {
	return IsLatinLetter(c) || c == '_'
}

func IsAlphanum(c byte) bool {
	return IsLatinLetterOrUnderscore(c) || IsDecDigit(c)
}

func IsDecDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func IsHexDigit(c byte) bool {
	if IsDecDigit(c) {
		return true
	}
	c = CapitalLatinLetter(c)
	return 'A' <= c && c <= 'F'
}

func IsDecDigitOrPeriod(c byte) bool {
	return IsDecDigit(c) || c == '.'
}

func DecDigitNum(c byte) uint8 {
	return c - '0'
}

// HexDigitNum transforms ASCII hexadecimal digit character (0 - 9, A - F, a - f)
// to its number value. Small (a - f) and capital (A - F) letters produce
// identical results.
func HexDigitNum(c byte) uint8 {
	if c <= '9' {
		return DecDigitNum(c)
	}
	return CapitalLatinLetter(c) - 'A' + 10
}

func ParseHexDigits(s string) uint64 {
	var v uint64
	for i := range len(s) {
		v <<= 4
		v += uint64(HexDigitNum(s[i]))
	}
	return v
}
