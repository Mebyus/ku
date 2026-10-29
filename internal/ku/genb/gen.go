package genb

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/enums/symk"
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
		nloc := g.snum
		narg := g.narg
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
			Nloc: nloc,
			Narg: narg,
			Nret: nret,
		})
	}

	return funs
}

type builder struct {
	body []ir.Stm

	// maps local variable symbol to its local slot
	mloc map[*stg.Symbol]uint32

	// maps function symbol to its index name
	fnames map[*stg.Symbol]uint32

	// number of local slots in function frame
	//
	// this is not the same as number of local
	// variables and function parameters function body
	snum uint32

	// number of local slots used by function parameters
	narg uint32
}

func (g *builder) reset() {
	g.snum = 0
	g.narg = 0
	clear(g.mloc)
	g.body = nil
}

func (g *builder) fun(def *stg.FunDef) {
	g.reset()
	g.body = g.block(&def.Body)
}

// add local variable to function frame
func (g *builder) add(s *stg.Symbol) {
	slot := g.snum
	g.mloc[s] = slot

	n := getSlotSize(s.Type)
	g.snum += n
	if s.Kind == symk.Param {
		g.narg += n
	}
}

// returns number of local slots needed to store a type
func getSlotSize(t *stg.Type) uint32 {
	switch t.Kind {
	case typk.Integer:
		// TODO: what about u64 and s64 integers?
		return 1
	case typk.Ref:
		return 1
	case typk.Span:
		return 2
	case typk.Array:
		def := t.Def.(*stg.Array)
		elem := def.Type
		size := max(elem.Size, 4)
		return (size * uint32(def.Num) / 4)
	default:
		panic(fmt.Sprintf("unexpected %s type", t.Kind))
	}
}

func (g *builder) block(block *stg.Block) []ir.Stm {
	for _, s := range block.Scope.Symbols {
		g.add(s)
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
	case *stg.Assign:
		return g.assign(s)
	case *stg.AssignVar:
		return &ir.Assign{
			Name: g.mloc[s.Symbol],
			Exp:  g.exp(s.Exp),
		}
	case *stg.InitVar:
		_, ok := s.Exp.(*stg.Zero)
		if ok {
			// we can skip zero initialization for local variable
			// due to how vm frames work
			return nil
		}

		return &ir.Assign{
			Name: g.mloc[s.Symbol],
			Exp:  g.exp(s.Exp),
		}
	case *stg.IfBlock:
		return &ir.IfBlock{
			Exp:  g.exp(s.Exp),
			Body: g.block(&s.Body),
		}
	case *stg.IfRet:
		var ret ir.Exp
		if s.Ret != nil {
			ret = g.exp(s.Ret)
		}
		return &ir.IfRet{
			Exp: g.exp(s.Exp),
			Ret: ret,
		}
	case *stg.While:
		return &ir.While{
			Exp:  g.exp(s.Exp),
			Body: g.block(&s.Body),
		}
	case *stg.InvokeFun:
		sym := s.Symbol
		name, ok := g.fnames[sym]
		if !ok {
			panic(fmt.Sprintf("unknown function \"%s\"", sym.Name))
		}
		return &ir.Invoke{
			Args: g.args(s.Args),
			Name: name,
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
	case *stg.SymAddr:
		s := e.Symbol
		if !s.Scope.Kind.IsLocal() {
			panic("stub")
		}
		name, ok := g.mloc[s]
		if !ok {
			panic(fmt.Sprintf("unknown local variable \"%s\"", s.Name))
		}

		return &ir.Addr{Name: name}
	case *stg.SliceArray:
		return g.sarr(e)
	case *stg.IndexSpan:
		// TODO: assign correct stride
		return &ir.Index{
			// can use span expression directly
			// since start pointer is stored in
			// first slot of the span struct
			Addr:   g.exp(e.Exp),
			Index:  g.exp(e.Index),
			Stride: 1,
		}
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
		return &ir.Call{
			Args: g.args(e.Args),
			Name: name,
		}
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}

// slice array
func (g *builder) sarr(s *stg.SliceArray) ir.Exp {
	a := g.exp(s.A)
	b := g.exp(s.B)

	switch e := s.Exp.(type) {
	case *stg.SymExp:
		symbol := e.Symbol
		if !symbol.Scope.Kind.IsLocal() {
			panic("stub")
		}

		slot, ok := g.mloc[symbol]
		if !ok {
			panic(fmt.Sprintf("unknown local variable \"%s\"", symbol.Name))
		}

		// TODO: check array element size
		// if it does not fit into one slot
		return &ir.SliceLocal{
			A:      a,
			B:      b,
			Name:   slot,
			Stride: 1,
		}
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}

}

func (g *builder) assign(a *stg.Assign) ir.Stm {
	switch t := a.Target.(type) {
	case *stg.SymDeref:
		s := t.Symbol
		if !s.Scope.Kind.IsLocal() {
			panic("stub")
		}
		name, ok := g.mloc[s]
		if !ok {
			panic(fmt.Sprintf("unknown local variable \"%s\"", s.Name))
		}

		return &ir.AssignAddr{
			Addr: &ir.Local{Name: name},
			Exp:  g.exp(a.Exp),
		}
	default:
		panic(fmt.Sprintf("unexpected %T expression", t))
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
