package lexer

import (
	"io"

	"github.com/mebyus/ku/goku/compiler/source"
	"github.com/mebyus/ku/goku/compiler/token"
)

type Stream interface {
	Lex() token.Token
}

// Flexer same as Stream, but can yield comment tokens.
type Flexer interface {
	Flex() token.Token
}

// Parrot implements Stream by yielding tokens from supplied list
type Parrot struct {
	toks []token.Token
	i    int
}

func FromTokens(toks []token.Token) *Parrot {
	return &Parrot{
		toks: toks,
	}
}

func (p *Parrot) Lex() token.Token {
	if p.i >= len(p.toks) {
		tok := token.Token{Kind: token.EOF}
		if len(p.toks) == 0 {
			return tok
		}
		pin := p.toks[len(p.toks)-1].Pin
		tok.Pin = pin
		return tok
	}
	tok := p.toks[p.i]
	p.i++
	return tok
}

func ListTokens(w io.Writer, s Stream, m source.PinMap) error {
	for {
		tok := s.Lex()
		_, err := io.WriteString(w, token.FormatTokenLine(m, tok))
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, "\n")
		if err != nil {
			return err
		}
		if tok.Kind == token.EOF {
			return nil
		}
	}
}
