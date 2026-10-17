package genb

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/stg"
	"github.com/mebyus/ku/internal/vm32b/ir"
)

func Gen(prog *stg.Program) []ir.Fun {
	if len(prog.Units) != 1 {
		panic("stub")
	}

	unit := prog.Units[0]
	// maps function symbol to its index name
	fnames := make(map[*stg.Symbol]uint32, len(unit.Funs))
	funs := make([]ir.Fun, 0, len(unit.Funs))
	for i, s := range unit.Funs {
		// TODO: move into separate loop, because we need
		// names for all functions in advance for translating
		// function calls
		fnames[s] = uint32(i)

		def := s.Def.(*stg.FunDef)
		if def.Result != nil {
			panic("stub")
		}

		var g builder
		g.fun(def)

		body := g.body
		if len(body) != 0 {
			_, ok := body[0].(*ir.Return)
			if !ok {
				body = append(body, &ir.Return{})
			}
		} else {
			// return with no expression for
			// functions which return nothing
			// we need it to accurately place
			// return instruction at the end
			// of each function
			body = append(body, &ir.Return{})
		}

		funs = append(funs, ir.Fun{
			Body: body,
			// TODO: number of args may be greater for complex param types
			Narg: uint32(len(def.Params)),
		})
	}

	return funs
}

type builder struct {
	body []ir.Stm

	// maps local variable symbol to its index name
	mloc map[*stg.Symbol]uint32
}

func (g *builder) reset() {
	if g.mloc == nil {
		g.mloc = make(map[*stg.Symbol]uint32)
	} else {
		clear(g.mloc)
	}
}

func (g *builder) fun(def *stg.FunDef) {
	g.reset()
	g.body = g.block(&def.Body)
}

func (g *builder) block(block *stg.Block) []ir.Stm {
	for _, s := range block.Scope.Symbols {
		name := uint32(len(g.mloc))
		g.mloc[s] = name
	}

	nodes := make([]ir.Stm, 0, len(block.Nodes))
	for _, node := range block.Nodes {
		s := g.stm(node)
		if s == nil {
			continue
		}

		nodes = append(nodes, s)
	}
	return nodes
}

func (g *builder) stm(s stg.Statement) ir.Stm {
	switch s := s.(type) {
	case *stg.Return:
		return nil
	case *stg.AssignVar:
		return &ir.Assign{
			Name: g.mloc[s.Symbol],
			Exp:  g.exp(s.Exp),
		}
	case *stg.Branch:
		return nil
	case *stg.While:
		return &ir.While{
			Exp:  g.exp(s.Exp),
			Body: g.block(&s.Body),
		}
	default:
		panic(fmt.Sprintf("unexpected %T statement", s))
	}
}

func (g *builder) exp(exp stg.Exp) ir.Exp {
	switch e := exp.(type) {
	case *stg.Integer:
		return &ir.Integer{Val: uint32(e.Val)}
	case *stg.SymExp:
		s := e.Symbol
		if !s.Scope.Kind.IsLocal() {
			panic("stub")
		}
		return &ir.Local{Name: g.mloc[s]}
	case *stg.BinExp:
		a := g.exp(e.A)
		b := g.exp(e.B)
		return &ir.BinExp{A: a, B: b, Op: e.Op.Kind}
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}
