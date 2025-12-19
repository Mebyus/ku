package baselex

import (
	"github.com/mebyus/ku/goku/compiler/sm"
)

type Lexer struct {
	text []byte

	mask sm.Pin

	// Byte offset into text.
	pos uint32

	// Mark byte offset into text.
	//
	// Mark is used to slice input text for token literals.
	mark uint32
}

func (lx *Lexer) Init(text *sm.Text) {
	lx.text = text.Data
	lx.mask = sm.PinTextMask(text.ID)
}

func (lx *Lexer) Pin() sm.Pin {
	return lx.mask | sm.Pin(lx.pos)
}

func (lx *Lexer) Eof() bool {
	return lx.pos >= uint32(len(lx.text))
}

// Returns byte at current lexer position.
func (lx *Lexer) Peek() byte {
	return lx.text[lx.pos]
}

// Returns byte after current lexer position.
// Returns 0 if next lexer position is outside of text.
func (lx *Lexer) Next() byte {
	p := lx.pos + 1
	if p >= uint32(len(lx.text)) {
		return 0
	}
	return lx.text[p]
}
