package builder

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/klaw/eval"
	"github.com/mebyus/ku/goku/klaw/parser"
)

type GenProgramConfig struct {
	// Path to main unit (relative to MainDir).
	Main string

	// Directory where to search for main unit.
	MainDir string

	// Directory where to search for regular units.
	SourceDir string

	Pool *srcmap.Pool

	BuildKind bk.Kind
}

func GenFromMain(out io.Writer, c *GenProgramConfig) error {
	unit, err := loadUnit(c.Pool, nil, filepath.Join(c.MainDir, c.Main))
	if err != nil {
		return err
	}
	unit.Main = true

	queue := NewUnitQueue()
	queue.AddUnit(unit)
	for {
		var item QueueItem
		ok := queue.Next(&item)
		if !ok {
			break
		}

		u, err := loadUnit(c.Pool, nil, filepath.Join(c.SourceDir, item.Path))
		if err != nil {
			return err
		}
		u.Path = item.Path
		queue.AddUnit(u)
	}

	units := queue.units
	m := make(map[string]uint32, len(units))
	for i, u := range units {
		if u.Main {
			continue
		}

		m[u.Path] = uint32(i)
		fmt.Println(u.Path)
	}

	return nil
}

func loadUnit(pool *srcmap.Pool, env eval.Env, path string) (*Unit, error) {
	text, err := pool.Load(filepath.Join(path, "unit.kub"))
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
	if len(u.Includes) == 0 {
		return nil, fmt.Errorf("unit does not include any source files")
	}

	texts, err := pool.LoadFromBase(path, u.Includes)
	if err != nil {
		return nil, err
	}
	return &Unit{
		Texts:   texts,
		Imports: u.Imports,
	}, nil
}
