package genb

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/enums/typk"
	"github.com/mebyus/ku/internal/ku/stg"
	"github.com/mebyus/ku/internal/vm32b/ir"
)

func Gen(prog *stg.Program) []ir.Fun {
	if len(prog.Units) != 1 {
		panic("stub")
	}

	unit := prog.Units[0]
	if len(unit.Funs) == 0 {
		panic("no code")
	}

	fnames := make(map[*stg.Symbol]uint32, len(unit.Funs))
	for i, s := range unit.Funs {
		fnames[s] = uint32(i)
	}
	g := builder{
		fnames: fnames,
		mloc:   make(map[*stg.Symbol]uint32),
	}

	// maps function symbol to its index name
	funs := make([]ir.Fun, 0, len(unit.Funs))
	for _, s := range unit.Funs {
		var nret uint32
		def := s.Def.(*stg.FunDef)
		if def.Result != nil {
			if def.Result.Kind != typk.Integer {
				panic("stub")
			}
			nret = 1
		}

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
			Nret: nret,
		})
	}

	return funs
}

type builder struct {
	body []ir.Stm

	// maps local variable symbol to its index name
	mloc map[*stg.Symbol]uint32

	// maps function symbol to its index name
	fnames map[*stg.Symbol]uint32
}

func (g *builder) reset() {
	clear(g.mloc)
	g.body = nil
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
	case *stg.DiscardExp:
		return &ir.Discard{Exp: g.exp(s.Exp)}
	case *stg.Return:
		return &ir.Return{Exp: g.exp(s.Exp)}
	case *stg.AssignVar:
		return &ir.Assign{
			Name: g.mloc[s.Symbol],
			Exp:  g.exp(s.Exp),
		}
	case *stg.CreateFixed:
		return &ir.Assign{
			Name: g.mloc[s.Symbol],
			Exp:  g.exp(s.Exp),
		}
	case *stg.IfBlock:
		return &ir.IfBlock{
			Exp:  g.exp(s.Exp),
			Body: g.block(&s.Body),
		}
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
		name, ok := g.mloc[s]
		if !ok {
			panic(fmt.Sprintf("unknown local variable \"%s\"", s.Name))
		}

		return &ir.Local{Name: name}
	case *stg.BinExp:
		a := g.exp(e.A)
		b := g.exp(e.B)
		return &ir.BinExp{A: a, B: b, Op: e.Op.Kind}
	case *stg.FunCall:
		s := e.Symbol
		name, ok := g.fnames[s]
		if !ok {
			panic(fmt.Sprintf("unknown function \"%s\"", s.Name))
		}
		return &ir.FunCall{
			Args: g.args(e.Args),
			Name: name,
		}
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}

func (g *builder) args(list []stg.Exp) []ir.Exp {
	if len(list) == 0 {
		return nil
	}

	args := make([]ir.Exp, 0, len(list))
	for _, a := range list {
		args = append(args, g.exp(a))
	}
	return args
}
