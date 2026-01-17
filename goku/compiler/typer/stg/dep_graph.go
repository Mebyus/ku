package stg

import (
	"slices"
)

// NodeId uniquely identifies node within a graph.
// It just stores node index inside graph nodes list.
type NodeId uint32

// CompId uniquely identifies component within a graph.
// It just stores component index inside graph components list.
type CompId uint32

type NodeFlag uint8

const (
	// Node contains self-loop, i.e. has direct dependency on itself.
	NodeLoop NodeFlag = 1 << iota

	// Isolated nodes do not have connections to other nodes.
	// Essentially such nodes are graph components with 1 vertex.
	NodeIsolated
)

// GraphNode describes a node in symbol dependency graph.
type GraphNode struct {
	// List of ancestor (dependency) node ids (indices).
	// For root nodes this list is always empty.
	// Sorted by node id value.
	//
	// These nodes correspond to symbols used by this node's symbol.
	anc []NodeId

	// List of descendant node ids (indices).
	// For pinnacle nodes this list is always empty.
	// Sorted by node id value.
	//
	// These nodes correspond to symbols which use this node's symbol.
	des []NodeId

	// Symbol id associated with this node.
	sid uint

	// Node belongs to this component.
	cid CompId

	id NodeId

	flags NodeFlag
}

// VertId uniquely identifies vertex within a component.
// It just stores vertex index inside component vertices list.
type VertId uint32

// Vertex describes a vertex in graph connected (non-trivial) component.
type Vertex struct {
	des []VertId

	// Symbol id associated with this vertex.
	sid uint

	nid NodeId

	id VertId
}

type CompFlag uint8

const (
	// Component with no clusters.
	CompNoCluster CompFlag = 1 << iota
)

// GraphComp graph connected component.
//
// By definition nodes from different components do not
// have connections between them. More formally:
//
//	if node A belongs to c1 (component 1) and node B belongs to c2 then
//	by definition there is no edge path A -> B or B -> A in this graph
//
// Components separate all graph nodes into equivalence classes with no
// intersections between them.
type GraphComp struct {
	verts []Vertex

	// vertices with no ancestors
	roots []VertId

	// vertices with no descendants
	pinns []VertId

	clusters []Cluster

	id CompId

	flags CompFlag
}

type ClusterId uint32

// cluster element
type celem struct {
	// symbol id
	sid uint

	vert VertId
}

type Cluster struct {
	elems []celem
}

// GraphBuilder helper object for building and ranking dependency graph.
type GraphBuilder struct {
	walker CompWalker

	// List of all graph nodes. In this slice node index directly
	// corresponds to its id.
	nodes []GraphNode

	// List of all non-trivial (with at least two nodes) connected
	// components within graph. In this slice component index directly
	// corresponds to its id.
	comps []GraphComp

	// List of isolated nodes.
	//
	// Isolated nodes do not have connections to other nodes.
	// Essentially such nodes are trivial graph components with 1 vertex.
	// We keep them separated from non-trivial components for the sake of
	// building algorithm efficiency.
	isolated []NodeId

	// List of self-loop symbol ids.
	loops []uint

	// Buffer for returning list of symbol ids to client code.
	// It is reused between various calls and should be used carefully.
	sl []uint

	// Stores visited flag for each node during BFS.
	// Index directly corresponds to node id (index).
	visited []bool

	// List of nodes for current BFS scan.
	wave []NodeId

	// List of nodes for next BFS scan.
	next []NodeId

	// Maps node id (index) to component vertex id (index).
	//
	// Note that each component has its own indexing for vertices.
	rm []VertId

	// Maps symbol id to its node id (index).
	m map[ /* symbol id */ uint]NodeId

	// Number of already visited nodes during BFS.
	vn uint32

	// Max component size among those that could have
	// clusters in them.
	maxCompSize uint32

	// Cluster iterator index for component inside graph.
	ci CompId

	// Cluster iterator index for cluster inside component.
	ck ClusterId
}

func (g *GraphBuilder) init() {
	g.m = make(map[uint]NodeId)
}

// Reset builder and preallocate resources for at least specified number
// of nodes.
func (g *GraphBuilder) reset(num int) {
	g.nodes = make([]GraphNode, 0, num)
	g.comps = nil

	g.vn = 0
	g.maxCompSize = 0
	g.ci = 0
	g.ck = 0
	clear(g.m)
	g.loops = g.loops[:0]
	g.visited = reset(g.visited, num)
	g.rm = reset(g.rm, num)
}

// Add a symbol id with the list of its dependency ids (ancestors).
func (g *GraphBuilder) Add(sid uint, deps []uint) {
	g.add(sid, deps)
}

// AddLoop same as Add, but mark corresponding node as self-loop.
func (g *GraphBuilder) AddLoop(sid uint, deps []uint) {
	id := g.add(sid, deps)
	g.nodes[id].flags |= NodeLoop
	g.loops = append(g.loops, sid)
}

