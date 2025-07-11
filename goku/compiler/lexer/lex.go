package lexer

import (
	"github.com/mebyus/ku/goku/compiler/char"
	"github.com/mebyus/ku/goku/compiler/token"
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
	if char.IsLatinLetterOrUnderscore(lx.c()) {
		return lx.word()
	}

	if char.IsDecDigit(lx.c()) {
		return lx.number()
	}

	if lx.c() == '"' {
		return lx.str()
	}

	if lx.c() == '\'' {
		return lx.rune()
	}

	if lx.c() == '#' {
		switch lx.n() {
		case '{':
			return lx.twoBytesToken(token.HashCurly)
		case '[':
			return lx.twoBytesToken(token.HashSquare)
		default:
			if char.IsLatinLetter(lx.n()) {
				return lx.static()
			}
			return lx.illegalByteToken()
		}
	}

	if lx.c() == '@' && lx.n() == '.' {
		return lx.label()
	}

	return lx.other()
}

func (lx *Lexer) rune() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.advance() // skip "'"
	if lx.eof() {
		tok.SetIllegalError(token.MalformedRune)
		tok.Data = "'"
		return tok
	}

	lx.start()
	if lx.c() == '\\' {
		// handle escape sequence
		var val uint64
		switch lx.n() {
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

			if !lx.eof() && lx.c() == '\'' {
				lx.advance()
			}
			tok.SetIllegalError(token.MalformedRune)
			tok.Data, _ = lx.take()
			return tok
		}

		lx.advance() // skip "\"
		lx.advance() // skip escape rune
		if lx.eof() {
			tok.SetIllegalError(token.MalformedRune)
			tok.Data, _ = lx.take()
			return tok
		}
		if lx.c() != '\'' {
			lx.advance()
			tok.SetIllegalError(token.MalformedRune)
			tok.Data, _ = lx.take()
			return tok
		}

		lx.advance() // skip "'"
		tok.Kind = token.Rune
		tok.Val = val
		return tok
	}

	if lx.n() == '\'' {
		// common case of ascii rune
		tok.Val = uint64(lx.c())
		tok.Kind = token.Rune
		lx.advance()
		lx.advance()
		return tok
	}

	// handle non-ascii runes
	for !lx.eof() && lx.c() != '\'' && lx.c() != '\n' {
		lx.advance()
	}

	data, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	if lx.eof() || lx.c() != '\'' {
		tok.SetIllegalError(token.MalformedRune)
		tok.Data = data
		return tok
	}
	lx.advance() // skip "'"

	tok.Kind = token.Rune
	tok.Data = data // TODO: parse rune val
	return tok
}

func (lx *Lexer) binNumber() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.advance() // skip '0'
	lx.advance() // skip 'b'

	lx.start()
	lx.skipBinDigits()

	if char.IsAlphanum(lx.c()) {
		lx.skipWord()
		data, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedBinaryInteger)
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
		tok.SetIllegalError(token.MalformedBinaryInteger)
		tok.Data = "0b"
		return tok
	}

	tok.Kind = token.BinInteger
	if lx.length() > 64 {
		lit, ok := lx.take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Data = lit
		return tok
	}

	tok.Val = char.ParseBinDigits(lx.view())
	return tok
}

func (lx *Lexer) octNumber() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.advance() // skip '0' byte
	lx.advance() // skip 'o' byte

	lx.start()
	lx.skipOctDigits()

	if char.IsAlphanum(lx.c()) {
		lx.skipWord()
		data, ok := lx.take()
		if ok {
			tok.SetIllegalError(token.MalformedOctalInteger)
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
		tok.SetIllegalError(token.MalformedOctalInteger)
		tok.Data = "0o"
		return tok
	}

	tok.Kind = token.OctInteger
	if lx.length() > 21 {
		data, ok := lx.take()
		if !ok {
			panic("unreachable due to previous checks")
		}
		tok.Data = data
		return tok
	}

	tok.Val = char.ParseOctDigits(lx.view())
	return tok
}

