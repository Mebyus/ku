package builder

import (
	"fmt"
	"io"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/parser"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/klaw/genc"
)

func GenTexts(pool *srcmap.Pool, out io.Writer, texts []*srcmap.Text) error {
	g := genc.Gen{State: &genc.State{
		Map:   pool,
		Debug: true,
		Test:  true,
	}}
	g.State.Init()

	for _, text := range texts {
		switch text.Ext {
		case ".c", ".h":
			_, err := text.WriteTo(out)
			if err != nil {
				return err
			}
		case ".ku":
			var err error
			p := parser.FromText(text)
			t, err := p.Text()
			if err != nil {
				return diag.Format(pool, err.(diag.Error))
			}

			g.Reset()
			g.Nodes(t)
			_, err = g.WriteTo(out)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown source file extension \"%s\" (%s)", text.Ext, text.Path)
		}
	}

	return nil
}
