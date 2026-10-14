package hg

import (
	"fmt"
	"slices"
)

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
	// only used for finding cycles
	scout scout

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
	//
	// It is also used as indicator of already visited vertices
	// when finding a cycle. For visited vertex stored value is 0.
	left []uint32

	c *Component
}

func (r *Ranker) Reset(c *Component) {
	num := len(c.verts)

	r.c = c
	r.wave = r.wave[:0]
	r.next = r.next[:0]
	r.rank = reset(r.rank, num)
	r.left = reset(r.left, num)
}

func (r *Ranker) RankOrFindCycle() (*Tower, []Symbol) {
	c := r.c
	if len(c.roots) == 0 || len(c.clusters) != 0 {
		// we already know there is a cycle
		// skip ranking attempt
		cycle := r.findcycle(0)
		if len(cycle) == 0 {
			panic("expected a cycle but not found one")
		}
		return nil, cycle
	}

	tower, _ := r.tryrank()
	if tower != nil {
		return tower, nil
	}

	// find vertex that was not visited yet
	for i := range len(c.verts) {
		v := vid(i)
		if r.left[v] != 0 {
			cycle := r.findcycle(v)
			if len(cycle) == 0 {
				panic("expected a cycle but not found one")
			}
			return nil, cycle
		}
	}

	panic("ranking failed but all vertices were visited")
}

// RankNoCycle rank component vertices assuming there are
// no clusters (and consequently no cycles).
//
// Panics if component has clusters.
func (r *Ranker) RankNoCycle() *Tower {
	c := r.c
	if len(c.clusters) != 0 {
		panic("component contains clusters")
	}
	if len(c.roots) == 0 {
		panic("no roots in component with no clusters")
	}

	tower, total := r.tryrank()
	if tower == nil {
		panic(fmt.Sprintf("not all vertices were ranked %d/%d", total, len(c.verts)))
	}

	return tower
}

// attempt to rank component vertices and report whether
// the attempt was successful
//
// Returns (*Tower, <number of vertives>) if ranking was successful.
//
// Returns (nil, <number of visited vertices>) if ranking failed
// and component contains a cycle.
func (r *Ranker) tryrank() (*Tower, int) {
	c := r.c
	if len(c.roots) == 0 {
		panic("cannot rank component with no roots")
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
		// not all vertices were ranked
		// that means component contains cycles
		return nil, total
	}
	if total > len(c.verts) {
		panic(fmt.Sprintf("too many reported vertices %d/%d", total, len(c.verts)))
	}

	return &Tower{
		Cohorts: cohorts,
		Total:   uint32(len(c.verts)),
	}, total
}

// find cycle inside the component starting from
// specified vertex
func (r *Ranker) findcycle(v vid) []Symbol {
	// To find a cycle we use scout (Ariadne thread) algorithm.
	// This is essentially DFS with cycle detection and backtracking.
	c := r.c
	s := &r.scout
	s.reset(c.verts)
	s.push(v)

	for len(s.stack) != 0 {
		step := s.step()

		switch step.act {
		case sasc:
			s.push(step.val)
		case sdes:
			s.pop()
		case scyc:
			list := s.getcycle(step.val)
			cycle := make([]Symbol, 0, len(list))
			for _, v := range list {
				cycle = append(cycle, c.verts[v].sym)
			}
			return cycle
		default:
			panic(step.act)
		}
	}

	return nil
}

func (s *scout) getcycle(v vid) []vid {
	for i, p := range s.stack {
		if p.vid == v {
			// list of vertex ids which form cycle
			var cycle []vid
			for k := len(s.stack) - 1; k >= i; k -= 1 {
				// add nodes in reverse order, since we walk graph by descending edges
				cycle = append(cycle, s.stack[k].vid)
			}
			if len(cycle) < 2 {
				panic("cycle must contain at least 2 vertices")
			}
			return cycle
		}
	}

	panic(fmt.Sprintf("vertex (vid=%d) not found in scout path", v))
}

// scout action type
type sact uint8

const (
	// scout ascends, path length increases
	sasc sact = iota

	// scout descends, path length decreases
	sdes

	// scout found cycle
	scyc
)

// describes scout position
type spos struct {
	// vertex position
	vid vid

	// Index of vertex descendant
	// (to clarify: this value is not a vertex index).
	//
	// This is an index into Vertex.des.
	//
	// Determines which vertex to visit next.
	dx uint32
}

// describes step scout wants to make
type sstep struct {
	// meaning depends on action type (act)
	//
	//	sasc: next node index
	//	sdes: not used
	//	scyc: vertex id in scout path stack
	// 		that connects start and end of the cycle
	val vid

	// action type, discriminator for val field
	act sact
}

// keeps track of scout path and what it already traversed
// inside component
type scout struct {
	verts []Vertex

	// contains current scout path, operates
	// as a stack data structure
	stack []spos

	// Contains include flag for each vertex.
	// Indicates whether vertex is part of the scout path.
	//
	// Element at specific index always corresponds to vertex with
	// the same index.
	inc []bool

	// Contains visited flag for each vertex.
	vis []bool
}

func (s *scout) reset(verts []Vertex) {
	n := len(verts)
	s.verts = verts
	s.stack = s.stack[:0]
	s.stack = slices.Grow(s.stack, n/4)
	s.inc = reset(s.inc, n)
	s.vis = reset(s.vis, n)
}

func (s *scout) step() sstep {
	tip := s.tip()
	des := s.verts[tip.vid].des

	if len(des) == 0 {
		return sstep{act: sdes}
	}

	i := tip.dx
	for i < uint32(len(des)) {
		// next node index in path chosen among current node descendants
		next := des[i]
		if !s.vis[next] {
			s.save(i)
			return sstep{
				val: next,
				act: sasc,
			}
		}

		if s.inc[next] {
			// cycle found
			return sstep{
				val: next,
				act: scyc,
			}
		}

		i += 1
	}

	return sstep{act: sdes}
}

// push new vertex onto the path stack
func (s *scout) push(v vid) {
	s.vis[v] = true
	s.inc[v] = true
	s.stack = append(s.stack, spos{vid: v})
}

func (s *scout) pop() {
	tip := s.tip()
	s.inc[tip.vid] = false
	s.stack = s.stack[:len(s.stack)-1]

	if len(s.stack) == 0 {
		return
	}

	s.save(s.tip().dx + 1)
}

func (s *scout) tip() spos {
	return s.stack[len(s.stack)-1]
}

// update descendant index at path tip
func (s *scout) save(x uint32) {
	tip := s.tip()
	tip.dx = x
	s.stack[len(s.stack)-1] = tip
}
