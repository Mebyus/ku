package genc

import (
	"io"
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
	var g genbuf
	err := g.file(filepath.Join(c.LangDir, "src/std/core/prelude.c"))
	if err != nil {
		return err
	}

	g.line()
	err = g.file(filepath.Join(c.LangDir, "src/std/core/os_linux_amd64_1.c"))
	if err != nil {
		return err
	}

	g.line()
	g.Gen(prog)

	switch c.Build.Goal {
	case kub.Exe:
		g.genMain(prog.Entry)
	case kub.IntExe:
		err = g.file(filepath.Join(c.LangDir, "src/std/core/int_exe_main.c"))
		if err != nil {
			return err
		}
	}

	_, err = g.buf.WriteTo(out)
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

		for _, t := range u.Types {
			g.typedef(t)
			g.line()
		}

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

// write file contents into stored output
func (g *genbuf) file(path string) error {
	_, err := g.buf.ReadFile(path)
	return err
}
