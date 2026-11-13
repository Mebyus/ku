package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/stg/hg"
)

func Prune(prog *Program) {
	// list of starting symbols for program-wide tree walk
	var list []*Symbol
	var nent int // total number of entrypoint symbols
	for _, u := range prog.Units {
		for _, s := range u.Exports {
			list = append(list, s)
		}
		if u.Entry != nil {
			nent += 1
			list = append(list, u.Entry)
			prog.Entry = u.Entry
		}
	}
	switch nent {
	case 0:
		// continue execution
	case 1:
		// continue execution
	default:
		prog.Pool.Report(0, 0, "program contains more than one entrypoint function")
	}
	if len(list) == 0 {
		prog.Pool.Report(0, 0, "program contains no entrypoint or exported function")
		return
	}

	seen := Walk(list)
	for _, u := range prog.Units {
		for _, s := range u.Funs {
			_, ok := seen[s.RawSym()]
			if !ok {
				s.Flags |= SymbolSkip
			}
		}
	}
}

// Walk symbol tree to gather set of all reachable unit-level
// symbols starting from a given list of symbols.
func Walk(list []*Symbol) map[hg.Symbol]struct{} {
	if len(list) == 0 {
		return nil
	}

	init := make([]*Symbol, 0, len(list))
	seen := make(map[hg.Symbol]struct{}, len(list))
	for _, s := range list {
		sym := s.RawSym()
		_, ok := seen[sym]
		if ok {
			continue
		}

		init = append(init, s)
		seen[sym] = struct{}{}
	}

	w := swalker{
		backlog: init,
		seen:    seen,
	}
	w.walk()
	return w.seen
}

// symbol walker that performs BFS for all reachable
// symbols from a given starting list
type swalker struct {
	// list of symbols waiting for inspection
	backlog []*Symbol

	// set of already seen symbols
	seen map[hg.Symbol]struct{}
}

func (w *swalker) walk() {
	for {
		s := w.next()
		if s == nil {
			return
		}

		w.walkFun(s)
	}
}

func (w *swalker) next() *Symbol {
	if len(w.backlog) == 0 {
		return nil
	}

	i := len(w.backlog) - 1
	s := w.backlog[i]
	w.backlog = w.backlog[:i]
	return s
}

func (w *swalker) add(s *Symbol) {
	if s.IsLocal() {
		return
	}

	sym := s.RawSym()
	_, ok := w.seen[sym]
	if ok {
		return
	}

	if s.Kind == symk.Fun {
		// other unit-level symbols do not need further inspection
		// since they do not have walkable dependencies
		//
		// TODO: add methods to condition
		w.backlog = append(w.backlog, s)
	}
	w.seen[sym] = struct{}{}
}

func (w *swalker) walkFun(s *Symbol) {
	def := s.Def.(*FunDef)
	w.walkBlock(&def.Body)
}

func (w *swalker) walkBlock(block *Block) {
	for _, node := range block.Nodes {
		w.walkStm(node)
	}
}

func (w *swalker) walkStm(s Statement) {
	switch s := s.(type) {
	case *Block:
		w.walkBlock(s)
	case *InitVar:
		w.walkExp(s.Exp)
	case *DiscardExp:
		w.walkExp(s.Exp)
	case *InvokeFun:
		w.walkInvokeFun(s)
	default:
		panic(fmt.Sprintf("unexpected %T statement", s))
	}
}

func (w *swalker) walkExp(exp Exp) {
	switch e := exp.(type) {
	case *InvExp, *Integer, *String, *Boolean:
		// do nothing
	case *SymExp:
		w.add(e.Symbol)
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}

func (w *swalker) walkInvokeFun(f *InvokeFun) {
	w.add(f.Symbol)
	for _, arg := range f.Args {
		w.walkExp(arg)
	}
}
