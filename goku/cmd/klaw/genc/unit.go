package genc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/klaw/builder"
	"github.com/mebyus/ku/goku/klaw/eval"
	"github.com/mebyus/ku/goku/klaw/parser"
)

func genFromUnit(path string) error {
	pool := srcmap.New()

	text, err := pool.Load(filepath.Join(path, "unit.klaw"))
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
	if len(u.Imports) != 0 {
		return fmt.Errorf("unit contains imports")
	}
	if len(u.Includes) == 0 {
		return fmt.Errorf("unit does not include any source files")
	}

	texts, err := pool.LoadFromBase(path, u.Includes)
	if err != nil {
		return err
	}

	return builder.GenTexts(pool, os.Stdout, texts)
}
