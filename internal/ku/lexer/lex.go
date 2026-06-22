package lexer

import (
	"github.com/mebyus/ku/internal/ku/char"
	"github.com/mebyus/ku/internal/ku/token"
)

// Lex reads next token from source text and writes it to given *Token.
// Token passed by pointer argument must be zeroed before call.
func (lx *Lexer) Lex(tok *token.Token) {
	if lx.eof() {
		lx.emitEOF(tok)
		return
	}

	lx.skipWhitespaceAndComments()
	if lx.eof() {
		lx.emitEOF(tok)
		return
	}

	tok.Pin = lx.pin()
	lx.lex(tok)
}

func (lx *Lexer) lex(tok *token.Token) {
	if char.IsLatinLetterOrUnderscore(lx.peek()) {
		lx.word(tok)
		return
	}

	if lx.peek() == '"' {
		lx.str(tok)
		return
	}

	if lx.peek() == '\'' {
		lx.rune(tok)
		return
	}


	if char.IsDecDigit(lx.peek()) {
		lx.number(tok)
		return
	}

	lx.other(tok)
}

func (lx *Lexer) word(tok *token.Token) {
	lx.start()
	lx.skipWord()
	word, ok := lx.take()
	if !ok {
		tok.SetError(token.LengthOverflow)
		return
	}

	kind, ok := token.Keyword(word)
	if ok {
		tok.Kind = kind
		return
	}

	tok.Kind = token.Word
	tok.Data = word
}

func (lx *Lexer) number(tok *token.Token) {
	if lx.peek() != '0' {
		lx.decNumber(tok)
		return
	}

	if lx.next() == 'b' {
		lx.binNumber(tok)
		return
	}

	if lx.next() == 'o' {
		lx.octNumber(tok)
		return
	}

	if lx.next() == 'x' {
		lx.hexNumber(tok)
		return
	}

	if lx.next() == '.' {
		lx.decNumber(tok)
		return
	}

	if char.IsAlphanum(lx.next()) {
		lx.emitIllegalWord(tok, token.MalformedDecimalInteger)
		return
	}

	// token is standalone number 0
	lx.advance()
	tok.Kind = token.Integer
	tok.Val = 0
	tok.Flags = token.DecInt
}

func (lx *Lexer) decNumber(tok *token.Token) {
	lx.start()
	lx.advance() // skip first digit
	scannedOnePeriod := false
	for !lx.eof() && char.IsDecDigitOrPeriod(lx.peek()) {
		if lx.peek() == '.' {
			if scannedOnePeriod || !char.IsDecDigit(lx.next()) {
				data, ok := lx.take()
				if ok {
					tok.SetError(token.MalformedDecimalFloat)
					tok.Data = data
				} else {
					tok.SetError(token.LengthOverflow)
				}
				return
			} else {
				scannedOnePeriod = true
			}
		}
		lx.advance()
	}

	if lx.isLengthOverflow() {
		tok.SetError(token.LengthOverflow)
		return
	}

	if !lx.eof() && char.IsAlphanum(lx.peek()) {
		lx.skipWord()
		data, ok := lx.take()
		if ok {
			tok.SetError(token.MalformedDecimalInteger)
			tok.Data = data
		} else {
			tok.SetError(token.LengthOverflow)
		}
		return
	}

	if !scannedOnePeriod {
		// decimal integer
		n, ok := char.ParseDecDigitsWithOverflowCheck(lx.view())
		if !ok {
			tok.SetError(token.DecimalIntegerOverflow)
			return
		}

		tok.Kind = token.Integer
		tok.Val = n
		tok.Flags = token.DecInt
		return
	}

	panic("stub")
	// tok.Kind = DecFloat
	// tok.Data, _ = lx.Take()
	// return tok
}

func (lx *Lexer) binNumber(tok *token.Token) {

}

func (lx *Lexer) octNumber(tok *token.Token) {

}

