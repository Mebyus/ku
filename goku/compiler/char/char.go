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

func IsSimpleWhitespace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\r'
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

func IsDecDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func IsOctDigit(c byte) bool {
	return '0' <= c && c <= '7'
}

func IsBinDigit(c byte) bool {
	return c == '0' || c == '1'
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

func IsAlphanum(c byte) bool {
	return IsLatinLetterOrUnderscore(c) || IsDecDigit(c)
}

func ToString(b byte) string {
	return string([]byte{b})
}

// ParseBinDigits interprets ASCII digit characters as digits of binary number
// and returns the number.
//
// Does not validate the input.
// Input slice must satisfy 1 <= len(s) <= 64.
// As a special case returns zero for nil or empty input slice.
func ParseBinDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v <<= 1
		v += uint64(DecDigitNum(d))
	}
	return v
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

func ParseHexDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v <<= 4
		v += uint64(HexDigitNum(d))
	}
	return v
}

func ParseOctDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v <<= 3
		v += uint64(DecDigitNum(d))
	}
	return v
}

func ParseDecDigits(s []byte) uint64 {
	var v uint64
	for _, d := range s {
		v *= 10
		v += uint64(DecDigitNum(d))
	}
	return v
}

const maxUint64DecString = "18446744073709551615"

// Maximum length of unsigned 64-bit integer formatted as decimal number.
const maxUint64DecLength = len(maxUint64DecString)

// ParseDecDigitsWithOverflowCheck interprets each byte in the slice as ASCII
// character for decimal digit of an integer number and returns the resulting integer.
//
// If integer does not fit into 64-bit unsigned integer, then (0, false) is returned.
// Otherwise returns (n, true).
//
// Examples:
//
//	['1', '2', '0'] => 120
//	['0']           => 0
func ParseDecDigitsWithOverflowCheck(s []byte) (uint64, bool) {
	if len(s) > maxUint64DecLength {
		return 0, false
	}
	if len(s) < maxUint64DecLength {
		return ParseDecDigits(s), true
	}

	var v uint64
	for i, d := range s {
		if d > maxUint64DecString[i] {
			return 0, false
		}
		v *= 10
		v += uint64(DecDigitNum(d))
	}
	return v, true
}
