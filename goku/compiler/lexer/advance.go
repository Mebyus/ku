package lexer

import "github.com/mebyus/ku/goku/compiler/char"

// Shift lexer scan position one byte forward.
func (lx *Lexer) advance() {
	lx.pos += 1
}

// Start recording new token data.
func (lx *Lexer) start() {
	lx.mark = lx.pos
}

// Returns a view (slice into scanned text) of recorded token data.
func (lx *Lexer) view() []byte {
	return lx.text[lx.mark:lx.pos]
}

// Returns recorded token data string (lit, true).
//
// Returns ("", false) in case token length overflowed max length.
func (lx *Lexer) take() (string, bool) {
	if lx.isLengthOverflow() {
		return "", false
	}
	return string(lx.view()), true
}

// Returns byte length of recorded token data.
func (lx *Lexer) length() uint32 {
	return lx.pos - lx.mark
}

const maxTokenByteLength = 1 << 12

func (lx *Lexer) isLengthOverflow() bool {
	return lx.length() > maxTokenByteLength
}

func (lx *Lexer) skipWord() {
	for !lx.eof() && char.IsAlphanum(lx.c()) {
		lx.advance()
	}
}

func (lx *Lexer) skipBinDigits() {
	for !lx.eof() && char.IsBinDigit(lx.c()) {
		lx.advance()
	}
}

func (lx *Lexer) skipOctDigits() {
	for !lx.eof() && char.IsOctDigit(lx.c()) {
		lx.advance()
	}
}

func (lx *Lexer) skipHexDigits() {
	for !lx.eof() && char.IsHexDigit(lx.c()) {
		lx.advance()
	}
}

func (lx *Lexer) skipWhitespaceAndComments() {
	for {
		lx.skipWhitespace()
		if lx.eof() {
			return
		}

		if lx.c() == '/' && lx.n() == '/' {
			lx.skipLineComment()
		} else if lx.c() == '/' && lx.n() == '*' {
			lx.skipBlockComment()
		} else {
			return
		}
	}
}

func (lx *Lexer) skipWhitespace() {
	for !lx.eof() && char.IsSimpleWhitespace(lx.c()) {
		lx.advance()
	}
}

func (lx *Lexer) skipLine() {
	for !lx.eof() && lx.c() != '\n' {
		lx.advance()
	}
	if lx.eof() {
		return
	}

	lx.advance() // skip newline byte
}

func (lx *Lexer) skipLineComment() {
	lx.advance() // skip '/'
	lx.advance() // skip '/'
	lx.skipLine()
}

func (lx *Lexer) skipBlockComment() {
	lx.advance() // skip '/'
	lx.advance() // skip '*'

	for !lx.eof() && !(lx.c() == '*' && lx.n() == '/') {
		lx.advance()
	}

	if lx.eof() {
		return
	}

	lx.advance() // skip '*'
	lx.advance() // skip '/'
}
