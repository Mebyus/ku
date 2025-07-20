package builder

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
	"github.com/mebyus/ku/goku/klaw/eval"
	"github.com/mebyus/ku/goku/klaw/parser"
)

type Unit struct {
	// Order of elements directly corresponds to file include order
	// in unit build file.
	Texts []*srcmap.Text

	// List of unit imports.
	Imports []eval.Import

	// Does not include unit or main directory prefix.
	Path origin.Path

	Main bool
}

type Module struct {
	Units []*Unit

	Name string

	Main *Unit
}

func loadUnitTexts(pool *srcmap.Pool, env *eval.Env, path string) ([]*srcmap.Text, error) {
	text, err := pool.Load(filepath.Join(path, "unit.klaw"))
	if err != nil {
		return nil, err
	}

	p := parser.FromText(text)
	unit, err := p.Unit()
	if err != nil {
		return nil, diag.Format(pool, err.(diag.Error))
	}
	u, err := eval.EvalUnit(env, unit)
	if err != nil {
		return nil, diag.Format(pool, err.(diag.Error))
	}
	if len(u.Imports) != 0 {
		return nil, fmt.Errorf("unit contains imports")
	}
	if len(u.Includes) == 0 {
		return nil, fmt.Errorf("unit does not include any source files")
	}

	return pool.LoadFromBase(path, u.Includes)
}

func GenUnit(out io.Writer, path string) error {
	pool := srcmap.New()
	texts, err := loadUnitTexts(pool, eval.NewEnv(), path)
	if err != nil {
		return err
	}

	return GenTexts(pool, out, texts)
}

func GenUnitWithTests(out io.Writer, path string) error {
	pool := srcmap.New()
	env := eval.NewEnv()
	env.TestExe = true

	texts, err := loadUnitTexts(pool, env, path)
	if err != nil {
		return err
	}

	return GenTextsWithTests(pool, out, texts)
}
