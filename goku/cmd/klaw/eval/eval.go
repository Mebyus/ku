package eval

import (
	"fmt"

	"github.com/mebyus/ku/goku/butler"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/klaw/eval"
	"github.com/mebyus/ku/goku/klaw/parser"
)

var Butler = &butler.Butler{
	Name: "eval",

	Short: "eval a given unit build file and print results",
	Usage: "[options] [file]",

	Exec: exec,
}

func exec(r *butler.Butler, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	path := files[0]
	return evalFile(path)
}

func evalFile(path string) error {
	pool := srcmap.New()
	text, err := pool.Load(path)
	if err != nil {
		return err
	}

	p := parser.FromText(text)
	unit, err := p.Unit()
	if err != nil {
		return diag.Format(pool, err.(diag.Error))
	}
	u, err := eval.EvalUnit(nil, unit)
	if err != nil {
		return diag.Format(pool, err.(diag.Error))
	}
	fmt.Printf("imports:  %v\n", u.Imports)
	fmt.Printf("includes: %v\n", u.Includes)
	return nil
}
