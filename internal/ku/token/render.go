package token

import (
	"io"

	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

type Stream interface {
	Lex(*Token)
}

// Render outputs formatted list of token from Stream to a given output writer.
func Render(pool *sx.Pool, out io.Writer, s Stream) error {
	var g fbuf.Buffer

	for {
		var tok Token
		s.Lex(&tok)

		pos := pool.DecodePin(tok.Pin).Pos.String()
		g.Puts(pos)
		g.RepSpace(12 - len(pos))

		tok.Format(&g)
		g.Line()

		if g.Len() >= 1<<13 {
			_, err := g.WriteTo(out)
			if err != nil {
				return err
			}
			g.Reset(0)
		}

		if tok.Kind == EOF {
			_, err := g.WriteTo(out)
			return err
		}
	}
}
