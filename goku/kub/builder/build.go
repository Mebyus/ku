package builder

import (
	"fmt"
	"io"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/parser"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/kub/genc"
)

type Context struct {
	Pool *srcmap.Pool

	// Include codegen for tests and gather test names.
	Test bool
}

func genTexts(c *Context, out io.Writer, texts []*srcmap.Text) error {
	g := genc.Gen{State: &genc.State{
		Map:   c.Pool,
		Debug: true,
		Test:  c.Test,
	}}
	g.State.Init()

	// List of test names
	var tests []string

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
				return diag.Format(c.Pool, err.(diag.Error))
			}

			g.Reset()
			g.Nodes(t)
			_, err = g.WriteTo(out)
			if err != nil {
				return err
			}
			if c.Test {
				tests = t.AppendTestNames(tests)
			}
		default:
			return fmt.Errorf("unknown source file extension \"%s\" (%s)", text.Ext, text.Path)
		}
	}

	g.Reset()
	g.NameBooks()
	_, err := g.WriteTo(out)
	if err != nil {
		return err
	}

	if !c.Test {
		return nil
	}

	g.Reset()
	g.MainTestDriver(tests)
	_, err = g.WriteTo(out)
	return err
}

func GenTexts(pool *srcmap.Pool, out io.Writer, texts []*srcmap.Text) error {
	c := Context{Pool: pool}
	return genTexts(&c, out, texts)
}

func GenTextsWithTests(pool *srcmap.Pool, out io.Writer, texts []*srcmap.Text) error {
	c := Context{
		Pool: pool,
		Test: true,
	}
	return genTexts(&c, out, texts)
}
