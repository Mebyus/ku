package lexer

import (
	"github.com/mebyus/ku/goku/compiler/baselex"
	"github.com/mebyus/ku/goku/compiler/char"
	"github.com/mebyus/ku/goku/vm/tokens"
)

func (lx *Lexer) Lex() tokens.Token {
	if lx.Eof() {
		return lx.emit(tokens.EOF)
	}

	lx.SkipWhitespaceAndComments()
	if lx.Eof() {
		return lx.emit(tokens.EOF)
	}

	return lx.lex()
}

func (lx *Lexer) lex() tokens.Token {
	if char.IsLatinLetterOrUnderscore(lx.Peek()) {
		return lx.word()
	}

	if char.IsDecDigit(lx.Peek()) {
		return lx.number()
	}

	if lx.Peek() == '#' {
		return lx.keyword()
	}

	if lx.Peek() == '@' && lx.Next() == '.' {
		return lx.label()
	}

	return lx.other()
}

func (lx *Lexer) number() tokens.Token {
	return tokens.FromBaseToken(lx.Number())
}

func (lx *Lexer) word() (tok tokens.Token) {
	tok.Pin = lx.Pin()

	lx.Start()
	lx.SkipWord()
	word, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(baselex.LengthOverflow)
		return tok
	}

	tok.Kind = tokens.Word
	tok.Data = word
	return tok
}

func (lx *Lexer) keyword() (tok tokens.Token) {
	tok.Pin = lx.Pin()

	lx.Advance() // skip '#'

	lx.Start()
	lx.SkipWord()
	data, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(baselex.LengthOverflow)
		return tok
	}
	if data == "" {
		tok.SetIllegalError(baselex.UnknownDirective)
		tok.Data = "#"
		return tok
	}

	switch data {
	case "fun":
		tok.Kind = tokens.Fun
	case "entry":
		tok.Kind = tokens.Entry
	case "data":
		tok.Kind = tokens.Data
	default:
		tok.SetIllegalError(baselex.UnknownDirective)
		tok.Data = "#" + data
	}

	tok.Flags = tokens.FlagKeyword
	return tok
}

func (lx *Lexer) label() (tok tokens.Token) {
	tok.Pin = lx.Pin()

	lx.Advance() // skip '@'
	lx.Advance() // skip '.'

	lx.Start()
	lx.SkipWord()
	data, ok := lx.Take()
	if !ok {
		tok.SetIllegalError(baselex.LengthOverflow)
		return tok
	}

	tok.Kind = tokens.Label
	tok.Data = data
	return tok
}

func (lx *Lexer) other() tokens.Token {
	switch lx.Peek() {
	case '{':
		return lx.oneByteToken(tokens.LeftCurly)
	case '}':
		return lx.oneByteToken(tokens.RightCurly)
	case '[':
		return lx.oneByteToken(tokens.LeftSquare)
	case ']':
		return lx.oneByteToken(tokens.RightSquare)
	case ';':
		return lx.oneByteToken(tokens.Semicolon)
	case '.':
		return lx.oneByteToken(tokens.Period)
	case ',':
		return lx.oneByteToken(tokens.Comma)
	case ':':
		return lx.oneByteToken(tokens.Colon)
	default:
		return lx.illegalByteToken()
	}
}

func (lx *Lexer) oneByteToken(k tokens.Kind) tokens.Token {
	tok := lx.emit(k)
	lx.Advance()
	return tok
}

func (lx *Lexer) illegalByteToken() tokens.Token {
	tok := lx.emit(tokens.Illegal)
	tok.Data = char.ToString(lx.Peek())
	lx.Advance()
	return tok
}
