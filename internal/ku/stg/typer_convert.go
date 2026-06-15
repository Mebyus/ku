package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/scok"
	"github.com/mebyus/ku/internal/ku/enums/symk"
)

func (t *Typer) convert() {
	for _, s := range t.unit.Funs {
		if s.IsStub() {
			continue
		}

		t.convertFun(s.Def.(*FunDef), &t.box.funs[s.Aux])
	}
}

func (t *Typer) convertFun(def *FunDef, f *ast.Fun) {
	t.sig = &def.Signature
	t.convertBlock(&def.Body, &f.Body)
}

// provides context for converting statements and expressions to STG form
type context struct {
	scope *Scope

	// if true means that expression must be evaluated at compile-time
	// and cannot contain runtime values inside
	static bool
}

func (t *Typer) convertBlock(block *Block, b *ast.Block) {
	block.Pin = b.Pin

	if len(b.Nodes) == 0 {
		return
	}

	nodes := make([]Statement, 0, len(b.Nodes))
	for i, n := range b.Nodes {
		s, exit := t.convertNode(&block.Scope, n)
		if s == nil {
			// skip empty statements
			continue
		}

		if exit && i != len(b.Nodes)-1 {
			// TODO: make this a warning
			// TODO: extract pin from statement here
			t.report(0, "dead code after statement that terminates further execution")
		}

		nodes = append(nodes, s)
	}

	if len(nodes) == 0 {
		// discard allocated nodes memory
		return
	}
	block.Nodes = nodes
}

// returned statement can be nil in case of error of if it was
// optimized out (for example it is the case for empty block statement)
//
// if exit return value equals true it means that the statement exits
// function execution (or it never stops like endless loop)
func (t *Typer) convertNode(c *Scope, s ast.Statement) (stm Statement, exit bool) {
	switch s := s.(type) {
	case *ast.Return:
		return t.convertReturn(c, s), true
	// case ast.Var:
	// 	return t.translateVar(s)
	case *ast.Const:
		return t.convertConst(c, s), false
	case *ast.If:
		return t.convertIf(c, s), false
	// case ast.Match:
	// 	return t.translateMatch(s)
	// case ast.Assign:
	// 	return t.translateAssign(s)
	// case ast.Invoke:
	// 	return t.translateInvoke(s)
	// case ast.DeferCall:
	// 	return t.translateDeferCall(s)
	// case ast.Loop:
	// 	return t.translateLoop(s)
	// case ast.While:
	// 	return t.translateWhile(s)
	// case ast.ForRange:
	// 	return t.translateForRange(s)
	// case ast.Must:
	// 	return t.translateMust(s)
	// case ast.Panic:
	// 	return t.translatePanic(s)
	// case ast.Break:
	// 	return t.translateBreak(s)
	// case ast.Gonext:
	// 	return t.translateGonext(s)
	// case ast.Stub:
	// 	return &stg.Stub{Pin: s.Pin}, nil
	// case ast.Never:
	// 	return &stg.Never{Pin: s.Pin}, nil
	case *ast.Block:
		if len(s.Nodes) == 0 {
			// block statement with no statements is equivalent to empty statement
			return nil, false
		}
		var block Block
		block.Scope.Init(scok.Block, c)
		t.convertBlock(&block, s)
		if len(block.Nodes) == 0 {
			// non empty AST block can still result in empty block
			return nil, false
		}

		// TODO: some block statements can result in terminating statement
		return &block, false
	default:
		panic(fmt.Sprintf("unexpected %T statement", s))
	}
}

func (t *Typer) convertConst(s *Scope, c *ast.Const) Statement {
	name := c.Name
	pin := c.Pin

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		t.report(pin, fmt.Sprintf("symbol named \"%s\" was already declared in this block", name))
		return nil
	}

	var typ *Type
	if c.Type != nil {
		typ = t.LookupType(&t.unit.Scope, c.Type)
	}

	exp := t.convertExp(&context{scope: s, static: true}, c.Exp)
	if typ != nil {
		// TODO: typecheck and correction of exp type
	} else {
		// TODO: assign typ according to exp type
	}

	symbol = s.New(symk.Const, name, pin)
	symbol.Type = typ
	symbol.Def = &StaticValue{Exp: exp}

	// constants are saved as symbols with known compile-time value inside scope
	// thus will not need separate constant definition statements in later stages
	return nil
}

