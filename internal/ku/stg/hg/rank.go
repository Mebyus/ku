package hg

import "fmt"

// Tower is the result of ranking a component.
//
// Tower is formed by a list of cohorts. Each cohort
// contains symbols with the same rank and cohort index
// inside the tower directly corresponds to cohort rank.
type Tower struct {
	// Index in this list equals cohort rank.
	Cohorts []Cohort

	// Total number of symbols in tower.
	Total uint32
}

// Cohort is a collection of symbols with the same rank.
type Cohort struct {
	// Always has at least one element.
	Symbols []Symbol

	// Lowest rank is 0. Only component roots have rank 0.
	Rank uint32
}

// Ranker is a reusable helper object which can rank
// vertices in a given component.
type Ranker struct {
	c *Component

	// List of vertices currently being scanned.
	wave []vid

	// List of vertices for next scan.
	next []vid

	// Stores rank of each vertex.
	//
	// Index in this list corresponds to vertex id (its index).
	rank []uint32

	// Contains counter for each vertex which indicates
	// how many ancestors are still unranked for a vertex with
	// particular index. Index directly corresponds to
	// vertex id (its index).
	left []uint32
}

func (r *Ranker) Reset(c *Component) {
	num := len(c.verts)

	r.c = c
	r.wave = r.wave[:0]
	r.next = r.next[:0]
	r.rank = reset(r.rank, num)
	r.left = reset(r.left, num)
}

// Rank vertices assuming there are no clusters.
//
// Panics if component has clusters.
func (r *Ranker) Rank() *Tower {
	c := r.c
	if len(c.clusters) != 0 {
		panic("component contains clusters")
	}
	if len(c.roots) == 0 {
		panic("no roots in cluster-free component")
	}

	// component is guaranteed to have at least
	// 2 cohorts: roots + pinnacles
	cohorts := make([]Cohort, 0, 2)

	// add root symbols as cohort with rank 0
	{
		roots := make([]Symbol, 0, len(c.roots))
		for _, v := range c.roots {
			roots = append(roots, c.verts[v].sym)
		}
		cohorts = append(cohorts, Cohort{Symbols: roots})
	}

	// number of nodes successfully ranked
	var total int

	for i, v := range c.verts {
		r.left[i] = v.na
	}

	r.wave = append(r.wave, c.roots...) // init starting wave
	for len(r.wave) != 0 {
		if debug {
			fmt.Printf("[rank] wave: %v\n", r.wave)
		}

		for _, v := range r.wave {
			// waiters are vertices still waiting to be ranked
			// we take descendants as waiters because after
			// current was already ranked we should check its
			// descendants, maybe some of them finished ranking
			waiters := c.verts[v].des
			if len(waiters) == 0 {
				continue
			}

			rank := r.rank[v] + 1

			for _, i := range waiters {
				r.left[i] -= 1

				if rank > r.rank[i] {
					// select highest rank from all vertices inside the wave
					r.rank[i] = rank
				}

				// check if waiter vertex has finished ranking
				if r.left[i] == 0 {
					for k := uint32(len(cohorts)); k <= rank; k += 1 {
						// add empty cohorts until we reach current rank
						cohorts = append(cohorts, Cohort{Rank: k})
					}
					cohorts[rank].Symbols = append(cohorts[rank].Symbols, c.verts[i].sym)

					// next wave is constructed from vertices that finished
					// ranking during this wave
					r.next = append(r.next, i)
				}
			}
		}

		total += len(r.wave)
		r.wave, r.next = r.next, r.wave
		r.next = r.next[:0]
	}

	if total < len(c.verts) {
		panic(fmt.Sprintf("not all vertices were ranked %d/%d", total, len(c.verts)))
	} else if total > len(c.verts) {
		panic(fmt.Sprintf("too many reported vertices %d/%d", total, len(c.verts)))
	}

	return &Tower{
		Cohorts: cohorts,
		Total:   uint32(len(c.verts)),
	}
}
