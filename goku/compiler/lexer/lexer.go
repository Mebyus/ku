package lexer

import (
	"github.com/mebyus/ku/goku/compiler/baselex"
	"github.com/mebyus/ku/goku/compiler/sm"
	"github.com/mebyus/ku/goku/compiler/token"
)

type Lexer struct {
	baselex.Lexer
}

func FromText(text *sm.Text) *Lexer {
	lx := Lexer{}
	lx.Init(text)
	return &lx
}

// FromBytes creates Lexer from raw text bytes. This function should be used
// for tests, since it does not set text id needed for pins.
func FromBytes(data []byte) *Lexer {
	return FromText(sm.NewText("", data))
}

// Create token (without literal) of specified kind at current lexer position.
//
// Does not advance lexer scan position.
func (lx *Lexer) emit(k token.Kind) token.Token {
	return token.Token{
		Kind: k,
		Pin:  lx.Pin(),
	}
}
