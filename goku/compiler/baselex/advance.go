package baselex

import "github.com/mebyus/ku/goku/compiler/char"

// Shift lexer scan position one byte forward.
func (lx *Lexer) Advance() {
	lx.pos += 1
}

// Start recording new token data.
func (lx *Lexer) Start() {
	lx.mark = lx.pos
}

// Returns a view (slice into scanned text) of recorded token data.
func (lx *Lexer) View() []byte {
	return lx.text[lx.mark:lx.pos]
}

// Returns recorded token data string (lit, true).
//
// Returns ("", false) in case token length overflowed max length.
func (lx *Lexer) Take() (string, bool) {
	if lx.IsLengthOverflow() {
		return "", false
	}
	return string(lx.View()), true
}

// Returns byte length of recorded token data.
func (lx *Lexer) Length() uint32 {
	return lx.pos - lx.mark
}

const maxTokenByteLength = 1 << 12

func (lx *Lexer) IsLengthOverflow() bool {
	return lx.Length() > maxTokenByteLength
}

func (lx *Lexer) SkipWord() {
	for !lx.Eof() && char.IsAlphanum(lx.Peek()) {
		lx.Advance()
	}
}

func (lx *Lexer) SkipBinDigits() {
	for !lx.Eof() && char.IsBinDigit(lx.Peek()) {
		lx.Advance()
	}
}

func (lx *Lexer) SkipOctDigits() {
	for !lx.Eof() && char.IsOctDigit(lx.Peek()) {
		lx.Advance()
	}
}

func (lx *Lexer) SkipHexDigits() {
	for !lx.Eof() && char.IsHexDigit(lx.Peek()) {
		lx.Advance()
	}
}

func (lx *Lexer) SkipWhitespaceAndComments() {
	for {
		lx.SkipWhitespace()
		if lx.Eof() {
			return
		}

		if lx.Peek() == '/' && lx.Next() == '/' {
			lx.SkipLineComment()
		} else if lx.Peek() == '/' && lx.Next() == '*' {
			lx.SkipBlockComment()
		} else {
			return
		}
	}
}

func (lx *Lexer) SkipWhitespace() {
	for !lx.Eof() && char.IsSimpleWhitespace(lx.Peek()) {
		lx.Advance()
	}
}

func (lx *Lexer) SkipLine() {
	for !lx.Eof() && lx.Peek() != '\n' {
		lx.Advance()
	}
	if lx.Eof() {
		return
	}

	lx.Advance() // skip newline byte
}

func (lx *Lexer) SkipLineComment() {
	lx.Advance() // skip '/'
	lx.Advance() // skip '/'
	lx.SkipLine()
}

func (lx *Lexer) SkipBlockComment() {
	lx.Advance() // skip '/'
	lx.Advance() // skip '*'

	for !lx.Eof() && !(lx.Peek() == '*' && lx.Next() == '/') {
		lx.Advance()
	}

	if lx.Eof() {
		return
	}

	lx.Advance() // skip '*'
	lx.Advance() // skip '/'
}
