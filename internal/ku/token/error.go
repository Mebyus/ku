package token

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

	DecimalIntegerOverflow
)

func (t *Token) SetError(code uint64) {
	t.Kind = INV
	t.Val = code
}
