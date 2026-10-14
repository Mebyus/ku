package builder

import "github.com/mebyus/ku/internal/ku/stg/hg"

func rankUnits(units []*unit, roots []uid) (*hg.Tower, []uid) {
	if len(units) == 0 {
		panic("empty list of units")
	}

	g := hg.ImpGraph{
		Units: make([]hg.Unit, 0, len(units)),
		Roots: roots,
	}
	for _, u := range units {
		g.Units = append(g.Units, hg.Unit{
			Des: u.des,
			Anc: uint32(len(u.anc)),
		})
	}
	return hg.RankOrFindCycle(&g)
}
