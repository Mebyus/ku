package token

import (
	"fmt"
	"io"

	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/char"
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

		cat := tok.Kind.Category()
		g.Puts(cat)
		g.RepSpace(12 - len(cat))

		var text string
		switch tok.Kind {
		case Word:
			text = tok.Data
		case Integer:
			text = tok.Data
		case String:
			// TODO: refactor into just fbuf usage with Len() storing and Truncate() rollback
			text = "\"" + char.Escape(tok.Data) + "\""
		}
		g.Puts(text)
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

// String returns token textual representation (without position info).
//
// Should be used mainly for error reporting.
func (t *Token) String() string {
	cat := t.Kind.Category()
	switch t.Kind {
	case INV, EOF:
		return cat
	case Word:
		return fmt.Sprintf("%s %s", cat, t.Data)
	case Integer:
		if t.Data != "" {
			return fmt.Sprintf("%s %s", cat, t.Data)
		}

		switch t.Val {
		case DecInt:
			return fmt.Sprintf("%s %d", cat, t.Val)
		case BinInt:
			return fmt.Sprintf("%s %b", cat, t.Val)
		case OctInt:
			return fmt.Sprintf("%s %o", cat, t.Val)
		case HexInt:
			return fmt.Sprintf("%s %X", cat, t.Val)
		default:
			return fmt.Sprintf("%s %d", cat, t.Val)
		}
	case String:
		return fmt.Sprintf("%s \"%s\"", cat, t.Data) // TODO: need to escape this
	}

	return cat // TODO: literals for keywords and other things
}
