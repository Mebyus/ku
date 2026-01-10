package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
)

// Prune tree of symbols starting from given nodes.
func Prune(symbols []*Symbol) {
	var p SymWalker
	p.init()

	for _, s := range symbols {
		p.add(s)
	}
	p.walk()

	fmt.Printf("tree contains %d symbol(s)\n", len(p.visited))
}

type SymWalker struct {
	ins Inspector

	// Symbols waiting their turn to be processed.
	backlog []*Symbol

	// Contains symbols already visited during graph walk.
	visited map[*Symbol]struct{}
}

func (w *SymWalker) init() {
	w.ins.init()

	w.visited = make(map[*Symbol]struct{})
}

func (w *SymWalker) add(s *Symbol) {
	_, ok := w.visited[s]
	if ok {
		return
	}

	w.visited[s] = struct{}{}
	w.backlog = append(w.backlog, s)
}

func (w *SymWalker) next() *Symbol {
	if len(w.backlog) == 0 {
		return nil
	}

	i := len(w.backlog) - 1
	s := w.backlog[i]
	w.backlog = w.backlog[:i]
	return s
}

func (w *SymWalker) walk() {
	for {
		sym := w.next()
		if sym == nil {
			return
		}

		symbols := w.ins.GetUsedSymbols(sym)
		for _, s := range symbols {
			w.add(s)
		}
	}
}

// Inspector is a helper object for inspecting how symbols depend on each other.
// Reuses internal buffers between
type Inspector struct {
	// set of used symbols for current walk.
	set map[*Symbol]struct{}
}

func NewInspector() *Inspector {
	n := &Inspector{}
	n.init()
	return n
}

func (n *Inspector) init() {
	n.set = make(map[*Symbol]struct{})
}

func (n *Inspector) add(s *Symbol) {
	if s.Kind == smk.Fun || s.Kind == smk.Method || (s.Kind == smk.Var && s.Scope.Kind == sck.Unit) {
		n.set[s] = struct{}{}
	}
}

func (n *Inspector) take() []*Symbol {
	if len(n.set) == 0 {
		return nil
	}
	list := make([]*Symbol, 0, len(n.set))
	for s := range n.set {
		list = append(list, s)
	}
	clear(n.set)
	return list
}

// Get global or unit level symbols used by the given symbol.
// Result can include itself for recursive symbols.
func (n *Inspector) GetUsedSymbols(symbol *Symbol) []*Symbol {
	switch symbol.Kind {
	case smk.Fun, smk.Method:
		n.inspectFun(symbol.Def.(*Fun))
	default:
		panic(fmt.Sprintf("unexpected %s (=%d) symbol \"%s\"", symbol.Kind, symbol.Kind, symbol.Name))
	}

	return n.take()
}

func (n *Inspector) inspectFun(fun *Fun) {
	n.inspectBlock(&fun.Body)
}

func (n *Inspector) inspectBlock(block *Block) {
	for _, s := range block.Nodes {
		n.inspectStatement(s)
	}
}

func (n *Inspector) inspectStatement(stm Statement) {
	switch s := stm.(type) {
	case *Break, *Stub, *Never:
		// do nothing
	case *Var:
		if s.Exp != nil {
			n.inspectExp(s.Exp)
		}
	case *Ret:
		if s.Exp != nil {
			n.inspectExp(s.Exp)
		}
	case *Block:
		n.inspectBlock(s)
	case *Assign:
		n.inspectExp(s.Exp)
		n.inspectExp(s.Target)
	case *If:
		n.inspectIf(s)
	case *Invoke:
		n.inspectExp(s.Call)
	case *DeferCall:
		n.inspectExp(s.Call)
	case *Loop:
		n.inspectBlock(&s.Body)
	case *While:
		n.inspectExp(s.Exp)
		n.inspectBlock(&s.Body)
	case *Must:
		n.inspectExp(s.Exp)
	default:
		panic(fmt.Sprintf("unexpected (%T) statement", s))
	}
}

func (n *Inspector) inspectIf(f *If) {
	for _, b := range f.Branches {
		n.inspectExp(b.Exp)
		n.inspectBlock(&b.Block)
	}

	if f.Else != nil {
		n.inspectBlock(f.Else)
	}
}

func (n *Inspector) inspectExp(exp Exp) {
	switch e := exp.(type) {
	case *Integer, *String, *Nil, *Rune, *Boolean:
		// do nothing
	case *SymExp:
		n.add(e.Symbol)
	case *Call:
		n.inspectCall(e)
	case *Pack:
		n.inspectExpList(e.List)
	case *Cast:
		n.inspectExp(e.Exp)
	case *Binary:
		n.inspectExp(e.A)
		n.inspectExp(e.B)
	case *DerefSelectField:
		n.inspectExp(e.Exp)
	default:
		panic(fmt.Sprintf("unexpected (%T) expression", e))
	}
}

func (n *Inspector) inspectExpList(list []Exp) {
	for _, exp := range list {
		n.inspectExp(exp)
	}
}

func (n *Inspector) inspectCall(c *Call) {
	n.add(c.Symbol)
	n.inspectExpList(c.Args)
}
