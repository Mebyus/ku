package baselex

import "github.com/mebyus/ku/goku/compiler/char"

func (lx *Lexer) Number() Token {
	if lx.Peek() != '0' {
		return lx.DecNumber()
	}

	if lx.Next() == 'b' {
		return lx.BinNumber()
	}

	if lx.Next() == 'o' {
		return lx.OctNumber()
	}

	if lx.Next() == 'x' {
		return lx.HexNumber()
	}

	if lx.Next() == '.' {
		return lx.DecNumber()
	}

	if char.IsAlphanum(lx.Next()) {
		return lx.IllegalWord(MalformedDecimalInteger)
	}

	tok := Token{
		Kind: DecInteger,
		Pin:  lx.Pin(),
		Val:  0,
	}
	lx.Advance()
	return tok
}

func (lx *Lexer) DecNumber() (tok Token) {
	tok.Pin = lx.Pin()

	lx.Start()
	lx.Advance() // skip first digit
	scannedOnePeriod := false
	for !lx.Eof() && char.IsDecDigitOrPeriod(lx.Peek()) {
		if lx.Peek() == '.' {
			if scannedOnePeriod || !char.IsDecDigit(lx.Next()) {
				data, ok := lx.Take()
				if ok {
					tok.SetIllegalError(MalformedDecimalFloat)
					tok.Data = data
				} else {
					tok.SetIllegalError(LengthOverflow)
				}
				return tok
			} else {
				scannedOnePeriod = true
			}
		}
		lx.Advance()
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(LengthOverflow)
		return tok
	}

	if !lx.Eof() && char.IsAlphanum(lx.Peek()) {
		lx.SkipWord()
		data, ok := lx.Take()
		if ok {
			tok.SetIllegalError(MalformedDecimalInteger)
			tok.Data = data
		} else {
			tok.SetIllegalError(LengthOverflow)
		}
		return tok
	}

	if !scannedOnePeriod {
		// decimal integer
		n, ok := char.ParseDecDigitsWithOverflowCheck(lx.View())
		if !ok {
			tok.SetIllegalError(DecimalIntegerOverflow)
			return tok
		}

		tok.Kind = DecInteger
		tok.Val = n
		return tok
	}

	tok.Kind = DecFloat
	tok.Data, _ = lx.Take()
	return tok
}

func (lx *Lexer) BinNumber() (tok Token) {
	tok.Pin = lx.Pin()

	lx.Advance() // skip '0'
	lx.Advance() // skip 'b'

	lx.Start()
	lx.SkipBinDigits()

	if char.IsAlphanum(lx.Peek()) {
		lx.SkipWord()
		data, ok := lx.Take()
		if ok {
			tok.SetIllegalError(MalformedBinaryInteger)
			tok.Data = data
		} else {
			tok.SetIllegalError(LengthOverflow)
		}
		return tok
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(LengthOverflow)
		return tok
	}
	if lx.Length() == 0 {
		tok.SetIllegalError(MalformedBinaryInteger)
		tok.Data = "0b"
		return tok
	}

	tok.Kind = BinInteger
	if lx.Length() > 64 {
		lit, ok := lx.Take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Data = lit
		return tok
	}

	tok.Val = char.ParseBinDigits(lx.View())
	return tok
}

func (lx *Lexer) OctNumber() (tok Token) {
	tok.Pin = lx.Pin()

	lx.Advance() // skip '0' byte
	lx.Advance() // skip 'o' byte

	lx.Start()
	lx.SkipOctDigits()

	if char.IsAlphanum(lx.Peek()) {
		lx.SkipWord()
		data, ok := lx.Take()
		if ok {
			tok.SetIllegalError(MalformedOctalInteger)
			tok.Data = data
		} else {
			tok.SetIllegalError(LengthOverflow)
		}
		return tok
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(LengthOverflow)
		return tok
	}
	if lx.Length() == 0 {
		tok.SetIllegalError(MalformedOctalInteger)
		tok.Data = "0o"
		return tok
	}

	tok.Kind = OctInteger
	if lx.Length() > 21 {
		data, ok := lx.Take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Data = data
		return tok
	}

	tok.Val = char.ParseOctDigits(lx.View())
	return tok
}

func (lx *Lexer) HexNumber() (tok Token) {
	tok.Pin = lx.Pin()

	lx.Advance() // skip "0"
	lx.Advance() // skip "x"

	lx.Start()
	lx.SkipHexDigits()

	if char.IsAlphanum(lx.Peek()) {
		lx.SkipWord()
		data, ok := lx.Take()
		if ok {
			tok.SetIllegalError(MalformedHexadecimalInteger)
			tok.Data = data
		} else {
			tok.SetIllegalError(LengthOverflow)
		}
		return tok
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(LengthOverflow)
		return tok
	}
	if lx.Length() == 0 {
		tok.SetIllegalError(MalformedHexadecimalInteger)
		tok.Data = "0x"
		return tok
	}

	tok.Kind = HexInteger
	if lx.Length() > 16 {
		lit, ok := lx.Take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Data = lit
		return tok
	}

	tok.Val = char.ParseHexDigits(lx.View())
	return tok
}

func (lx *Lexer) IllegalWord(code uint64) (tok Token) {
	tok.Pin = lx.Pin()

	lx.Start()
	lx.SkipWord()
	data, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(LengthOverflow)
		return tok
	}

	tok.SetIllegalError(code)
	tok.Data = data
	return tok
}