func (t *Typer) convertIf(s *Scope, f *ast.If) Statement {
	exp := t.convertExp(&context{scope: s}, f.Exp)
	// TODO: check that exp has boolean type

	m := If{Exp: exp}
	m.Body.Scope.Init(scok.Branch, s)
	t.convertBlock(&m.Body, &f.Body)

	if f.Else == nil || len(f.Else.Nodes) == 0 {
		return &m
	}

	var block Block
	block.Scope.Init(scok.Branch, s)
	t.convertBlock(&block, f.Else)
	if len(block.Nodes) != 0 {
		m.Else = &block
	}

	return &m
}

func (t *Typer) convertReturn(s *Scope, r *ast.Return) Statement {
	pin := r.Pin

	if t.sig.Never {
		t.report(pin, "return in function which never returns according to declaration")
	} else {
		if t.sig.Result == nil {
			if r.Exp != nil {
				t.report(pin, "non-empty return expression in function with declared void result")
			}
		} else {
			if r.Exp == nil {
				t.report(pin, "return with no expression in function that must return something")
			}
		}
	}

	var exp Exp
	if r.Exp != nil {
		exp = t.convertExp(&context{scope: s}, r.Exp)
		t.checkReturnType(t.sig.Result, exp)
		// TODO: typecheck return type against function result type
		// maybe we also need to adjust type of static values here?
	}
	return &Return{
		Pin: pin,
		Exp: exp,
	}
}

func (t *Typer) convertExp(c *context, exp ast.Exp) Exp {
	switch e := exp.(type) {
	// case ast.Nil:
	// return s.Types.MakeNil(e.Pin), nil
	case *ast.Integer:
		return t.makeInteger(e.Pin, e.Val)
	case *ast.True:
		return t.makeBoolean(e.Pin, true)
	case *ast.False:
		return t.makeBoolean(e.Pin, false)
	case *ast.SymExp:
		return t.convertSymExp(c, e)
	case *ast.BinExp:
		return t.convertBinExp(c, e)
	case *ast.ParenExp:
		return t.convertExp(c, e.Exp)
	case *ast.ErrorExp:
		return t.makeInvExp(e.Pin)
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}

func (t *Typer) convertSymExp(c *context, exp *ast.SymExp) Exp {
	name := exp.Name
	pin := exp.Pin

	symbol := c.scope.Lookup(name)
	if symbol == nil {
		t.report(pin, fmt.Sprintf("unknown symbol \"%s\" used as expression", name))
		return t.makeInvExp(pin)
	}

	switch symbol.Kind {
	case symk.Const:
		return symbol.Def.(*StaticValue).Exp
	case
		// symk.Var,
		// symk.Loop,
		symk.Param:

		if c.static {
			t.report(pin, fmt.Sprintf("runtime value symbol \"%s\" used in compile-time expression", name))
			return t.makeInvExp(pin)
		}

		return &SymExp{
			pin:    pin,
			typ:    symbol.Type,
			Symbol: symbol,
		}
	default:
		t.report(pin, fmt.Sprintf("symbol \"%s\" cannot be used as operand or expression", name))
		return t.makeInvExp(pin)
	}
}

func (t *Typer) convertBinExp(c *context, exp *ast.BinExp) Exp {
	a := t.convertExp(c, exp.A)
	b := t.convertExp(c, exp.B)

	e := &BinExp{
		A:   a,
		B:   b,
		Op:  exp.Op,
		pin: a.Pin(), // TODO: maybe we need separate pin for binary expression instead of first operand
	}
	t.checkBinExpType(e)
	return e
}
