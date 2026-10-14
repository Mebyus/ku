package hg

type UnitID = vid

// ImpGraph carries minimal data to construct and rank
// unit import graph.
type ImpGraph struct {
	// Unit id (Symbol) must match with its index in this slice.
	Units []Unit

	// List of root units.
	//
	// Root units have no imports.
	Roots []UnitID
}

// Unit intermediate structure for ranking unit dependency graph.
type Unit struct {
	// List of unit descendants.
	Des []UnitID

	// Number of unit ancestors.
	Anc uint32
}

// RankOrFindCycle either ranks import graph or finds a cycle in it.
//
// Returns (*Tower, nil) when graph was ranked successfully.
// Returns (nil, cycle) when graph contains a cycle.
//
// Returned cycle is just a list of unit ids which import
// each other in ascending order. This means that list
// of units in returned cycle:
//
//	[A, B, C]
//
// Imports each other like this:
//
//	A -> B -> C
//	A imports B
//	B imports C
//	C imports A
func RankOrFindCycle(g *ImpGraph) (*Tower, []UnitID) {
	var c Component
	fillCompFromUnits(&c, g)

	var r Ranker
	r.Reset(&c)

	tower, cycle := r.RankOrFindCycle()
	if len(cycle) == 0 {
		return tower, nil
	}

	list := make([]UnitID, 0, len(cycle))
	for _, id := range cycle {
		list = append(list, UnitID(id))
	}
	return nil, list
}

func fillCompFromUnits(c *Component, g *ImpGraph) {
	c.roots = g.Roots
	c.verts = make([]Vertex, 0, len(g.Units))
	for i, u := range g.Units {
		c.verts = append(c.verts, Vertex{
			des: u.Des,
			sym: Symbol(i),
			id:  vid(i),
			na:  u.Anc,
		})
	}
}