func (lx *Lexer) decNumber() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.start()
	lx.advance() // skip first digit
	scannedOnePeriod := false
	for !lx.eof() && char.IsDecDigitOrPeriod(lx.c()) {
		if lx.c() == '.' {
			if scannedOnePeriod || !char.IsDecDigit(lx.n()) {
				data, ok := lx.take()
				if ok {
					tok.SetIllegalError(token.MalformedDecimalFloat)
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

	if !lx.eof() && char.IsAlphanum(lx.c()) {
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

	tok.Kind = token.DecFloat
	tok.Data, _ = lx.take()
	return tok
}

func (lx *Lexer) hexNumber() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.advance() // skip "0"
	lx.advance() // skip "x"

	lx.start()
	lx.skipHexDigits()

	if char.IsAlphanum(lx.c()) {
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
	if lx.c() != '0' {
		return lx.decNumber()
	}

	if lx.n() == 'b' {
		return lx.binNumber()
	}

	if lx.n() == 'o' {
		return lx.octNumber()
	}

	if lx.n() == 'x' {
		return lx.hexNumber()
	}

	if lx.n() == '.' {
		return lx.decNumber()
	}

	if char.IsAlphanum(lx.n()) {
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

func (lx *Lexer) static() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.advance() // skip '#'

	lx.start()
	lx.skipWord()
	data, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	switch data {
	case "must":
		tok.Kind = token.StaticMust
	case "typeid":
		tok.Kind = token.TypeId
	case "error":
		tok.Kind = token.ErrorId
	case "enum":
		tok.Kind = token.Enum
	case "size":
		tok.Kind = token.Size
	case "never":
		tok.Kind = token.Never
	case "stub":
		tok.Kind = token.Stub
	case "build":
		tok.Kind = token.Build
	case "debug":
		tok.Kind = token.Debug
	case "lookup":
		tok.Kind = token.Lookup
	default:
		tok.SetIllegalError(token.UnknownDirective)
	}

	return tok
}

func (lx *Lexer) word() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	if !char.IsAlphanum(lx.n()) {
		// word is 1 character long
		c := lx.c()
		lx.advance() // skip single (start) character

		if c == '_' {
			tok.Kind = token.Underscore
		} else {
			tok.Kind = token.Word
			tok.Data = char.ToString(c)
		}
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
		tok.Flags = token.FlagKeyword
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

	if lx.c() == '"' {
		// common case of empty string literal
		lx.advance() // skip quote
		tok.Kind = token.String
		return tok
	}

	var fills uint64 // number of fill places inside the string
	lx.start()
	for !lx.eof() && lx.c() != '"' && lx.c() != '\n' {
		if lx.c() == '\\' && lx.n() == '"' {
			// do not stop if we encounter escape sequence
			lx.advance() // skip "\"
			lx.advance() // skip quote
		} else if lx.c() == '$' && lx.n() == '{' {
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

	if lx.c() != '"' {
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
		tok.Data = data
		tok.Kind = token.FillString
		return tok
	}

	data, ok = char.Unescape(data)
	if !ok {
		tok.SetIllegalError(token.BadEscapeInString)
		return tok
	}

	tok.Kind = token.String
	tok.Data = data
	return tok
}

func (lx *Lexer) label() token.Token {
	var tok token.Token
	tok.Pin = lx.pin()

	lx.advance() // skip '@'
	lx.advance() // skip '.'

	lx.start()
	lx.skipWord()
	data, ok := lx.take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	switch data {
	case "next":
		tok.Kind = token.LabelNext
	case "out":
		tok.Kind = token.LabelOut
	default:
		tok.Data = data
		panic("arbitrary labels not implemented: " + data)
	}
	return tok
}

func (lx *Lexer) other() token.Token {
	switch lx.c() {
	case '(':
		return lx.oneByteToken(token.LeftParen)
	case ')':
		return lx.oneByteToken(token.RightParen)
	case '{':
		return lx.oneByteToken(token.LeftCurly)
	case '}':
		return lx.oneByteToken(token.RightCurly)
	case '[':
		if lx.n() == ']' {
			return lx.twoBytesToken(token.Chunk)
		}
		if lx.n() == '_' {
			pin := lx.pin()
			lx.advance() // skip "["
			if lx.n() != ']' {
				lx.advance() // skip "_"
				return token.Token{
					Pin:  pin,
					Kind: token.Illegal,
					Data: "[_",
				}
			}
			lx.advance() // skip "_"
			lx.advance() // skip "]"
			return token.Token{
				Pin:  pin,
				Kind: token.AutoLen,
			}
		}
		if lx.n() == '*' {
			pin := lx.pin()
			lx.advance() // skip "["
			if lx.n() != ']' {
				return token.Token{
					Pin:  pin,
					Kind: token.LeftSquare,
				}
			}
			lx.advance() // skip "*"
			lx.advance() // skip "]"
			return token.Token{
				Pin:  pin,
				Kind: token.ArrayPointer,
			}
		}
		if lx.n() == '^' {
			pin := lx.pin()
			lx.advance() // skip "["
			if lx.n() != ']' {
				return token.Token{
					Pin:  pin,
					Kind: token.LeftSquare,
				}
			}
			lx.advance() // skip "^"
			lx.advance() // skip "]"
			return token.Token{
				Pin:  pin,
				Kind: token.CapBuffer,
			}
		}
		return lx.oneByteToken(token.LeftSquare)
	case ']':
		return lx.oneByteToken(token.RightSquare)
	case '<':
		switch lx.n() {
		case '=':
			return lx.twoBytesToken(token.LessOrEqual)
		case '<':
			pin := lx.pin()
			lx.advance() // skip "<"
			lx.advance() // skip "<"
			if lx.c() == '=' {
				lx.advance() // skip "="
				return token.Token{
					Pin:  pin,
					Kind: token.LeftShiftAssign,
				}
			}
			return token.Token{
				Pin:  pin,
				Kind: token.LeftShift,
			}
		case '-':
			return lx.twoBytesToken(token.LeftArrow)
		default:
			return lx.oneByteToken(token.LeftAngle)
		}
	case '>':
		switch lx.n() {
		case '=':
			return lx.twoBytesToken(token.GreaterOrEqual)
		case '>':
			pin := lx.pin()
			lx.advance() // skip ">"
			lx.advance() // skip ">"
			if lx.c() == '=' {
				lx.advance() // skip "="
				return token.Token{
					Pin:  pin,
					Kind: token.RightShiftAssign,
				}
			}
			return token.Token{
				Pin:  pin,
				Kind: token.RightShift,
			}
		default:
			return lx.oneByteToken(token.RightAngle)
		}
	case '+':
		if lx.n() == '=' {
			return lx.twoBytesToken(token.AddAssign)
		}
		return lx.oneByteToken(token.Plus)
	case '-':
		if lx.n() == '=' {
			return lx.twoBytesToken(token.SubAssign)
		}
		return lx.oneByteToken(token.Minus)
	case ',':
		return lx.oneByteToken(token.Comma)
	case '=':
		switch lx.n() {
		case '=':
			return lx.twoBytesToken(token.Equal)
		case '>':
			return lx.twoBytesToken(token.RightArrow)
		default:
			return lx.oneByteToken(token.Assign)
		}
	case ':':
		if lx.n() == '=' {
			return lx.twoBytesToken(token.Walrus)
		}
		return lx.oneByteToken(token.Colon)
	case ';':
		return lx.oneByteToken(token.Semicolon)
	case '.':
		switch lx.n() {
		case '.':
			pin := lx.pin()
			lx.advance() // skip "."
			lx.advance() // skip "."
			if lx.c() != '.' {
				return token.Token{
					Pin:  pin,
					Kind: token.Illegal,
				}
			}
			lx.advance() // skip "."
			return token.Token{
				Pin:  pin,
				Kind: token.Ellipsis,
			}
		case '&':
			return lx.twoBytesToken(token.Address)
		case '*':
			pin := lx.pin()
			lx.advance()
			lx.advance()
			if lx.eof() || lx.c() != '.' {
				return token.Token{
					Pin:  pin,
					Kind: token.Deref,
				}
			}

			lx.advance()
			return token.Token{
				Pin:  pin,
				Kind: token.DerefSelect,
			}
		case '{':
			return lx.twoBytesToken(token.Tweak)
		case '[':
			return lx.twoBytesToken(token.DerefIndex)
		case '(':
			return lx.twoBytesToken(token.BagSelect)
		case '!':
			return lx.twoBytesToken(token.Insist)
		case '?':
			return lx.twoBytesToken(token.Chain)
		default:
			return lx.oneByteToken(token.Period)
		}
	case '%':
		return lx.oneByteToken(token.Percent)
	case '*':
		return lx.oneByteToken(token.Asterisk)
	case '&':
		if lx.n() == '&' {
			return lx.twoBytesToken(token.And)
		}
		return lx.oneByteToken(token.Ampersand)
	case '/':
		if lx.n() == '=' {
			return lx.twoBytesToken(token.DivAssign)
		}
		return lx.oneByteToken(token.Slash)
	case '!':
		if lx.n() == '=' {
			return lx.twoBytesToken(token.NotEqual)
		}
		return lx.oneByteToken(token.Not)
	case '?':
		return lx.oneByteToken(token.Quest)
	case '^':
		return lx.oneByteToken(token.Caret)
	case '|':
		if lx.n() == '|' {
			return lx.twoBytesToken(token.Or)
		}
		return lx.oneByteToken(token.Pipe)
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
	tok.Data = char.ToString(byte(lx.c()))
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
