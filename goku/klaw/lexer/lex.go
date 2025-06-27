package lexer

import (
	"github.com/mebyus/ku/goku/compiler/char"
	"github.com/mebyus/ku/goku/klaw/token"
)

func (lx *Lexer) Lex() token.Token {
	if lx.eof() {
		return lx.emit(token.EOF)
	}

	lx.skipWhitespaceAndComments()
	if lx.eof() {
		return lx.emit(token.EOF)
	}

	return lx.lex()
}

func (lx *Lexer) lex() token.Token {
	if char.IsLatinLetterOrUnderscore(lx.peek()) {
		return lx.word()
	}

	if char.IsDecDigit(lx.peek()) {
		return lx.number()
	}

	if lx.peek() == '"' {
		return lx.str()
	}

	return lx.other()
}

func (lx *Lexer) decNumber() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.start()
	lx.advance() // skip first digit
	scannedOnePeriod := false
	for !lx.eof() && char.IsDecDigitOrPeriod(lx.peek()) {
		if lx.peek() == '.' {
			if scannedOnePeriod || !char.IsDecDigit(lx.next()) {
				data, ok := lx.take()
				if ok {
					tok.SetIllegalError(token.DecimalIntegerOverflow)
					tok.Data = data
				} else {
					tok.SetIllegalError(token.LengthOverflow)
				}
				return tok
			} else {
				scannedOnePeriod = true
			}
		}
		lx.advance()
	}

	if lx.isLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	if !lx.eof() && char.IsAlphanum(lx.peek()) {
		lx.skipWord()
		data, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedDecimalInteger)
			tok.Data = data
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return tok
	}

	if !scannedOnePeriod {
		// decimal integer
		n, ok := char.ParseDecDigitsWithOverflowCheck(lx.view())
		if !ok {
			tok.SetIllegalError(token.DecimalIntegerOverflow)
			return tok
		}

		tok.Kind = token.DecInteger
		tok.Val = n
		return tok
	}

	tok.SetIllegalError(token.MalformedDecimalInteger)
	return tok
}

func (lx *Lexer) hexNumber() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.advance() // skip "0"
	lx.advance() // skip "x"

	lx.start()
	lx.skipHexDigits()

	if char.IsAlphanum(lx.peek()) {
		lx.skipWord()
		data, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedHexadecimalInteger)
			tok.Data = data
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return tok
	}

	if lx.isLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}
	if lx.length() == 0 {
		tok.SetIllegalError(token.MalformedHexadecimalInteger)
		tok.Data = "0x"
		return tok
	}

	tok.Kind = token.HexInteger
	if lx.length() > 16 {
		lit, ok := lx.take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Data = lit
		return tok
	}

	tok.Val = char.ParseHexDigits(lx.view())
	return tok
}

func (lx *Lexer) number() (tok token.Token) {
	if lx.peek() != '0' {
		return lx.decNumber()
	}

	if lx.next() == 'x' {
		return lx.hexNumber()
	}

	if lx.next() == '.' {
		return lx.decNumber()
	}

	if char.IsAlphanum(lx.next()) {
		return lx.illegalWord(token.MalformedDecimalInteger)
	}

	tok = token.Token{
		Kind: token.DecInteger,
		Pin:  lx.pin(),
		Val:  0,
	}
	lx.advance()
	return
}

func (lx *Lexer) word() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	if !char.IsAlphanum(lx.next()) {
		// word is 1 character long
		c := lx.peek()
		lx.advance() // skip single (start) character

		tok.Kind = token.Word
		tok.Data = char.ToString(c)
		return tok
	}

	// word is at least 2 characters long
	lx.start()
	lx.skipWord()
	word, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	kind, ok := token.Keyword(word)
	if ok {
		tok.Kind = kind
		return tok
	}

	tok.Kind = token.Word
	tok.Data = word
	return tok
}

// Scan string literal.
func (lx *Lexer) str() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.advance() // skip quote
	if lx.eof() {
		tok.SetIllegalError(token.MalformedString)
		tok.Data = "\""
		return tok
	}

	if lx.peek() == '"' {
		// common case of empty string literal
		lx.advance() // skip quote
		tok.Kind = token.String
		return tok
	}

	var fills uint64 // number of fill places inside the string
	lx.start()
	for !lx.eof() && lx.peek() != '"' && lx.peek() != '\n' {
		if lx.peek() == '\\' && lx.next() == '"' {
			// do not stop if we encounter escape sequence
			lx.advance() // skip "\"
			lx.advance() // skip quote
		} else if lx.peek() == '$' && lx.next() == '{' {
			fills += 1
			lx.advance() // skip "$"
			lx.advance() // skip "{"
		} else {
			lx.advance()
		}
	}

	if lx.eof() {
		data, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedString)
			tok.Data = data
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return tok
	}

	if lx.peek() != '"' {
		data, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedString)
			tok.Data = data
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return tok
	}

	data, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	lx.advance() // skip quote

	if fills != 0 {
		tok.SetIllegalError(token.MalformedString)
		return tok
	}

	data, ok = char.Unescape(data)
	if !ok {
		tok.SetIllegalError(token.MalformedString)
		return tok
	}

	tok.Kind = token.String
	tok.Data = data
	return tok
}

func (lx *Lexer) other() token.Token {
	switch lx.peek() {
	case '(':
		return lx.oneByteToken(token.LeftParen)
	case ')':
		return lx.oneByteToken(token.RightParen)
	case '{':
		return lx.oneByteToken(token.LeftCurly)
	case '}':
		return lx.oneByteToken(token.RightCurly)
	case '<':
		switch lx.next() {
		case '=':
			return lx.twoBytesToken(token.LessOrEqual)
		default:
			return lx.oneByteToken(token.LeftAngle)
		}
	case '>':
		switch lx.next() {
		case '=':
			return lx.twoBytesToken(token.GreaterOrEqual)
		default:
			return lx.oneByteToken(token.RightAngle)
		}
	case '=':
		switch lx.next() {
		case '=':
			return lx.twoBytesToken(token.Equal)
		default:
			return lx.oneByteToken(token.Assign)
		}
	case ';':
		return lx.oneByteToken(token.Semicolon)
	case '.':
		return lx.oneByteToken(token.Period)
	case '!':
		if lx.next() == '=' {
			return lx.twoBytesToken(token.NotEqual)
		}
		return lx.oneByteToken(token.Not)
	default:
		return lx.illegalByteToken()
	}
}

func (lx *Lexer) oneByteToken(k token.Kind) token.Token {
	tok := lx.emit(k)
	lx.advance()
	return tok
}

func (lx *Lexer) twoBytesToken(k token.Kind) token.Token {
	tok := lx.emit(k)
	lx.advance()
	lx.advance()
	return tok
}

func (lx *Lexer) illegalByteToken() token.Token {
	tok := lx.emit(token.Illegal)
	tok.Data = char.ToString(byte(lx.peek()))
	lx.advance()
	return tok
}

func (lx *Lexer) illegalWord(code uint64) token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.start()
	lx.skipWord()
	data, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	tok.SetIllegalError(code)
	tok.Data = data
	return tok
}
