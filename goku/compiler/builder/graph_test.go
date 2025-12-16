package builder

import (
	"slices"
	"testing"

	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
	"github.com/mebyus/ku/goku/graphs"
)

type GraphTestCase struct {
	// Input.
	Units []*stg.Unit

	// Output.
	Cohorts [][]string

	// Output.
	Cycle []string

	Name string
}

func imports(list ...string) []srcmap.ImportSite {
	if len(list) == 0 {
		return nil
	}
	imports := make([]srcmap.ImportSite, 0, len(list))
	for _, s := range list {
		imports = append(imports, srcmap.ImportSite{Path: origin.Local(s)})
	}
	stg.SortImports(imports)
	return imports
}

func unit(p string, list ...string) *stg.Unit {
	return &stg.Unit{
		Path:    origin.Local(p),
		Imports: imports(list...),
	}
}

func dumpCycleForTest(c *graphs.Cycle, units []*stg.Unit) []string {
	if c == nil {
		return nil
	}

	c.Shift()
	list := make([]string, 0, len(c.Nodes))
	for _, n := range c.Nodes {
		list = append(list, units[n].Path.Import)
	}
	return list
}

func dumpCohortsForTest(g *graphs.Graph, units []*stg.Unit) [][]string {
	var r [][]string
	for _, c := range g.Cohorts {
		var list []string
		for _, i := range c {
			list = append(list, units[i].Path.Import)
		}
		r = append(r, list)
	}
	return r
}

func TestGraph(t *testing.T) {
	tests := []GraphTestCase{
		{
			Name: "1 single unit",
			Units: []*stg.Unit{
				unit("foo"),
			},
			Cohorts: [][]string{
				{"foo"},
			},
		},
		{
			Name: "2 two independent units",
			Units: []*stg.Unit{
				unit("foo"),
				unit("bar"),
			},
			Cohorts: [][]string{
				{"foo", "bar"},
			},
		},
		{
			Name: "3 two units",
			Units: []*stg.Unit{
				unit("foo"),
				unit("bar", "foo"),
			},
			Cohorts: [][]string{
				{"foo"},
				{"bar"},
			},
		},
		{
			Name: "4 two units in cycle",
			Units: []*stg.Unit{
				unit("foo", "bar"),
				unit("bar", "foo"),
			},
			Cycle: []string{"bar", "foo"},
		},
		{
			Name: "5 cycle",
			Units: []*stg.Unit{
				unit("fmt"),
				unit("foo", "bar", "fmt"),
				unit("bar", "foo"),
				unit("kar", "bar"),
			},
			Cycle: []string{"bar", "foo"},
		},
		{
			Name: "6 cycle",
			Units: []*stg.Unit{
				unit("foo"),
				unit("bar"),
				unit("fmt", "foo", "bar"),
				unit("kar", "fmt", "b"),
				unit("main", "kar"),
				unit("b", "c"),
				unit("a", "b"),
				unit("c", "a"),
			},
			Cycle: []string{"a", "b", "c"},
		},
		{
			Name: "7 cohort",
			Units: []*stg.Unit{
				unit("foo"),
				unit("bar"),
				unit("fmt", "foo", "bar"),
				unit("kar", "fmt", "bar"),
				unit("main", "kar"),
			},
			Cohorts: [][]string{
				{"foo", "bar"},
				{"fmt"},
				{"kar"},
				{"main"},
			},
		},
	}

	for _, tt := range tests {
		stg.SortAndOrderUnits(tt.Units)
		for _, c := range tt.Cohorts {
			slices.Sort(c)
		}

		var b Bundle
		b.Units = tt.Units
		cycle := dumpCycleForTest(b.makeGraph(), tt.Units)
		if !slices.Equal(cycle, tt.Cycle) {
			t.Errorf("%s: Cycle got = %v, want %v", tt.Name, cycle, tt.Cycle)
			continue
		}
		if len(cycle) != 0 {
			// no need to check cohorts if graph has cycles
			continue
		}

		cohorts := dumpCohortsForTest(&b.Graph, tt.Units)
		if len(cohorts) != len(tt.Cohorts) {
			t.Errorf("%s: len(Cohorts) got = %d, want %d", tt.Name, len(b.Graph.Cohorts), len(tt.Cohorts))
			continue
		}

		for i := range len(tt.Cohorts) {
			a := cohorts[i]
			b := tt.Cohorts[i]
			if !slices.Equal(a, b) {
				t.Errorf("%s: Cohorts[%d] got = %v, want %v", tt.Name, i, a, b)
				break
			}
		}
	}
}
