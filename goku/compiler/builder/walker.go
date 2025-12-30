package builder

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/parser"
	"github.com/mebyus/ku/goku/compiler/sm"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func Walk(cfg WalkConfig, init ...QueueItem) (*Bundle, diag.Error) {
	w := Walker{
		WalkConfig: cfg,

		pool: sm.New(),
	}
	w.Bundle.Pool = w.pool

	err := w.WalkFrom(init...)
	if err != nil {
		return nil, err
	}

	cycle := w.Bundle.makeGraph()
	if cycle != nil {
		return nil, &diag.ImportCycleError{
			Sites: convertImportCycle(cycle, w.Bundle.Units),
		}
	}

	return &w.Bundle, nil
}

type BaseDirs struct {
	// Base directory for standard library units lookup.
	Std string

	// Base directory for local units lookup.
	Loc string
}

type WalkConfig struct {
	Dir BaseDirs
}

type Walker struct {
	Bundle Bundle

	WalkConfig

	pool *sm.Pool
}

func (w *Walker) WalkFrom(init ...QueueItem) diag.Error {
	if len(init) == 0 {
		panic("no init items")
	}

	q := NewUnitQueue()

	for _, item := range init {
		q.Add(item)
	}

	for {
		var item QueueItem
		if !q.Next(&item) {
			w.Bundle.Units = q.Sorted()
			return nil
		}

		u, err := w.AnalyzeUnit(item)
		if err != nil {
			return err
		}

		if u.HasMain() {
			if u.DiscoveryIndex != 0 {
				panic("not implemented")
				// return fmt.Errorf("main unit [%s] cannot be imported", u.Path)
			}
			if w.Bundle.Main != nil {
				panic("multiple main units in uwalk graph")
			}
			w.Bundle.Main = u
		}
		q.AddUnit(u)
	}
}

func (w *Walker) AnalyzeUnit(item QueueItem) (*stg.Unit, diag.Error) {
	path := item.Path

	dir, err := w.Resolve(path)
	if err != nil {
		return nil, err
	}
	files, loadErr := w.pool.LoadDir(dir, &sm.DirScanParams{IncludeTestFiles: item.IncludeTestFiles})
	if loadErr != nil {
		return nil, &diag.SimpleMessageError{
			Pin:  item.Pin,
			Text: fmt.Sprintf("load unit \"%s\": %s", item.Path, loadErr),
		}
	}

	var imports []sm.ImportSite
	parsers := make([]*parser.Parser, 0, len(files))
	pset := sm.NewPathSet()
	for _, file := range files {
		p := parser.FromText(file)
		build, err := p.Build()
		if err != nil {
			return nil, err
		}
		if build != nil {
			panic("not implemented")
		}
		blocks, err := p.ImportBlocks()
		if err != nil {
			return nil, err
		}

		for _, block := range blocks {
			for _, m := range block.Imports {
				p := sm.UnitPath{
					Origin: block.Origin,
					Import: m.String.Val,
				}
				if p == path {
					return nil, &diag.SimpleMessageError{
						Pin:  m.String.Pin,
						Text: fmt.Sprintf("unit \"%s\" imports itself", p),
					}
				}

				if pset.Has(p) {
					return nil, &diag.SimpleMessageError{
						Pin:  m.String.Pin,
						Text: fmt.Sprintf("multiple imports of the same unit \"%s\"", p),
					}
				}
				pset.Add(p)

				imports = append(imports, sm.ImportSite{
					Path: p,
					Name: m.Name.Str,
					Pin:  m.Name.Pin,
				})
			}
		}
		parsers = append(parsers, p)
	}

	stg.SortImports(imports)
	w.Bundle.Source = append(w.Bundle.Source, parsers)
	return &stg.Unit{
		Path:    path,
		Imports: imports,
	}, nil
}

// Resolve returns system path to directory which contains unit source files.
func (w *Walker) Resolve(path sm.UnitPath) (string, diag.Error) {
	o := path.Origin
	switch path.Origin {
	case 0:
		panic("empty path")
	case sm.Std:
		return w.Dir.Std + "/" + path.Import, nil
	case sm.Pkg:
		panic("not implemented")
	case sm.Loc:
		return w.Dir.Loc + "/" + path.Import, nil
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) origin", o, o))
	}
}
