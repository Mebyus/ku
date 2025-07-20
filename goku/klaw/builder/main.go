package builder

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
	"github.com/mebyus/ku/goku/graphs"
	"github.com/mebyus/ku/goku/klaw/eval"
	"github.com/mebyus/ku/goku/klaw/parser"
)

type GenProgramConfig struct {
	// Path to main unit (relative to MainDir).
	Main string

	// Directory where to search for standard library units.
	RootDir string

	// Directory where to search for main unit.
	MainDir string

	// Directory where to search for regular units.
	SourceDir string

	Pool *srcmap.Pool

	BuildKind bk.Kind
}

func (c *GenProgramConfig) resolveUnitPath(p origin.Path) string {
	switch p.Origin {
	case 0:
		panic("unspecified origin")
	case origin.Std:
		return filepath.Join(c.RootDir, "src/std", p.Import)
	case origin.Loc:
		return filepath.Join(c.SourceDir, p.Import)
	default:
		panic(fmt.Sprintf("unexpected origin \"%s\" (=%d)", p.Origin, p.Origin))
	}
}

func GenFromMain(out io.Writer, c *GenProgramConfig) error {
	if c.RootDir == "" {
		panic("empty root dir")
	}
	if c.SourceDir == "" {
		panic("empty source dir")
	}

	env := eval.NewEnv()
	env.Exe = true

	unit, err := loadUnit(c.Pool, env, filepath.Join(c.MainDir, c.Main))
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

		u, err := loadUnit(c.Pool, env, c.resolveUnitPath(item.Path))
		if err != nil {
			return err
		}
		u.Path = item.Path
		queue.AddUnit(u)
	}

	units := queue.units
	m := make(map[origin.Path]uint32, len(units))
	for i, u := range units {
		if u.Main {
			continue
		}

		m[u.Path] = uint32(i)
	}

	var g graphs.Graph
	g.Nodes = make([]graphs.Node, len(units))
	g.Rank = make([]uint32, len(units))

	for i, unit := range units {
		// i = unit.Index inside this loop, because we sorted
		// and indexed units beforehand

		g.Nodes[i].Anc = make([]uint32, 0, len(unit.Imports))
		for _, s := range unit.Imports {
			u, ok := m[s.Path]
			if !ok {
				panic(fmt.Sprintf("imported unit \"%s\" not found", s.Path))
			}
			if u == uint32(i) {
				panic(fmt.Sprintf("unit \"%s\" imported itself", s.Path))
			}

			g.Nodes[i].AddAnc(u)
			g.Nodes[u].AddDes(uint32(i))
		}

		if len(unit.Imports) == 0 {
			g.Roots = append(g.Roots, uint32(i))
		}
	}

	var s graphs.Scout
	cycle := s.RankOrFindCycle(&g)
	if cycle != nil {
		return fmt.Errorf("import cycle: %v", cycle.Nodes)
	}

	var texts []*srcmap.Text
	for _, c := range g.Cohorts {
		for _, i := range c {
			u := units[i]
			texts = append(texts, u.Texts...)
		}
	}

	return GenTexts(c.Pool, out, texts)
}

func loadUnit(pool *srcmap.Pool, env *eval.Env, path string) (*Unit, error) {
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
