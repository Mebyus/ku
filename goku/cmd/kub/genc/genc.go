package genc

import (
	"fmt"
	"os"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/parser"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/kub/builder"
	"github.com/mebyus/ku/goku/kub/genc"
)

var Butler = &butler.Butler{
	Name: "genc",

	Short: "Generate C code from a given Ku source file",
	Usage: "[options] [file]",

	Exec: exec,
}

func exec(r *butler.Butler, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	path := files[0]
	return gen(path)
}

func gen(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return builder.GenUnit(os.Stdout, path)
	}
	return genFromFile(path)
}

func genFromFile(path string) error {
	pool := srcmap.New()
	text, err := pool.Load(path)
	if err != nil {
		return err
	}

	p := parser.FromText(text)
	t, err := p.Text()
	if err != nil {
		return diag.Format(pool, err.(diag.Error))
	}
	g := genc.Gen{State: &genc.State{
		Map:   pool,
		Debug: true,
		Test:  true,
	}}
	g.State.Init()
	g.Nodes(t)
	_, err = g.WriteTo(os.Stdout)
	return err
}