func (g *GraphBuilder) add(sid uint, deps []uint) NodeId {
	id := g.get(sid)
	g.nodes[id].anc = g.getList(deps)
	return id
}

// Get node id (index) by its symbol id. If corresponding node
// does not exist yet, then create and store it before returning
// node id.
func (g *GraphBuilder) get(sid uint) NodeId {
	id, ok := g.m[sid]
	if ok {
		return id
	}

	id = NodeId(len(g.nodes))
	g.m[sid] = id
	g.nodes = append(g.nodes, GraphNode{
		sid: sid,
		id:  id,
	})
	return id
}

func (g *GraphBuilder) getList(sids []uint) []NodeId {
	if len(sids) == 0 {
		return nil
	}

	list := make([]NodeId, 0, len(sids))
	for _, sid := range sids {
		list = append(list, g.get(sid))
	}
	slices.Sort(list)
	return list
}

// Returns list of self-loop symbol ids.
//
// Callers must not write to or take ownership of returned slice.
func (g *GraphBuilder) getLoops() []uint {
	if len(g.loops) == 0 {
		return nil
	}
	return g.loops
}

// Iterator over clusters inside the graph.
//
// Each call returns a list of symbol ids which form non-trivial
// (with at least 2 symbols) cluster.
//
// Iterator stops when empty slice is returned.
func (g *GraphBuilder) cluster() []uint {
	for {
		if int(g.ci) >= len(g.comps) {
			return nil
		}
		if int(g.ck) >= len(g.comps[g.ci].clusters) {
			g.ci += 1
			g.ck = 0
			continue
		}

		elems := g.comps[g.ci].clusters[g.ck].elems
		g.ck += 1

		g.sl = slices.Grow(g.sl[:0], len(elems))
		for _, elem := range elems {
			g.sl = append(g.sl, elem.sid)
		}
		return g.sl
	}
}

func (g *GraphBuilder) fillDescendants() {
	for i, node := range g.nodes {
		// Inside this loop the following is true:
		//
		//	node.id == i
		//
		// Since node.id is just its index in nodes slice.
		for _, id := range node.anc {
			g.nodes[id].des = append(g.nodes[id].des, NodeId(i))
		}
	}
}

// Split graph into connected components.
func (g *GraphBuilder) split() {
	g.markIsolatedNodes()

	if len(g.isolated) >= len(g.nodes) {
		// all nodes are isolated
		// we do not need second pass with BFS
		return
	}

	// proceed with BFS on all non-isolated nodes
	for i := range len(g.nodes) {
		if int(g.vn) >= len(g.nodes) {
			// all nodes already visited, no need
			// to start new BFS
			return
		}

		id := NodeId(i)
		if !g.visited[i] {
			g.bfs(id)
		}
	}
}

func (g *GraphBuilder) markIsolatedNodes() {
	for i := range len(g.nodes) {
		node := g.nodes[i]

		if len(node.anc)+len(node.des) == 0 {
			id := NodeId(i)
			g.nodes[i].flags |= NodeIsolated
			g.isolated = append(g.isolated, id)
			g.visit(id)
		}
	}
}

// mark node as visited during BFS
func (g *GraphBuilder) visit(id NodeId) {
	g.visited[id] = true
	g.vn += 1
}

func (g *GraphBuilder) remap(nodes []NodeId) []VertId {
	if len(nodes) == 0 {
		return nil
	}

	list := make([]VertId, 0, len(nodes))
	for _, n := range nodes {
		list = append(list, g.rm[n])
	}
	return list
}

// map (walk through) next connected component starting from specified node.
func (g *GraphBuilder) bfs(id NodeId) {
	cid := CompId(len(g.comps))

	// list of vertices in component being mapped
	var verts []Vertex

	var roots []VertId
	var pinns []VertId

	// reset the slice from previous BFS, but keep
	// underlying memory
	g.wave = g.wave[:0]

	// setup initial scan wave
	g.wave = append(g.wave, id)

	for {
		for _, n := range g.wave {
			g.nodes[n].cid = cid
			vid := VertId(len(verts))
			g.rm[n] = vid

			anc := g.nodes[n].anc
			if len(anc) == 0 {
				roots = append(roots, vid)
			}
			des := g.nodes[n].des
			if len(des) == 0 {
				pinns = append(pinns, vid)
			}
			verts = append(verts, Vertex{
				id:  vid,
				nid: n,
				sid: g.nodes[n].sid,
			})

			for _, id := range anc {
				if !g.visited[id] {
					g.next = append(g.next, id)
					g.visit(id)
				}
			}
			for _, id := range des {
				if !g.visited[id] {
					g.next = append(g.next, id)
					g.visit(id)
				}
			}
		}

		if len(g.next) == 0 {
			// we gathered all vertices that belong to current component
			// do vertex descendants remap before exiting
			if len(verts) < 2 {
				panic("non-trivial connected component must have at least 2 vertices")
			}

			for k := range len(verts) {
				id := verts[k].nid
				verts[k].des = g.remap(g.nodes[id].des)
			}

			var flags CompFlag
			if len(verts) == len(roots)+len(pinns) {
				// Component has only 2 cohorts, namely:
				//
				//	0 - roots
				//	1 - pinnacles
				//
				// Such component cannot have clusters in it.
				flags |= CompNoCluster
			} else {
				size := uint32(len(verts))
				if size > g.maxCompSize {
					g.maxCompSize = size
				}
			}

			g.comps = append(g.comps, GraphComp{
				verts: verts,
				roots: roots,
				pinns: pinns,
				id:    cid,
				flags: flags,
			})
			return
		}

		// swap scan buffers and prepare for next wave
		g.wave, g.next = g.next, g.wave
		g.next = g.next[:0]
	}
}

