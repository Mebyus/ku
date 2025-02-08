package builder

import (
	"fmt"
	"slices"

	"github.com/mebyus/ku/goku/compiler/source/origin"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

type Node struct {
	// List of ancestor nodes indices. For root nodes this list is always empty.
	// Sorted by node index.
	//
	// These nodes correspond to units imported by this node's unit.
	Anc []uint32

	// List of descendant nodes indices. For pinnacle nodes this list is always empty.
	// Sorted by link node index.
	//
	// These nodes correspond to units which import this node's unit.
	Des []uint32
}

func (n *Node) addAnc(i uint32) {
	n.Anc = append(n.Anc, i)
}

func (n *Node) addDes(i uint32) {
	n.Des = append(n.Des, i)
}

// Graph
//
// Each node is named by its index inside the Nodes slice.
type Graph struct {
	// Index in this slice corresponds to order of import units.
	Nodes []Node

	// Stores rank of each node.
	//
	// Index in this slice corresponds to node index.
	Rank []uint32

	// List of node indices.
	Roots []uint32

	// Each cohort is a list of node indices.
	Cohorts [][]uint32
}

type Cycle struct {
	// Contains node indices.
	// Always has at least two nodes.
	// Always starts with minimal index in the cycle.
	Nodes []uint32
}

// Performs a cyclical shift of nodes inside the cycle.
// After this operation the node with the minimum index will be placed first.
func (c *Cycle) shift() {
	m := 0 // index of minimal element in nodes
	v := c.Nodes[m]
	for i, n := range c.Nodes {
		if n < v {
			v = n
			m = i
		}
	}

	if m == 0 {
		return
	}

	l := len(c.Nodes)
	nodes := make([]uint32, l)
	copy(nodes, c.Nodes[m:])
	copy(nodes[l-m:], c.Nodes[:m])
	c.Nodes = nodes
}

// Fills Unit.Imports.Units according to import paths.
func (b *Bundle) mapGraphNodes() {
	m := make(map[origin.Path]*stg.Unit, len(b.Units))
	b.Graph.Nodes = make([]Node, len(b.Units))
	b.Graph.Rank = make([]uint32, len(b.Units))
	b.Context.Map = m

	for _, unit := range b.Units {
		m[unit.Path] = unit
	}

	for i, unit := range b.Units {
		// i = unit.Index inside this loop, because we sorted
		// and indexed units beforehand

		for _, s := range unit.Imports {
			u, ok := m[s.Path]
			if !ok {
				panic(fmt.Sprintf("imported unit \"%s\" not found", s.Path))
			}
			if u == unit {
				panic("unit imported itself")
			}

			b.Graph.Nodes[i].addAnc(u.Index)
			b.Graph.Nodes[u.Index].addDes(uint32(i))
		}

		if len(unit.Imports) == 0 {
			b.Graph.Roots = append(b.Graph.Roots, uint32(i))
		}
	}
}

// Scout traverses graph of unit imports. If graph has a cycle than Scout
// reports it. Otherwise each graph node is ranked and placed into corresponding
// cohort.
//
// This is mostly a helper struct that carries various state while we analyze the graph.
type Scout struct {
	// This is a stack. It stores current path of the scout inside the graph.
	// This is used to detect cycles and do backtracking.
	path []ScoutPos

	// Indicates whether a node is part of the scout path.
	//
	// Element at specific index always corresponds to node with
	// the same index.
	contains []bool

	// Indicates whether a node was already visited or not.
	//
	// Element at specific index always corresponds to node with
	// the same index.
	visited []bool

	// Graph that is being traversed.
	// Used as final result if graph does not contain cycles.
	*Graph
}

type ScoutPos struct {
	// Node index.
	i uint32

	// Index of node descendant (to clarify: this value is not a node index).
	//
	// This is an index into Node.Des.
	des uint32
}

func (s *Scout) rankOrFindCycle(g *Graph) *Cycle {
	s.Graph = g
	s.visited = make([]bool, len(g.Nodes))

	if len(g.Roots) == 0 {
		return s.findCycle()
	}

	// number of nodes successfully ranked
	var total int

	// how many ancestors are still unranked for a node with
	// particular index
	left := make([]uint32, len(g.Nodes))
	for i, n := range g.Nodes {
		left[i] = uint32(len(n.Anc))
	}

	// nodes to scan in this wave
	wave := slices.Clone(g.Roots)
	g.Cohorts = make([][]uint32, 0, 2)
	g.Cohorts = append(g.Cohorts, g.Roots)

	// buffer for next wave
	var next []uint32

	for len(wave) != 0 {
		for _, i := range wave {
			s.visited[i] = true

			waiters := g.Nodes[i].Des
			if len(waiters) == 0 {
				continue
			}

			// rank that will be passed to waiters
			rank := g.Rank[i] + 1

			for _, j := range waiters {
				left[j] -= 1

				if rank > g.Rank[j] {
					// select highest rank from all nodes inside the wave
					g.Rank[j] = rank
				}

				// check if waiter node has finished ranking
				if left[j] == 0 {
					for k := uint32(len(g.Cohorts)); k <= rank; k += 1 {
						g.Cohorts = append(g.Cohorts, nil)
					}
					g.Cohorts[rank] = append(g.Cohorts[rank], j)

					// next wave is constructed from nodes that finished
					// ranking during this wave
					next = append(next, j)
				}
			}
		}

		total += len(wave)
		wave, next = next, wave
		next = next[:0]
	}

	if total < len(g.Nodes) {
		return s.findCycle()
	}

	for _, c := range g.Cohorts[1:] {
		slices.Sort(c)
	}

	return nil
}

func (s *Scout) findCycle() *Cycle {
	g := s.Graph
	s.contains = make([]bool, len(g.Nodes))

	for i := range len(g.Nodes) {
		i := uint32(i)
		if !s.visited[i] {
			c := s.traverse(i)
			if c != nil {
				c.shift()
				return c
			}
		}
	}

	panic("unreachable")
}

func (s *Scout) traverse(i uint32) *Cycle {
	s.push(i)

	for len(s.path) != 0 {
		step := s.step()

		switch step.kind {
		case ascend:
			s.push(step.val)
		case descend:
			s.pop()
		case cycle:
			return s.handle(step.val)
		default:
			panic(step.kind)
		}
	}

	return nil
}

type StepKind uint8

const (
	// scout ascends, its path length increases
	ascend StepKind = iota

	// scout descends, its path length decreases
	descend

	// scout found cycle in its path
	cycle
)

type ScoutStep struct {
	// meaning depends on kind
	//
	//	ascend: next node index
	//	descend: ignored
	//	cycle: node index from scout path, that connects start and end of the cycle
	val uint32

	kind StepKind
}

func (s *Scout) step() ScoutStep {
	tip := s.tip()
	des := s.Nodes[tip.i].Des

	if len(des) == 0 {
		return ScoutStep{kind: descend}
	}

	j := tip.des
	for j < uint32(len(des)) {
		// next node index in path chosen among current node descendants
		next := des[j]
		if !s.visited[next] {
			return ScoutStep{
				val:  next,
				kind: ascend,
			}
		}

		if s.contains[next] {
			// cycle found
			return ScoutStep{
				val:  next,
				kind: cycle,
			}
		}

		j += 1
		s.save(j)
	}

	return ScoutStep{kind: descend}
}

func (s *Scout) handle(i uint32) *Cycle {
	for j, pos := range s.path {
		if pos.i == i {
			// list of Node indices which form cycle
			var nodes []uint32
			for k := len(s.path) - 1; k >= j; k -= 1 {
				// add nodes in reverse order, since we walk graph by descend edges
				nodes = append(nodes, s.path[k].i)
			}
			if len(nodes) < 2 {
				panic("unreachable")
			}
			return &Cycle{Nodes: nodes}
		}
	}

	panic("unreachable")
}

// push node index onto the path stack
func (s *Scout) push(i uint32) {
	s.visited[i] = true
	s.contains[i] = true
	s.path = append(s.path, ScoutPos{i: i})
}

func (s *Scout) pop() {
	tip := s.tip()
	s.contains[tip.i] = false
	s.path = s.path[:len(s.path)-1]

	if len(s.path) == 0 {
		return
	}

	tip = s.tip()
	tip.des += 1
	s.path[len(s.path)-1] = tip
}

func (s *Scout) tip() ScoutPos {
	return s.path[len(s.path)-1]
}

// update descendant index at path tip
func (s *Scout) save(d uint32) {
	tip := s.tip()
	tip.des = d
	s.path[len(s.path)-1] = tip
}

func convertImportCycle(c *Cycle, units []*stg.Unit) []stg.ImportSite {
	if len(c.Nodes) < 2 {
		panic("bad cycle data")
	}

	sites := make([]stg.ImportSite, 0, len(c.Nodes))
	for i := 0; i < len(c.Nodes)-1; i += 1 {
		j := c.Nodes[i]
		k := c.Nodes[i+1]

		u := units[j]
		next := units[k]

		s, ok := u.FindImportSite(next.Path)
		if !ok {
			panic(fmt.Sprintf("unable to find \"%s\" import inside \"%s\"", next.Path, u.Path))
		}

		sites = append(sites, s)
	}
	return sites
}
