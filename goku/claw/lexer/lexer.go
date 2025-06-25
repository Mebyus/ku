package lexer

import (
	"github.com/mebyus/ku/goku/claw/token"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type Lexer struct {
	text []byte

	mask srcmap.Pin

	// Byte offset into text.
	pos uint32

	// Mark byte offset into text.
	//
	// Mark is used to slice input text for token literals.
	mark uint32
}

func FromText(text *srcmap.Text) *Lexer {
	return &Lexer{
		text: text.Data,
		mask: srcmap.PinTextMask(text.ID),
	}
}

// FromBytes creates Lexer from raw text bytes. This function should be used
// for tests, since it does not set text id needed for pins.
func FromBytes(data []byte) *Lexer {
	return &Lexer{text: data}
}

func (lx *Lexer) pin() srcmap.Pin {
	return lx.mask | srcmap.Pin(lx.pos)
}

// Create token (without literal) of specified kind at current lexer position.
//
// Does not advance lexer scan position.
func (lx *Lexer) emit(k token.Kind) token.Token {
	return token.Token{
		Kind: k,
		Pin:  lx.pin(),
	}
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
