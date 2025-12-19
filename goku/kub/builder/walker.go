package builder

import (
	"fmt"
	"path/filepath"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/sm"
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

func (w *Walker) WalkFrom(items ...sm.QueueItem) ([]*sm.Unit, error) {
	if len(items) == 0 {
		panic("no init items")
	}

	q := sm.NewUnitQueue()
	for _, item := range items {
		q.Add(item)
	}

	for {
		var item sm.QueueItem
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

func (w *Walker) loadUnit(path sm.UnitPath) (*sm.Unit, error) {
	p := w.resolveUnitPath(path)
	unit, err := loadUnit(w.Pool, w.Env, p)
	if err != nil {
		return nil, err
	}
	unit.Path = path

	return unit, nil
}

// Transforms unit path to directory path where unit is stored.
func (w *Walker) resolveUnitPath(path sm.UnitPath) string {
	switch path.Origin {
	case 0:
		panic("unspecified origin")
	case sm.Std:
		return filepath.Join(w.RootDir, "src/std", path.Import)
	case sm.Loc:
		return filepath.Join(w.UnitDir, path.Import)
	case sm.Main:
		return filepath.Join(w.MainDir, path.Import)
	default:
		panic(fmt.Sprintf("unexpected origin \"%s\" (=%d)", path.Origin, path.Origin))
	}
}

func loadUnit(pool *sm.Pool, env *eval.Env, path string) (*sm.Unit, error) {
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
	return &sm.Unit{
		Texts:   texts,
		Imports: u.Imports,
	}, nil
}