func (g *GraphBuilder) discoverClusters() {
	if len(g.comps) == 0 || g.maxCompSize == 0 {
		return
	}

	g.walker.reset(int(g.maxCompSize))
	for i := range len(g.comps) {
		c := &g.comps[i]

		if c.flags&CompNoCluster == 0 {
			g.walker.init(c)
			g.walker.Walk()
		}
	}
}

// Build builds a graph from previously added nodes data.
func (g *GraphBuilder) Build() {
	g.fillDescendants()
	g.split()
	g.discoverClusters()
}

// VertPath keeps track of vertices path inside graph component.
// Operates as stack data structure.
type VertPath struct {
	// stack elements, each element is a vertex id (index).
	s []VertId

	// used to check whether or not vertex is present in the stack (by vertex index)
	m []bool
}

func (s *VertPath) reset(num int) {
	s.s = s.s[:0]
	s.s = slices.Grow(s.s, num/2) // prealloc some space, to reduce number of reallocations
	s.m = reset(s.m, num)
}

func (s *VertPath) has(id VertId) bool {
	return s.m[id]
}

func (s *VertPath) push(id VertId) {
	s.s = append(s.s, id)

	// mark added element as present in stack
	s.m[id] = true
}

// Pop removes element stored on top of the stack and returns it
// to the caller.
func (s *VertPath) pop() VertId {
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
func (s *VertPath) top() VertId {
	return s.s[s.tip()]
}

func (s *VertPath) tip() int {
	return len(s.s) - 1
}

// CompWalker helper object whuch keeps track of internal state
// for clusters discovery inside graph component.
type CompWalker struct {
	path VertPath

	// Stores discovery step number of visited vertices.
	//
	// Equals 0 for vertices which are not yet visited.
	disc []uint32

	// Earliest visited vertex (the vertex with minimum
	// discovery step number) that can be reached from
	// subtree rooted with vertex at corresponding id (index).
	low []uint32

	// component currently being processed
	comp *GraphComp

	// keeps track on number of steps happened during traversal
	step uint32
}

func (w *CompWalker) init(comp *GraphComp) {
	w.reset(len(comp.verts))
	w.comp = comp
}

func (w *CompWalker) reset(num int) {
	w.step = 0
	w.disc = reset(w.disc, num)
	w.low = reset(w.low, num)
	w.path.reset(num)
}

// Walk implements Tarjan’s algorithm for searching Strongly Connected Components
// inside directed graph.
func (w *CompWalker) Walk() {
	for _, i := range w.comp.roots {
		w.walk(i)
	}

	if int(w.step) >= len(w.comp.verts) {
		// all vertices have been discovered
		return

	}

	// we need additional scan due to root cluster(s)
	// inside the component
	for i := range len(w.comp.verts) {
		id := VertId(i)
		if w.disc[i] == 0 {
			w.walk(id)
		}

		if int(w.step) >= len(w.comp.verts) {
			// all vertices have been discovered
			return
		}
	}
}

// recursive depth-first walk starting from specified vertex
func (w *CompWalker) walk(v VertId) {
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

		// cluster number of newly discovered cluster
		cid := ClusterId(len(w.comp.clusters))
		// w.comp.verts[id].cluster = num

		// cluster has at least 2 vertices by definition
		list := make([]celem, 0, 2)
		list = append(list, celem{sid: w.comp.verts[id].sid, vert: id})

		for id != v {
			id = w.path.pop()
			// w.comp.verts[id].cluster = num
			list = append(list, celem{sid: w.comp.verts[id].sid, vert: id})
		}

		// TODO: do we need to sort this list for consistent testing?
		// slices.Sort(list)
		_ = cid

		w.comp.clusters = append(w.comp.clusters, Cluster{
			elems: list,
		})
	}
}

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
