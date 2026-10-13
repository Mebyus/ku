package hg

import "fmt"

// Split graph into connected components.
func (g *Builder) split() {
	g.markisols()

	if len(g.isols) >= len(g.nodes) {
		// all nodes are isolated
		// we do not need second pass with BFS
		return
	}

	// proceed with BFS on all non-isolated nodes
	for i := range len(g.nodes) {
		if g.vn >= uint32(len(g.nodes)) {
			// all nodes already visited, no need
			// to start new BFS
			return
		}

		id := nid(i)
		if !g.vis[i] {
			g.bfs(id)
		}
	}
}

// mark all isolated nodes
func (g *Builder) markisols() {
	for i := range len(g.nodes) {
		node := g.nodes[i]

		if len(node.anc)+len(node.des) == 0 {
			id := nid(i)
			g.nodes[i].flags |= nisol
			g.isols = append(g.isols, id)
			g.visit(id)
		}
	}
}

// mark node as visited during BFS
func (g *Builder) visit(id nid) {
	g.vis[id] = true
	g.vn += 1
}

func (g *Builder) remap(nodes []nid) []vid {
	if len(nodes) == 0 {
		return nil
	}

	list := make([]vid, 0, len(nodes))
	for _, n := range nodes {
		list = append(list, g.rm[n])
	}
	return list
}

// map (walk through) next connected component starting from specified node.
func (g *Builder) bfs(id nid) {
	cid := cid(len(g.comps))

	if debug {
		fmt.Printf("[hg] cid %d bfs from node %d\n", cid, id)
	}

	// list of vertices in component being mapped
	var verts []Vertex

	var roots []vid
	var pinns []vid

	// reset the slice from previous BFS, but keep
	// underlying memory
	g.wave = g.wave[:0]

	// setup initial scan wave
	g.wave = append(g.wave, id)
	g.visit(id)

	for {
		if debug {
			fmt.Printf("[hg] wave: %v\n", g.wave)
		}

		for _, n := range g.wave {
			g.nodes[n].cid = cid
			vid := vid(len(verts))
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
				sym: g.nodes[n].sym,
				na:  uint32(len(anc)),
			})

			for _, id := range anc {
				if !g.vis[id] {
					g.next = append(g.next, id)
					g.visit(id)
				}
			}
			for _, id := range des {
				if !g.vis[id] {
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

			var flags cflag
			if len(verts) <= len(roots)+len(pinns) {
				// Component has only 2 cohorts, namely:
				//
				//	0 - roots
				//	1 - pinnacles
				//
				// Such component cannot have clusters in it.
				flags |= nocluster
			} else {
				size := uint32(len(verts))
				if size > g.maxcsize {
					g.maxcsize = size
				}
			}

			if debug {
				fmt.Printf("[hg] component %d\n", cid)
				for _, v := range verts {
					fmt.Printf("   | %d: %d\n", v.id, v.nid)
				}
			}

			g.comps = append(g.comps, Component{
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
