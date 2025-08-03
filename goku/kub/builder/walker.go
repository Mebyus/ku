package builder

import (
	"fmt"
	"path/filepath"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
	"github.com/mebyus/ku/goku/kub/eval"
	"github.com/mebyus/ku/goku/kub/parser"
)

type ResolveConfig struct {
	// Directory where to search for standard library units.
	RootDir string

	// Directory where to search for main unit.
	MainDir string

	// Directory where to search for regular units.
	UnitDir string
}

type Walker struct {
	*BuildConfig
}

func (w *Walker) WalkFrom(items ...srcmap.QueueItem) ([]*srcmap.Unit, error) {
	if len(items) == 0 {
		panic("no init items")
	}

	q := srcmap.NewUnitQueue()
	for _, item := range items {
		q.Add(item)
	}

	for {
		var item srcmap.QueueItem
		ok := q.Next(&item)
		if !ok {
			break
		}

		u, err := w.loadUnit(item.Path)
		if err != nil {
			return nil, err
		}
		u.Path = item.Path
		q.AddUnit(u)
	}

	return q.Units(), nil
}

func (w *Walker) loadUnit(path origin.Path) (*srcmap.Unit, error) {
	p := w.resolveUnitPath(path)
	unit, err := loadUnit(w.Pool, w.Env, p)
	if err != nil {
		return nil, err
	}
	unit.Path = path

	return unit, nil
}

// Transforms unit path to directory path where unit is stored.
func (w *Walker) resolveUnitPath(path origin.Path) string {
	switch path.Origin {
	case 0:
		panic("unspecified origin")
	case origin.Std:
		return filepath.Join(w.RootDir, "src/std", path.Import)
	case origin.Loc:
		return filepath.Join(w.UnitDir, path.Import)
	case origin.Main:
		return filepath.Join(w.MainDir, path.Import)
	default:
		panic(fmt.Sprintf("unexpected origin \"%s\" (=%d)", path.Origin, path.Origin))
	}
}

func loadUnit(pool *srcmap.Pool, env *eval.Env, path string) (*srcmap.Unit, error) {
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
	return &srcmap.Unit{
		Texts:   texts,
		Imports: u.Imports,
	}, nil
}
