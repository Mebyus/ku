package genc

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mebyus/ku/internal/ku/builder/kub"
	"github.com/mebyus/ku/internal/ku/stg"
)

type Config struct {
	Build   kub.Build
	LangDir string
}

func Gen(c *Config, out io.Writer, prog *stg.Program) error {
	err := copyFile(out, filepath.Join(c.LangDir, "src/std/core/prelude.c"))
	if err != nil {
		return err
	}

	var buf Buffer
	buf.nl()
	buf.Gen(prog)
	_, err = out.Write(buf.out)

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

type Buffer struct {
	// output buffer.
	out []byte

	// name prefix for top-level symbols
	prefix string

	// maps type to its name in codegen output
	types map[*stg.Type]string

	// indentation level
	ilevel int
}

func (g *Buffer) Gen(prog *stg.Program) {
	g.types = make(map[*stg.Type]string)

	for _, u := range prog.Units {
		g.prefix = u.LinkName + "_"

		// forward declarations for all unit functions
		for _, f := range u.Funs {
			if f.Skip() {
				continue
			}

			g.funstub(f)
			g.nl()
		}

		// definitions for all unit functions except stubs
		for _, f := range u.Funs {
			if f.Skip() || f.IsStub() {
				continue
			}

			g.fun(f)
			g.nl()
		}
	}
}

func (g *Buffer) putb(b byte) {
	g.out = append(g.out, b)
}

func (g *Buffer) puts(s string) {
	g.out = append(g.out, s...)
}

// put decimal formatted integer into output buffer
func (g *Buffer) putn(n uint64) {
	g.puts(strconv.FormatUint(n, 10))
}

func (g *Buffer) space() {
	g.putb(' ')
}

func (g *Buffer) semi() {
	g.putb(';')
}

func (g *Buffer) nl() {
	g.putb('\n')
}

func (g *Buffer) indent() {
	for range g.ilevel {
		g.putb('\t')
	}
}

func (g *Buffer) inc() {
	g.ilevel += 1
}

func (g *Buffer) dec() {
	g.ilevel -= 1
}
