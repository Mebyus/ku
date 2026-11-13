package genc

import (
	"io"
	"os"
	"path/filepath"

	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/builder/kub"
	"github.com/mebyus/ku/internal/ku/stg"
)

type Config struct {
	Build   kub.Build
	LangDir string
}

func WriteProg(c *Config, out io.Writer, prog *stg.Program) error {
	err := copyFile(out, filepath.Join(c.LangDir, "src/std/core/prelude.c"))
	if err != nil {
		return err
	}

	var g genbuf
	g.line()
	g.Gen(prog)
	_, err = g.buf.WriteTo(out)

	if c.Build.Goal == kub.IntExe {
		err = copyFile(out, filepath.Join(c.LangDir, "src/std/core/int_exe_main.c"))
		if err != nil {
			return err
		}
	}

	return err
}

func copyFile(w io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	return err
}

type genbuf struct {
	buf fbuf.Buffer

	// name prefix for top-level symbols
	prefix string

	// maps type to its name in codegen output
	types map[*stg.Type]string
}

func (g *genbuf) Gen(prog *stg.Program) {
	g.types = make(map[*stg.Type]string)

	for _, u := range prog.Units {
		g.prefix = u.LinkName + "_"

		// forward declarations for all unit functions
		for _, f := range u.Funs {
			if f.Skip() {
				continue
			}

			g.funstub(f)
			g.line()
		}

		// definitions for all unit functions except stubs
		for _, f := range u.Funs {
			if f.Skip() || f.IsStub() {
				continue
			}

			g.fun(f)
			g.line()
		}
	}
}

func (g *genbuf) putb(b byte) {
	g.buf.Putb(b)
}

func (g *genbuf) puts(s string) {
	g.buf.Puts(s)
}

// put decimal formatted integer into output buffer
func (g *genbuf) putn(n uint64) {
	g.buf.Putn(n)
}

func (g *genbuf) space() {
	g.buf.Space()
}

func (g *genbuf) semi() {
	g.buf.Semi()
}

func (g *genbuf) line() {
	g.buf.Line()
}

func (g *genbuf) indent() {
	g.buf.Indent()
}

func (g *genbuf) inc() {
	g.buf.Inc()
}

func (g *genbuf) dec() {
	g.buf.Dec()
}
