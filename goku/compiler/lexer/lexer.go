package lexer

import (
	"github.com/mebyus/ku/goku/compiler/source"
	"github.com/mebyus/ku/goku/compiler/token"
)

type Lexer struct {
	text []byte

	mask source.Pin

	// Byte offset into text.
	pos uint32

	// Mark byte offset into text.
	//
	// Mark is used to slice input text for token literals.
	mark uint32
}

func FromText(text *source.Text) *Lexer {
	return &Lexer{
		text: text.Data,
		mask: source.PinTextMask(text.ID),
	}
}

// FromBytes creates Lexer from raw text bytes. This function should be used
// for tests, since it does not set text id needed for pins.
func FromBytes(data []byte) *Lexer {
	return &Lexer{text: data}
}

func (lx *Lexer) pin() source.Pin {
	return lx.mask | source.Pin(lx.pos)
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
