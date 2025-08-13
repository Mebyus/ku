package baselex

import "strconv"

const (
	LengthOverflow = iota + 1
	NonPrintableByte

	MalformedString
	MalformedRune
	MalformedBinaryInteger
	MalformedOctalInteger
	MalformedDecimalInteger
	MalformedDecimalFloat
	MalformedHexadecimalInteger
	MalformedBlockComment
	MalformedMacro
	MalformedEnv

	DecimalIntegerOverflow

	BadEscapeInString

	UnknownDirective
)

func FormatErrorValue(v uint64) string {
	// TODO: human-readable text for errors
	return strconv.FormatUint(v, 10)
}

func (t *Token) SetIllegalError(code uint64) {
	t.Kind = Illegal
	t.Val = code
}
