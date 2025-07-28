package builder

import (
	"fmt"
	"io"
	"os"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bk"
	"github.com/mebyus/ku/goku/compiler/enums/bm"
	"github.com/mebyus/ku/goku/compiler/parser"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
	"github.com/mebyus/ku/goku/graphs"
	"github.com/mebyus/ku/goku/kub/eval"
	"github.com/mebyus/ku/goku/kub/genc"
)

type BuildConfig struct {
	ResolveConfig

	Pool *srcmap.Pool

	BuildKind bk.Kind

	BuildMode bm.Mode
}

func genUnitsToFile(config *BuildConfig, out string, items []srcmap.QueueItem) error {
	genOut, err := os.Create(out)
	if err != nil {
		return err
	}
	defer genOut.Close()

	return genUnits(config, genOut, items)
}

// Generate C code for specified units (by their paths in items).
func genUnits(config *BuildConfig, out io.Writer, items []srcmap.QueueItem) error {
	env := eval.NewEnv()
	env.BuildKind = config.BuildKind
	env.BuildMode = config.BuildMode

	walker := Walker{
		BuildConfig: config,
		Env:         env,
	}
	units, err := walker.WalkFrom(items...)
	if err != nil {
		return err
	}

	m := make(map[origin.Path]uint32, len(units))
	for i, u := range units {
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

	return genTexts(config, out, texts)
}

func genTexts(c *BuildConfig, out io.Writer, texts []*srcmap.Text) error {
	g := genc.Gen{State: &genc.State{
		Map:   c.Pool,
		Debug: c.BuildKind == bk.Debug,
		Test:  c.BuildMode == bm.TestExe,
	}}
	g.State.Init()

	// List of test names
	var tests []string

	for _, text := range texts {
		switch text.Ext {
		case ".c", ".h":
			_, err := text.WriteTo(out)
			if err != nil {
				return err
			}
		case ".ku":
			var err error
			p := parser.FromText(text)
			t, err := p.Text()
			if err != nil {
				return diag.Format(c.Pool, err.(diag.Error))
			}

			g.Reset()
			g.Nodes(t)
			_, err = g.WriteTo(out)
			if err != nil {
				return err
			}
			if c.BuildMode == bm.TestExe {
				tests = t.AppendTestNames(tests)
			}
		default:
			return fmt.Errorf("unknown source file extension \"%s\" (%s)", text.Ext, text.Path)
		}
	}

	g.Reset()
	g.NameBooks()
	_, err := g.WriteTo(out)
	if err != nil {
		return err
	}

	if c.BuildMode != bm.TestExe {
		return nil
	}

	g.Reset()
	g.MainTestDriver(tests)
	_, err = g.WriteTo(out)
	return err
}
