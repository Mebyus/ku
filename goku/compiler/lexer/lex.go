package lexer

import (
	"github.com/mebyus/ku/goku/compiler/char"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (lx *Lexer) Lex() token.Token {
	if lx.Eof() {
		return lx.emit(token.EOF)
	}

	lx.SkipWhitespaceAndComments()
	if lx.Eof() {
		return lx.emit(token.EOF)
	}

	return lx.lex()
}

func (lx *Lexer) lex() token.Token {
	if char.IsLatinLetterOrUnderscore(lx.Peek()) {
		return lx.word()
	}

	if char.IsDecDigit(lx.Peek()) {
		return lx.number()
	}

	if lx.Peek() == '"' {
		return lx.str()
	}

	if lx.Peek() == '\'' {
		return lx.rune()
	}

	if lx.Peek() == '#' {
		switch lx.Next() {
		case '{':
			return lx.twoBytesToken(token.HashCurly)
		case '[':
			return lx.twoBytesToken(token.HashSquare)
		case ':':
			return lx.env()
		default:
			if char.IsLatinLetter(lx.Next()) {
				return lx.static()
			}
			return lx.illegalByteToken()
		}
	}

	if lx.Peek() == '@' && lx.Next() == '.' {
		return lx.label()
	}

	return lx.other()
}

func (lx *Lexer) rune() (tok token.Token) {
	tok.Pin = lx.Pin()

	lx.Advance() // skip "'"
	if lx.Eof() {
		tok.SetIllegalError(token.MalformedRune)
		tok.Data = "'"
		return tok
	}

	lx.Start()
	if lx.Peek() == '\\' {
		// handle escape sequence
		var val uint64
		switch lx.Next() {
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
			lx.Advance() // skip "\"
			lx.Advance() // skip unknown escape rune

			if !lx.Eof() && lx.Peek() == '\'' {
				lx.Advance()
			}
			tok.SetIllegalError(token.MalformedRune)
			tok.Data, _ = lx.Take()
			return tok
		}

		lx.Advance() // skip "\"
		lx.Advance() // skip escape rune
		if lx.Eof() {
			tok.SetIllegalError(token.MalformedRune)
			tok.Data, _ = lx.Take()
			return tok
		}
		if lx.Peek() != '\'' {
			lx.Advance()
			tok.SetIllegalError(token.MalformedRune)
			tok.Data, _ = lx.Take()
			return tok
		}

		lx.Advance() // skip "'"
		tok.Kind = token.Rune
		tok.Val = val
		return tok
	}

	if lx.Next() == '\'' {
		// common case of ascii rune
		tok.Val = uint64(lx.Peek())
		tok.Kind = token.Rune
		lx.Advance()
		lx.Advance()
		return tok
	}

	// handle non-ascii runes
	for !lx.Eof() && lx.Peek() != '\'' && lx.Peek() != '\n' {
		lx.Advance()
	}

	data, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	if lx.Eof() || lx.Peek() != '\'' {
		tok.SetIllegalError(token.MalformedRune)
		tok.Data = data
		return tok
	}
	lx.Advance() // skip "'"

	runes := []rune(data)
	if len(runes) != 1 {
		tok.SetIllegalError(token.MalformedRune)
		tok.Data = data
		return tok
	}

	tok.Kind = token.Rune
	tok.Val = uint64(runes[0])
	return tok
}

func (lx *Lexer) binNumber() (tok token.Token) {
	tok.Pin = lx.Pin()

	lx.Advance() // skip '0'
	lx.Advance() // skip 'b'

	lx.Start()
	lx.SkipBinDigits()

	if char.IsAlphanum(lx.Peek()) {
		lx.SkipWord()
		data, ok := lx.Take()
		if ok {
			tok.SetIllegalError(token.MalformedBinaryInteger)
			tok.Data = data
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return tok
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}
	if lx.Length() == 0 {
		tok.SetIllegalError(token.MalformedBinaryInteger)
		tok.Data = "0b"
		return tok
	}

	tok.Kind = token.BinInteger
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

func (lx *Lexer) octNumber() token.Token {
	var tok token.Token
	tok.Pin = lx.Pin()

	lx.Advance() // skip '0' byte
	lx.Advance() // skip 'o' byte

	lx.Start()
	lx.SkipOctDigits()

	if char.IsAlphanum(lx.Peek()) {
		lx.SkipWord()
		data, ok := lx.Take()
		if ok {
			tok.SetIllegalError(token.MalformedOctalInteger)
			tok.Data = data
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return tok
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}
	if lx.Length() == 0 {
		tok.SetIllegalError(token.MalformedOctalInteger)
		tok.Data = "0o"
		return tok
	}

	tok.Kind = token.OctInteger
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

func (lx *Lexer) decNumber() token.Token {
	var tok token.Token
	tok.Pin = lx.Pin()

	lx.Start()
	lx.Advance() // skip first digit
	scannedOnePeriod := false
	for !lx.Eof() && char.IsDecDigitOrPeriod(lx.Peek()) {
		if lx.Peek() == '.' {
			if scannedOnePeriod || !char.IsDecDigit(lx.Next()) {
				data, ok := lx.Take()
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
		lx.Advance()
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	if !lx.Eof() && char.IsAlphanum(lx.Peek()) {
		lx.SkipWord()
		data, ok := lx.Take()
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
		n, ok := char.ParseDecDigitsWithOverflowCheck(lx.View())
		if !ok {
			tok.SetIllegalError(token.DecimalIntegerOverflow)
			return tok
		}

		tok.Kind = token.DecInteger
		tok.Val = n
		return tok
	}

	tok.Kind = token.DecFloat
	tok.Data, _ = lx.Take()
	return tok
}

func (lx *Lexer) hexNumber() token.Token {
	var tok token.Token
	tok.Pin = lx.Pin()

	lx.Advance() // skip "0"
	lx.Advance() // skip "x"

	lx.Start()
	lx.SkipHexDigits()

	if char.IsAlphanum(lx.Peek()) {
		lx.SkipWord()
		data, ok := lx.Take()
		if ok {
			tok.SetIllegalError(token.MalformedHexadecimalInteger)
			tok.Data = data
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return tok
	}

	if lx.IsLengthOverflow() {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}
	if lx.Length() == 0 {
		tok.SetIllegalError(token.MalformedHexadecimalInteger)
		tok.Data = "0x"
		return tok
	}

	tok.Kind = token.HexInteger
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

func (lx *Lexer) number() (tok token.Token) {
	if lx.Peek() != '0' {
		return lx.decNumber()
	}

	if lx.Next() == 'b' {
		return lx.binNumber()
	}

	if lx.Next() == 'o' {
		return lx.octNumber()
	}

	if lx.Next() == 'x' {
		return lx.hexNumber()
	}

	if lx.Next() == '.' {
		return lx.decNumber()
	}

	if char.IsAlphanum(lx.Next()) {
		return lx.illegalWord(token.MalformedDecimalInteger)
	}

	tok = token.Token{
		Kind: token.DecInteger,
		Pin:  lx.Pin(),
		Val:  0,
	}
	lx.Advance()
	return
}

func (lx *Lexer) env() token.Token {
	var tok token.Token
	tok.Pin = lx.Pin()

	lx.Advance() // skip '#'
	lx.Advance() // skip ':'

	if !char.IsLatinLetter(lx.Peek()) {
		tok.SetIllegalError(token.MalformedEnv)
		tok.Data = "#:"
		return tok
	}

	lx.Start()
	for !lx.Eof() && (char.IsAlphanum(lx.Peek()) || lx.Peek() == '.') {
		if lx.Peek() == '.' && lx.Next() == '.' {
			data, ok := lx.Take()
			if !ok {
				tok.SetIllegalError(token.LengthOverflow)
				return tok
			}
			tok.SetIllegalError(token.MalformedEnv)
			tok.Data = data
			return tok
		}
		lx.Advance()
	}

	data, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	tok.Kind = token.Env
	tok.Data = data

	return tok
}

func (lx *Lexer) static() token.Token {
	var tok token.Token
	tok.Pin = lx.Pin()

	lx.Advance() // skip '#'

	lx.Start()
	lx.SkipWord()
	data, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	switch data {
	case "must":
		tok.Kind = token.StaticMust
	case "if":
		tok.Kind = token.StaticIf
	case "typeid":
		tok.Kind = token.TypeId
	case "error":
		tok.Kind = token.ErrorId
	case "enum":
		tok.Kind = token.Enum
	case "size":
		tok.Kind = token.Size
	case "check":
		tok.Kind = token.Check
	case "len":
		tok.Kind = token.Len
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
		tok.Data = "#" + data
	}

	return tok
}

func (lx *Lexer) word() token.Token {
	var tok token.Token
	tok.Pin = lx.Pin()

	if !char.IsAlphanum(lx.Next()) {
		// word is 1 character long
		c := lx.Peek()
		lx.Advance() // skip single (start) character

		if c == '_' {
			tok.Kind = token.Underscore
		} else {
			tok.Kind = token.Word
			tok.Data = char.ToString(c)
		}
		return tok
	}

	// word is at least 2 characters long
	lx.Start()
	lx.SkipWord()
	word, ok := lx.Take()
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
	tok.Pin = lx.Pin()

	lx.Advance() // skip quote
	if lx.Eof() {
		tok.SetIllegalError(token.MalformedString)
		tok.Data = "\""
		return tok
	}

	if lx.Peek() == '"' {
		// common case of empty string literal
		lx.Advance() // skip quote
		tok.Kind = token.String
		return tok
	}

	var fills uint64 // number of fill places inside the string
	lx.Start()
	for !lx.Eof() && lx.Peek() != '"' && lx.Peek() != '\n' {
		if lx.Peek() == '\\' && lx.Next() == '"' {
			// do not stop if we encounter escape sequence
			lx.Advance() // skip "\"
			lx.Advance() // skip quote
		} else if lx.Peek() == '$' && lx.Next() == '{' {
			fills += 1
			lx.Advance() // skip "$"
			lx.Advance() // skip "{"
		} else {
			lx.Advance()
		}
	}

	if lx.Eof() {
		data, ok := lx.Take()
		if ok {
			tok.SetIllegalError(token.MalformedString)
			tok.Data = data
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return tok
	}

	if lx.Peek() != '"' {
		data, ok := lx.Take()
		if ok {
			tok.SetIllegalError(token.MalformedString)
			tok.Data = data
		} else {
			tok.SetIllegalError(token.LengthOverflow)
		}
		return tok
	}

	data, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	lx.Advance() // skip quote

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
	tok.Pin = lx.Pin()

	lx.Advance() // skip '@'
	lx.Advance() // skip '.'

	lx.Start()
	lx.SkipWord()
	data, ok := lx.Take()
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
	switch lx.Peek() {
	case '(':
		return lx.oneByteToken(token.LeftParen)
	case ')':
		return lx.oneByteToken(token.RightParen)
	case '{':
		return lx.oneByteToken(token.LeftCurly)
	case '}':
		return lx.oneByteToken(token.RightCurly)
	case '[':
		if lx.Next() == ']' {
			return lx.twoBytesToken(token.Chunk)
		}
		if lx.Next() == '_' {
			pin := lx.Pin()
			lx.Advance() // skip "["
			if lx.Next() != ']' {
				lx.Advance() // skip "_"
				return token.Token{
					Pin:  pin,
					Kind: token.Illegal,
					Data: "[_",
				}
			}
			lx.Advance() // skip "_"
			lx.Advance() // skip "]"
			return token.Token{
				Pin:  pin,
				Kind: token.AutoLen,
			}
		}
		if lx.Next() == '*' {
			pin := lx.Pin()
			lx.Advance() // skip "["
			if lx.Next() != ']' {
				return token.Token{
					Pin:  pin,
					Kind: token.LeftSquare,
				}
			}
			lx.Advance() // skip "*"
			lx.Advance() // skip "]"
			return token.Token{
				Pin:  pin,
				Kind: token.ArrayPointer,
			}
		}
		if lx.Next() == '&' {
			pin := lx.Pin()
			lx.Advance() // skip "["
			if lx.Next() != ']' {
				return token.Token{
					Pin:  pin,
					Kind: token.LeftSquare,
				}
			}
			lx.Advance() // skip "&"
			lx.Advance() // skip "]"
			return token.Token{
				Pin:  pin,
				Kind: token.ArrayRef,
			}
		}
		if lx.Next() == '^' {
			pin := lx.Pin()
			lx.Advance() // skip "["
			if lx.Next() != ']' {
				return token.Token{
					Pin:  pin,
					Kind: token.LeftSquare,
				}
			}
			lx.Advance() // skip "^"
			lx.Advance() // skip "]"
			return token.Token{
				Pin:  pin,
				Kind: token.CapBuffer,
			}
		}
		return lx.oneByteToken(token.LeftSquare)
	case ']':
		return lx.oneByteToken(token.RightSquare)
	case '<':
		switch lx.Next() {
		case '=':
			return lx.twoBytesToken(token.LessOrEqual)
		case '<':
			pin := lx.Pin()
			lx.Advance() // skip "<"
			lx.Advance() // skip "<"
			if lx.Peek() == '=' {
				lx.Advance() // skip "="
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
		switch lx.Next() {
		case '=':
			return lx.twoBytesToken(token.GreaterOrEqual)
		case '>':
			pin := lx.Pin()
			lx.Advance() // skip ">"
			lx.Advance() // skip ">"
			if lx.Peek() == '=' {
				lx.Advance() // skip "="
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
		if lx.Next() == '=' {
			return lx.twoBytesToken(token.AddAssign)
		}
		return lx.oneByteToken(token.Plus)
	case '-':
		if lx.Next() == '=' {
			return lx.twoBytesToken(token.SubAssign)
		}
		return lx.oneByteToken(token.Minus)
	case ',':
		return lx.oneByteToken(token.Comma)
	case '=':
		switch lx.Next() {
		case '=':
			return lx.twoBytesToken(token.Equal)
		case '>':
			return lx.twoBytesToken(token.RightArrow)
		default:
			return lx.oneByteToken(token.Assign)
		}
	case ':':
		if lx.Next() == '=' {
			return lx.twoBytesToken(token.Walrus)
		}
		return lx.oneByteToken(token.Colon)
	case ';':
		return lx.oneByteToken(token.Semicolon)
	case '.':
		switch lx.Next() {
		case '.':
			pin := lx.Pin()
			lx.Advance() // skip "."
			lx.Advance() // skip "."
			if lx.Peek() != '.' {
				return token.Token{
					Pin:  pin,
					Kind: token.Illegal,
				}
			}
			lx.Advance() // skip "."
			return token.Token{
				Pin:  pin,
				Kind: token.Ellipsis,
			}
		case '&':
			return lx.twoBytesToken(token.Address)
		case '*':
			pin := lx.Pin()
			lx.Advance()
			lx.Advance()
			if lx.Eof() || lx.Peek() != '.' {
				return token.Token{
					Pin:  pin,
					Kind: token.Deref,
				}
			}

			lx.Advance()
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
		if lx.Next() == '&' {
			return lx.twoBytesToken(token.And)
		}
		if lx.Next() == '=' {
			return lx.twoBytesToken(token.AndAssign)
		}
		return lx.oneByteToken(token.Ampersand)
	case '/':
		if lx.Next() == '=' {
			return lx.twoBytesToken(token.DivAssign)
		}
		return lx.oneByteToken(token.Slash)
	case '!':
		if lx.Next() == '=' {
			return lx.twoBytesToken(token.NotEqual)
		}
		return lx.oneByteToken(token.Not)
	case '?':
		return lx.oneByteToken(token.Quest)
	case '^':
		return lx.oneByteToken(token.Caret)
	case '|':
		if lx.Next() == '|' {
			return lx.twoBytesToken(token.Or)
		}
		if lx.Next() == '=' {
			return lx.twoBytesToken(token.OrAssign)
		}
		return lx.oneByteToken(token.Pipe)
	default:
		return lx.illegalByteToken()
	}
}

func (lx *Lexer) oneByteToken(k token.Kind) token.Token {
	tok := lx.emit(k)
	lx.Advance()
	return tok
}

func (lx *Lexer) twoBytesToken(k token.Kind) token.Token {
	tok := lx.emit(k)
	lx.Advance()
	lx.Advance()
	return tok
}

func (lx *Lexer) illegalByteToken() token.Token {
	tok := lx.emit(token.Illegal)
	tok.Data = char.ToString(byte(lx.Peek()))
	lx.Advance()
	return tok
}

func (lx *Lexer) illegalWord(code uint64) token.Token {
	var tok token.Token
	tok.Pin = lx.Pin()

	lx.Start()
	lx.SkipWord()
	data, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(token.LengthOverflow)
		return tok
	}

	tok.SetIllegalError(code)
	tok.Data = data
	return tok
}