func (lx *Lexer) hexNumber(tok *token.Token) {
	lx.advance() // skip "0"
	lx.advance() // skip "x"

	lx.start()
	lx.skipHexDigits()

	if char.IsAlphanum(lx.peek()) {
		lx.skipWord()
		data, ok := lx.take()
		if ok {
			tok.SetError(token.MalformedHexadecimalInteger)
			tok.Data = data
		} else {
			tok.SetError(token.LengthOverflow)
		}
		return
	}

	if lx.isLengthOverflow() {
		tok.SetError(token.LengthOverflow)
		return
	}
	if lx.length() == 0 {
		tok.SetError(token.MalformedHexadecimalInteger)
		tok.Data = "0x"
		return
	}

	tok.Kind = token.Integer
	if lx.length() > 16 {
		lit, ok := lx.take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Data = lit
		return
	}

	tok.Val = char.ParseHexDigits(lx.view())
	tok.Flags = token.HexInt
}

func (lx *Lexer) rune(tok *token.Token) {
	lx.advance() // skip "'"
	if lx.eof() {
		tok.SetError(token.MalformedRune)
		tok.Data = "'"
		return
	}

	lx.start()
	if lx.peek() == '\\' {
		// handle escape sequence
		var val uint64
		switch lx.next() {
		case '\\':
			val = '\\'
		case 'n':
			val = '\n'
		case 't':
			val = '\t'
		case 'r':
			val = '\r'
		case '\'':
			val = '\''
		default:
			lx.advance() // skip "\"
			lx.advance() // skip unknown escape rune

			if !lx.eof() && lx.peek() == '\'' {
				lx.advance()
			}
			tok.SetError(token.MalformedRune)
			tok.Data, _ = lx.take()
			return
		}

		lx.advance() // skip "\"
		lx.advance() // skip escape rune
		if lx.eof() {
			tok.SetError(token.MalformedRune)
			tok.Data, _ = lx.take()
			return
		}
		if lx.peek() != '\'' {
			lx.advance()
			tok.SetError(token.MalformedRune)
			tok.Data, _ = lx.take()
			return
		}

		lx.advance() // skip "'"
		tok.Kind = token.Rune
		tok.Val = val
		return
	}

	if lx.next() == '\'' {
		// common case of ascii rune
		tok.Val = uint64(lx.peek())
		tok.Kind = token.Rune
		lx.advance() // skip rune character
		lx.advance() // skip "'"
		return
	}

	// handle non-ascii runes
	for !lx.eof() && lx.peek() != '\'' && lx.peek() != '\n' {
		lx.advance()
	}

	data, ok := lx.take()
	if !ok {
		tok.SetError(token.LengthOverflow)
		return
	}

	if lx.eof() || lx.peek() != '\'' {
		tok.SetError(token.MalformedRune)
		tok.Data = data
		return
	}
	lx.advance() // skip "'"

	runes := []rune(data)
	if len(runes) != 1 {
		tok.SetError(token.MalformedRune)
		tok.Data = data
		return
	}

	tok.Kind = token.Rune
	tok.Val = uint64(runes[0])
}

func (lx *Lexer) str(tok *token.Token) {
	lx.advance() // skip quote
	if lx.eof() {
		tok.SetError(token.MalformedString)
		tok.Data = "\""
		return
	}

	if lx.peek() == '"' {
		// common case of empty string literal
		lx.advance() // skip quote
		tok.Kind = token.String
		return
	}

	var esc uint64 // number of escapes inside the string
	lx.start()
	for !lx.eof() && lx.peek() != '"' && lx.peek() != '\n' {
		if lx.peek() == '\\' {
			esc += 1
		}

		if lx.peek() == '\\' && lx.next() == '"' {
			// do not stop if we encounter escape sequence
			lx.advance() // skip "\"
			lx.advance() // skip quote
		} else {
			lx.advance()
		}
	}

	if lx.eof() {
		data, ok := lx.take()
		if ok {
			tok.SetError(token.MalformedString)
			tok.Data = data
		} else {
			tok.SetError(token.LengthOverflow)
		}
		return
	}

	if lx.peek() != '"' {
		data, ok := lx.take()
		if ok {
			tok.SetError(token.MalformedString)
			tok.Data = data
		} else {
			tok.SetError(token.LengthOverflow)
		}
		return
	}

	data, ok := lx.take()
	if !ok {
		tok.SetError(token.LengthOverflow)
		return
	}

	lx.advance() // skip quote

	if esc != 0 {
		// data, ok = char.Unescape(data)
		// if !ok {
		// 	tok.SetError(token.BadEscapeInString)
		// 	return
		// }
		panic("stub")
	}

	tok.Kind = token.String
	tok.Data = data
}

