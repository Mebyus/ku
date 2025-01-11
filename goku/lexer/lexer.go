package lexer

import (
	"github.com/mebyus/ku/goku/source"
)

type Lexer struct {
	text []byte

	mask source.Pin

	// Byte offset into text.
	pos uint32
}

func FromText(text *source.Text) *Lexer {
	return &Lexer{
		text: text.Data,
		mask: source.PinTextMask(text.ID),
	}
}

func (lx *Lexer) pin() source.Pin {
	return lx.mask | source.Pin(lx.pos)
}

func (lx *Lexer) eof() bool {
	return lx.pos >= uint32(len(lx.text))
}

// Returns byte at current lexer position.
func (lx *Lexer) c() byte {
	return lx.text[lx.pos]
}

// Returns byte after current lexer position.
// Returns 0 if next lexer position is outside of text.
func (lx *Lexer) n() byte {
	p := lx.pos + 1
	if p >= uint32(len(lx.text)) {
		return 0
	}
	return lx.text[p]
}
