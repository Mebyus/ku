package token

const (
	LengthOverflow = iota + 1
	NonPrintableByte

	MalformedString
	MalformedDecimalInteger
	MalformedHexadecimalInteger
	MalformedBlockComment

	DecimalIntegerOverflow
)

func (t *Token) SetIllegalError(code uint64) {
	t.Kind = Illegal
	t.Val = code
}
