package char

func ParseDecDigits(s string) uint64 {
	var v uint64
	for i := range len(s) {
		d := s[i]
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
func ParseDecDigitsWithOverflowCheck(s string) (uint64, bool) {
	if len(s) > maxUint64DecLength {
		return 0, false
	}
	if len(s) < maxUint64DecLength {
		return ParseDecDigits(s), true
	}

	var v uint64
	for i := range len(s) {
		d := s[i]
		if d > maxUint64DecString[i] {
			return 0, false
		}
		v *= 10
		v += uint64(DecDigitNum(d))
	}
	return v, true
}
