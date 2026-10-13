package hg

import "slices"

func (g *Builder) discover() {
	if len(g.comps) == 0 || g.maxcsize == 0 {
		// we proved there are no clusters in this graph
		// with previous checks
		return
	}

	g.walker.reset(int(g.maxcsize))
	for i := range len(g.comps) {
		c := &g.comps[i]

		if c.flags&nocluster == 0 {
			g.walker.init(c)
			g.walker.Walk()
		}
	}
}

// Component walker is a helper object which keeps track of
// internal state for clusters discovery inside graph component.
type cwalker struct {
	path vpath

	// Stores discovery step number of visited vertices.
	//
	// Equals 0 for vertices which are not yet visited.
	disc []uint32

	// Earliest visited vertex (the vertex with minimum
	// discovery step number) that can be reached from
	// subtree rooted with vertex at corresponding id (index).
	low []uint32

	// component currently being processed
	comp *Component

	// keeps track on number of steps happened during traversal
	step uint32
}

func (w *cwalker) init(comp *Component) {
	w.reset(len(comp.verts))
	w.comp = comp
}

func (w *cwalker) reset(num int) {
	w.step = 0
	w.disc = reset(w.disc, num)
	w.low = reset(w.low, num)
	w.path.reset(num)
}

// Walk implements Tarjan’s algorithm for searching strongly connected components
// inside directed graph (we call such components clusters for brewity).
func (w *cwalker) Walk() {
	for _, i := range w.comp.roots {
		w.walk(i)
	}

	if w.step >= uint32(len(w.comp.verts)) {
		// all vertices have been discovered
		return

	}

	// we need additional scan due to root cluster(s)
	// inside the component
	for i := range len(w.comp.verts) {
		id := vid(i)
		if w.disc[i] == 0 {
			w.walk(id)
		}

		if w.step >= uint32(len(w.comp.verts)) {
			// all vertices have been discovered
			return
		}
	}
}

// recursive depth-first walk starting from specified vertex
func (w *cwalker) walk(v vid) {
	w.step += 1
	w.disc[v] = w.step
	w.low[v] = w.step
	w.path.push(v)

	for _, i := range w.comp.verts[v].des {
		if w.disc[i] == 0 {
			// if vertex is not yet visited, traverse its subtree
			w.walk(i)

			// after subtree traversal current vertex
			// should have the lowest low discovery step
			// of all its descendant vertices
			w.low[v] = min(w.low[v], w.low[i])
		} else if w.path.has(i) {
			// this vertex is already present in stack,
			// thus forming a cycle, we must update
			// low discovery step of subtree start
			w.low[v] = min(w.low[v], w.disc[i])
		}
	}

	if w.low[v] == w.disc[v] {
		// we found head vertex of the cluster,
		// pop the stack until reaching head
		id := w.path.pop()

		if id == v {
			// do not keep track of trivial clusters (that contains one vertex)
			// we already handle that with self-loops
			//
			// TODO: perhaps we should simply panic here, since clusters with
			// only one vertex are possible only for self-loop nodes
			return
		}

		// cluster id of newly discovered cluster
		cid := len(w.comp.clusters)

		// cluster has at least 2 vertices by definition
		elems := make([]Element, 0, 2)
		elems = append(elems, Element{
			sym: w.comp.verts[id].sym,
			nid: w.comp.verts[id].nid,
			vid: id,
		})

		for id != v {
			id = w.path.pop()
			elems = append(elems, Element{
				sym: w.comp.verts[id].sym,
				nid: w.comp.verts[id].nid,
				vid: id,
			})
		}

		// TODO: do we need to sort this list for consistent testing?
		// slices.Sort(list)

		// TODO: do we need to store cluster index (id) somewhere?
		_ = cid

		w.comp.clusters = append(w.comp.clusters, Cluster{
			elems: elems,
		})
	}
}

// Keeps track of vertices path inside graph component.
// Operates as stack data structure.
type vpath struct {
	// stack elements, each element is a vertex id (index).
	s []vid

	// used to check whether or not vertex is present in the stack (by vertex index).
	m []bool
}

func (s *vpath) reset(num int) {
	s.s = s.s[:0]
	s.s = slices.Grow(s.s, num/2) // prealloc some space, to reduce number of reallocations
	s.m = reset(s.m, num)
}

func (s *vpath) has(id vid) bool {
	return s.m[id]
}

func (s *vpath) push(id vid) {
	s.s = append(s.s, id)

	// mark added element as present in stack
	s.m[id] = true
}

// Pop removes element stored on top of the stack and returns it
// to the caller.
func (s *vpath) pop() vid {
	id := s.top()

	// shrink stack by one element, but keep underlying space
	// for future use
	s.s = s.s[:s.tip()]

	// mark removed element as no longer present in stack
	s.m[id] = false

	return id
}

// Top returns element stored on top of the stack. Does not alter
// stack state.
func (s *vpath) top() vid {
	return s.s[s.tip()]
}

func (s *vpath) tip() int {
	return len(s.s) - 1
}