func (lx *Lexer) other(tok *token.Token) {
	switch lx.peek() {
	case '(':
		lx.emitOneByteToken(tok, token.LeftParen)
	case ')':
		lx.emitOneByteToken(tok, token.RightParen)
	case '{':
		lx.emitOneByteToken(tok, token.LeftCurly)
	case '}':
		lx.emitOneByteToken(tok, token.RightCurly)
	case '[':
		lx.emitOneByteToken(tok, token.LeftSquare)
	case ']':
		lx.emitOneByteToken(tok, token.RightSquare)
	case ';':
		lx.emitOneByteToken(tok, token.Semicolon)
	case ':':
		if lx.next() == '=' {
			lx.emitTwoBytesToken(tok, token.Walrus)
			return
		}
		lx.emitOneByteToken(tok, token.Colon)
	case ',':
		lx.emitOneByteToken(tok, token.Comma)
	case '.':
		if lx.next() == '{' {
			lx.advance() // skip "."
			if lx.next() == '}' {
				lx.advance() // skip "{"
				lx.advance() // skip "}"
				tok.Kind = token.AlterZero
				return
			}
			lx.advance() // skip "{"
			tok.Kind = token.Alter
			return
		}

		lx.emitOneByteToken(tok, token.Period)
	case '|':
		if lx.next() == '|' {
			lx.emitTwoBytesToken(tok, token.BoolOr)
			return
		}
		lx.emitOneByteToken(tok, token.Pipe)
	case '&':
		if lx.next() == '&' {
			lx.emitTwoBytesToken(tok, token.BoolAnd)
			return
		}
		lx.emitOneByteToken(tok, token.Ampersand)
	case '*':
		lx.emitOneByteToken(tok, token.Asterisk)
	case '/':
		lx.emitOneByteToken(tok, token.Slash)
	case '+':
		lx.emitOneByteToken(tok, token.Plus)
	case '-':
		if lx.next() == '>' {
			lx.emitTwoBytesToken(tok, token.RightArrow)
			return
		}
		lx.emitOneByteToken(tok, token.Minus)
	case '=':
		if lx.next() == '=' {
			lx.emitTwoBytesToken(tok, token.Equal)
			return
		}
		lx.emitOneByteToken(tok, token.Assign)
	case '!':
		if lx.next() == '=' {
			lx.emitTwoBytesToken(tok, token.NotEqual)
			return
		}
		lx.emitOneByteToken(tok, token.Exclam)
	case '<':
		if lx.next() == '=' {
			lx.emitTwoBytesToken(tok, token.LessOrEqual)
			return
		}
		lx.emitOneByteToken(tok, token.LeftAngle)
	case '>':
		if lx.next() == '=' {
			lx.emitTwoBytesToken(tok, token.GreaterOrEqual)
			return
		}
		lx.emitOneByteToken(tok, token.RightAngle)
	default:
		lx.emitInvalidBytesToken(tok)
	}
}

func (lx *Lexer) emitOneByteToken(tok *token.Token, kind token.Kind) {
	tok.Kind = kind

	lx.advance()
}

func (lx *Lexer) emitTwoBytesToken(tok *token.Token, kind token.Kind) {
	tok.Kind = kind

	lx.advance()
	lx.advance()
}

func (lx *Lexer) emitIllegalWord(tok *token.Token, code uint64) {
	lx.start()
	lx.skipWord()
	data, ok := lx.take()
	if !ok {
		tok.SetError(token.LengthOverflow)
		return
	}

	tok.SetError(code)
	tok.Data = data
}

func (lx *Lexer) emitInvalidBytesToken(tok *token.Token) {
	tok.Val = uint64(lx.peek())
	tok.SetError(token.NonPrintableByte)

	lx.advance() // enshure we consume first invalid byte even if it is printable
	for !lx.eof() && !char.IsTextByte(lx.peek()) {
		lx.advance()
	}
}

func (lx *Lexer) emitEOF(tok *token.Token) {
	tok.Pin = lx.pin()
	tok.Kind = token.EOF
}
