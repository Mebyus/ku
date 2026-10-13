package hg

import (
	"slices"
)

const debug = false

// possibly reallocate and clear (set to zero) given slice
// and return a slice with at least specified number of
// zeroed elements
func reset[T any](s []T, num int) []T {
	if len(s) == 0 {
		return make([]T, num)
	}

	clear(s)
	if len(s) < num {
		s = append(s, make([]T, num-len(s))...)
	}
	return s
}

type Builder struct {
	walker cwalker

	// List of all graph nodes. In this slice node index directly
	// corresponds to its id.
	nodes []Node

	// List of all non-trivial (with at least two nodes) connected
	// components within graph. In this slice component index directly
	// corresponds to its id.
	comps []Component

	// List of isolated nodes.
	//
	// Isolated nodes do not have connections to other nodes.
	// Essentially such nodes are trivial graph components with 1 vertex.
	// We keep them separated from non-trivial components for the sake of
	// building algorithm efficiency.
	isols []nid

	// List of self-loop symbols.
	loops []Symbol

	// Stores visited flag for each node during BFS.
	// Index directly corresponds to node id (index).
	vis []bool

	// List of nodes for current BFS scan.
	wave []nid

	// List of nodes for next BFS scan.
	next []nid

	// Maps node id (index) to component vertex id (index).
	//
	// Note that each component has its own indexing for vertices.
	rm []vid

	// Maps symbol id to its node id (index).
	m map[Symbol]nid

	// Number of already visited nodes during BFS.
	vn uint32

	// Max component size (number of vertices) among those that
	// could have clusters in them.
	maxcsize uint32
}

// Reset builder and preallocate resources for
// at least specified number of nodes.
func (g *Builder) Reset(num int) {
	if g.m == nil {
		g.m = make(map[Symbol]nid, max(32, num))
	} else {
		clear(g.m)
	}

	g.nodes = g.nodes[:0]
	g.nodes = slices.Grow(g.nodes, num)

	g.loops = g.loops[:0]
	g.isols = g.isols[:0]
	g.comps = g.comps[:0]
	g.vis = reset(g.vis, num)
	g.rm = reset(g.rm, num)

	g.vn = 0
	g.maxcsize = 0

	// g.ci = 0
	// g.ck = 0
}

// Build builds a graph from previously added nodes data.
func (g *Builder) Build() {
	g.fill()
	g.split()
	g.discover()
}

// Add a symbol with the list of its dependencies (ancestors).
func (g *Builder) Add(sym Symbol, deps []Symbol) {
	g.add(sym, deps)
}

// AddLoop same as Add, but mark corresponding node as self-loop.
func (g *Builder) AddLoop(sym Symbol, deps []Symbol) {
	id := g.add(sym, deps)
	g.nodes[id].flags |= nloop
	g.loops = append(g.loops, sym)
}

// Add symbol and its dependencies stored inside DepSet helper object.
func (g *Builder) Take(s *DepSet) {
	deps := s.list()

	if s.loop {
		g.AddLoop(s.sym, deps)
	} else {
		g.Add(s.sym, deps)
	}
}

// Isolated returns a copy of all isolated symbols inside the graph.
// Caller can take ownership of returned slice.
func (g *Builder) Isolated() []Symbol {
	if len(g.isols) == 0 {
		return nil
	}

	list := make([]Symbol, 0, len(g.isols))
	for _, id := range g.isols {
		list = append(list, g.nodes[id].sym)
	}
	return list
}

// Symbols returns list of symbols stored in graph nodes.
// Symbol index in returned slice equals to node id (its index).
//
// Used for formatting and debugging purposes.
func (g *Builder) Symbols() []Symbol {
	if len(g.nodes) == 0 {
		return nil
	}

	list := make([]Symbol, 0, len(g.nodes))
	for _, node := range g.nodes {
		list = append(list, node.sym)
	}
	return list
}

// Components returns internal buffer with stored graph components.
//
// Caller must not take ownership or modify returned slice.
func (g *Builder) Components() []Component {
	return g.comps
}

// Returns node id of newly created node.
func (g *Builder) add(sym Symbol, deps []Symbol) nid {
	id := g.get(sym)
	g.nodes[id].anc = g.getlist(deps)
	return id
}

// Get node id (index) by its symbol id. If corresponding node
// does not exist yet, then create and store it before returning
// node id.
func (g *Builder) get(sym Symbol) nid {
	id, ok := g.m[sym]
	if ok {
		return id
	}

	id = nid(len(g.nodes))
	g.m[sym] = id
	g.nodes = append(g.nodes, Node{
		sym: sym,
		id:  id,
	})
	return id
}

// Generate or get already stored node ids for each symbol in the given list.
// Returned list is sorted in ascending order.
func (g *Builder) getlist(syms []Symbol) []nid {
	if len(syms) == 0 {
		return nil
	}

	list := make([]nid, 0, len(syms))
	for _, sym := range syms {
		list = append(list, g.get(sym))
	}
	slices.Sort(list)
	return list
}

// fill descendants for each node
func (g *Builder) fill() {
	for i, node := range g.nodes {
		// Inside this loop the following is true:
		//
		//	node.id == i
		//
		// Since node.id is just its index in nodes slice.
		for _, id := range node.anc {
			// Resulting list of descendants in each node
			// will be sorted in ascending order since each
			// node with index i (from outer loop) is visited
			// only once and in ascending order.
			g.nodes[id].des = append(g.nodes[id].des, nid(i))
		}
	}
}
