package lexer

import (
	"fmt"
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
	t []token.Token
	i int
}

func FromTokens(tokens []token.Token) *Parrot {
	return &Parrot{t: tokens}
}

// Reset this stream to the first token.
func (p *Parrot) Reset() {
	p.i = 0
}

func (p *Parrot) Lex() token.Token {
	if p.i < len(p.t) {
		tok := p.t[p.i]
		p.i++
		return tok
	}

	// Always return EOF if original list of tokens was exhausted
	tok := token.Token{Kind: token.EOF}
	if len(p.t) == 0 {
		return tok
	}
	pin := p.t[len(p.t)-1].Pin
	tok.Pin = pin
	return tok
}

func Render(w io.Writer, s Stream, m source.PinMap) error {
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

// ListTokens reads tokens from a stream until EOF token is reached.
//
// Each token obtained from the stream is saved in returned list.
// List always includes EOF token as last element.
func ListTokens(s Stream) []token.Token {
	var list []token.Token
	for {
		tok := s.Lex()
		list = append(list, tok)
		if tok.Kind == token.EOF {
			return list
		}
	}
}

// Compare two streams of tokens. Returns error upon encountering the first pair
// of unequal tokens.
//
// Returns nil if two streams yeild the same tokens (including EOF).
func Compare(a, b Stream) error {
	var n int
	for {
		n += 1
		t1 := a.Lex()
		t2 := b.Lex()

		err := token.Compare(t1, t2)
		if err != nil {
			return fmt.Errorf("compare tokens (n=%d) %s", n, err)
		}

		if t1.Kind == token.EOF {
			return nil
		}
	}
}
