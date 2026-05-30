package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/scok"
)

func (t *Typer) convert() {
	for _, s := range t.unit.Funs {
		t.convertFun(s.Def.(*FunDef), t.box.funs[s.Aux])
	}
}

func (t *Typer) convertFun(def *FunDef, f *ast.Fun) {
	t.sig = &def.Signature
	t.convertBlock(&def.Body, &f.Body)
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
func (t *Typer) convertNode(scope *Scope, s ast.Statement) (stm Statement, exit bool) {
	switch s := s.(type) {
	case *ast.Return:
		return t.convertReturn(s), true
	// case ast.Var:
	// 	return t.translateVar(s)
	// case ast.Const:
	// 	return t.translateConst(s)
	// case ast.If:
	// 	return t.translateIf(s)
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
		block.Scope.Init(scok.Block, scope)
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

func (t *Typer) convertReturn(r *ast.Return) Statement {
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
		// TODO: assign expression here
	}
	return &Return{
		Pin: r.Pin,
		Exp: exp,
	}

	// exp, err := t.scope.TranslateExp(&stg.Hint{}, r.Exp)
	// if err != nil {
	// 	return nil, err
	// }
	// err = t.checkResultExp(t.sig.Result, exp)
	// if err != nil {
	// 	return nil, err
	// }

	// return &stg.Ret{Exp: exp}, nil
}
