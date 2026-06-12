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
		lx.emitOneByteToken(tok, token.Assign)
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
	for !char.IsTextByte(lx.peek()) {
		lx.advance()
	}
}

func (lx *Lexer) emitEOF(tok *token.Token) {
	tok.Pin = lx.pin()
	tok.Kind = token.EOF
}
