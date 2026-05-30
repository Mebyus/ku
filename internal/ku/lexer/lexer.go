package lexer

import "github.com/mebyus/ku/internal/ku/sx"

type Lexer struct {
	text []byte

	mask sx.Pin

	// Byte offset into text.
	pos uint32

	// Mark byte offset into text.
	//
	// Mark is used to slice input text for token literals.
	mark uint32
}

func FromText(text *sx.Text) *Lexer {
	lx := Lexer{}
	lx.init(text)
	return &lx
}

func (lx *Lexer) init(text *sx.Text) {
	lx.text = text.Data
	lx.mask = sx.PinTextMask(text.ID)
}

func (lx *Lexer) pin() sx.Pin {
	return lx.mask | sx.Pin(lx.pos)
}

func (lx *Lexer) eof() bool {
	return lx.pos >= uint32(len(lx.text))
}

// Returns byte at current lexer position.
func (lx *Lexer) peek() byte {
	return lx.text[lx.pos]
}

// Returns byte after current lexer position.
// Returns 0 if next lexer position is outside of text.
func (lx *Lexer) next() byte {
	p := lx.pos + 1
	if p >= uint32(len(lx.text)) {
		return 0
	}
	return lx.text[p]
}
