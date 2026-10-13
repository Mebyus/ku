// package hg (name stands for "hoist graph") implements dependency graph building
// and analysis algorithms for addressing various complications that unit-level symbol
// hoisting brings into the language.
package hg

// Symbol uniquely identifies symbol within program unit.
//
// This package does not provide any facilities for generating
// these ids. It is responsibility of the package user to attach
// unique id for each symbol.
//
// We use integers for symbol identification for the sake of
// code efficiency.
type Symbol uint

// Node id uniquely identifies node within a graph.
// It just stores node index inside graph nodes list.
type nid uint32

// Component id uniquely identifies component within a graph.
// It just stores component index inside graph components list.
type cid uint32

// stores node flags
type nflag uint8

const (
	// Node contains self-loop, i.e. it has direct dependency on itself.
	nloop nflag = 1 << iota

	// Isolated nodes do not have connections to other nodes.
	// Essentially such nodes are graph components with 1 vertex.
	nisol
)

// Node represents a symbol as node in dependency graph.
type Node struct {
	// List of ancestor (dependency) node ids (indices).
	// For root nodes this list is always empty.
	// Sorted by node id value.
	//
	// These nodes correspond to symbols used by this node's symbol.
	anc []nid

	// List of descendant node ids (indices).
	// For pinnacle nodes this list is always empty.
	// Sorted by node id value.
	//
	// These nodes correspond to symbols which use this node's symbol.
	des []nid

	// Symbol id associated with this node.
	sym Symbol

	// Node belongs to this component.
	// Not applicable to isolated nodes.
	cid cid

	// Node id.
	// Stores node index within list of graph nodes.
	id nid

	flags nflag
}

// Vertex id uniquely identifies vertex within a component.
// It just stores vertex index inside component vertices list.
type vid uint32

// Vertex describes a vertex in non-trivial graph component.
type Vertex struct {
	// List of descendant vertices.
	des []vid

	// Symbol id associated with this vertex.
	sym Symbol

	nid nid

	// Vertex id.
	// Stores vertex index within list of component vertices.
	id vid

	// Number of ancestors for this vertex.
	//
	// Equals 0 for root vertices.
	na uint32
}

// stores component flags
type cflag uint8

const (
	// Component with no clusters.
	nocluster cflag = 1 << iota
)

// Component represents connected component within a graph.
//
// By definition nodes from different components do not
// have connections between them. More formally:
//
//	if node A belongs to c1 (component 1) and node B belongs to c2 then
//	by definition there is no edge path A -> B or B -> A in this graph
//
// Components separate all graph nodes into equivalence classes with no
// intersections between them.
type Component struct {
	verts []Vertex

	// List of root vertex ids.
	//
	// Root vertex has no ancestors.
	roots []vid

	// List of pinnacle vertex ids
	//
	// Pinnacle vertex has no descendants.
	pinns []vid

	// List of clusters that constitute this component.
	clusters []Cluster

	id cid

	flags cflag
}

// Symbols returns list of symbols that form the component.
//
// Should be used mostly for formatting and debugging purposes.
func (c *Component) Symbols() []Symbol {
	list := make([]Symbol, 0, len(c.verts))
	for _, v := range c.verts {
		list = append(list, v.sym)
	}
	return list
}

// Element describes vertex as part of the cluster.
type Element struct {
	sym Symbol

	vid vid

	nid nid
}

// Cluster represents strongly connected component within a graph.
//
// By definition for any two elements A and B inside a cluster there exist
// at least one of each edge paths A -> B and B -> A. In other words
// inside a cluster any element is reachable from any other element.
type Cluster struct {
	elems []Element
}
