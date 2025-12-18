package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) translateSymbols() diag.Error {
	for _, f := range t.funs {
		err := t.translateSymbolBody(f)
		if err != nil {
			return err
		}
	}
	for _, m := range t.methods {
		err := t.translateSymbolBody(m)
		if err != nil {
			return err
		}
	}
	return nil
}

// get function or method body
func (t *Typer) getBody(s *stg.Symbol) ast.Block {
	switch s.Kind {
	case smk.Fun:
		return t.box.Fun(s.Aux).Body
	case smk.Method:
		return t.box.Method(s.Aux).Body
	default:
		panic(fmt.Sprintf("unexpected %s (=%d) symbol", s.Kind, s.Kind))
	}
}

func (t *Typer) translateSymbolBody(s *stg.Symbol) diag.Error {
	def := s.Def.(*stg.Fun)
	t.sig = &def.Signature

	err := t.translateBlock(&def.Body, t.getBody(s))
	if err != nil {
		return err
	}
	return nil
}

func (t *Typer) translateBlock(block *stg.Block, b ast.Block) diag.Error {
	if len(b.Nodes) == 0 {
		return nil
	}

	oldScope := t.scope
	t.scope = &block.Scope

	nodes := make([]stg.Statement, 0, len(b.Nodes))
	for i, n := range b.Nodes {
		s, err := t.translateStatement(n)
		if err != nil {
			return err
		}
		if s == nil {
			// skip empty statements
			continue
		}

		_, ok := s.(*stg.Ret)
		if ok && i != len(b.Nodes)-1 {
			return &diag.SimpleMessageError{
				Pin:  n.Span().Pin,
				Text: "dead code after return",
			}
		}

		nodes = append(nodes, s)
	}

	t.scope = oldScope
	if len(nodes) == 0 {
		return nil
	}
	block.Nodes = nodes
	return nil
}

func (t *Typer) translateStatement(stm ast.Statement) (stg.Statement, diag.Error) {
	switch s := stm.(type) {
	case ast.Ret:
		return t.translateRet(s)
	case ast.Var:
		return t.translateVar(s)
	case ast.If:
		return t.translateIf(s)
	case ast.Assign:
		return t.translateAssign(s)
	case ast.Block:
		if len(s.Nodes) == 0 {
			// block statement with no statements is equivalent to empty statement
			return nil, nil
		}
		var block stg.Block
		block.Scope.Init(sck.Block, t.scope)
		err := t.translateBlock(&block, s)
		if err != nil {
			return nil, err
		}
		if len(block.Nodes) == 0 {
			// non empty AST block can still result in empty block
			return nil, nil
		}
		return &block, nil
	default:
		panic(fmt.Sprintf("unexpected %s (=%d) statement (%T)", s.Kind(), s.Kind(), s))
	}
}

func (t *Typer) translateAssign(a ast.Assign) (*stg.Assign, diag.Error) {
	switch r := a.Target.(type) {
	case ast.Symbol:
		return t.translateAssignSymbol(r, a)
	default:
		panic(fmt.Sprintf("unexpected %s (=%d) target expression (%T)", r.Kind(), r.Kind(), r))
	}
}

func (t *Typer) translateAssignSymbol(symbol ast.Symbol, a ast.Assign) (*stg.Assign, diag.Error) {
	name := symbol.Name
	pin := symbol.Pin

	s := t.scope.Lookup(name)
	if s == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", name),
		}
	}
	if s.Kind != smk.Var && s.Kind != smk.Param {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("cannot assign to %s symbol \"%s\"", s.Kind, name),
		}
	}

	exp, err := t.translateExp(a.Value)
	if err != nil {
		return nil, err
	}
	err = t.checkAssignTypes(s.Type, exp)
	if err != nil {
		return nil, err
	}

	return &stg.Assign{
		Symbol: s,
		Exp:    exp,
	}, nil
}

func (t *Typer) checkAssignTypes(want *stg.Type, exp stg.Exp) diag.Error {
	if exp.Type() == want {
		return nil
	}

	return nil
}

func (t *Typer) translateVar(v ast.Var) (*stg.Var, diag.Error) {
	typ, err := t.ctx.Types.Lookup(t.scope, v.Type)
	if err != nil {
		return nil, err
	}

	name := v.Name.Str
	pin := v.Name.Pin

	if t.scope.Has(name) {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("symbol \"%s\" was already declared in this block", name),
		}
	}
	s := t.scope.Alloc(smk.Var, name, pin)
	s.Type = typ

	return &stg.Var{
		Symbol: s,
	}, nil
}

func (t *Typer) translateIf(f ast.If) (*stg.If, diag.Error) {
	var s stg.If

	err := t.translateBranch(&s.If, f.If)
	if err != nil {
		return nil, err
	}

	if len(f.ElseIfs) != 0 {
		s.ElseIfs = make([]stg.Branch, 0, len(f.ElseIfs))
		for i := range len(f.ElseIfs) {
			err := t.translateBranch(&s.ElseIfs[i], f.ElseIfs[i])
			if err != nil {
				return nil, err
			}
		}
	}

	if f.Else != nil && len(f.Else.Nodes) != 0 {
		var block stg.Block
		block.Scope.Init(sck.Branch, t.scope)
		err := t.translateBlock(&block, *f.Else)
		if err != nil {
			return nil, err
		}
		if len(block.Nodes) != 0 {
			s.Else = &block
		}
	}
	return &s, nil
}

func (t *Typer) translateBranch(branch *stg.Branch, f ast.IfClause) diag.Error {
	exp, err := t.translateExp(f.Exp)
	if err != nil {
		return err
	}
	typ := exp.Type()
	if typ.Kind != tpk.Boolean {
		return &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("branch condition yields %s value, not a boolean", typ),
		}
	}

	branch.Block.Scope.Init(sck.Branch, t.scope)
	err = t.translateBlock(&branch.Block, f.Body)
	if err != nil {
		return err
	}

	branch.Exp = exp
	return nil
}

func (t *Typer) translateRet(r ast.Ret) (*stg.Ret, diag.Error) {
	if t.sig.Never {
		return nil, &diag.SimpleMessageError{
			Pin:  r.Pin,
			Text: "return in function with never result",
		}
	}
	if t.sig.Result == nil {
		if r.Exp == nil {
			return &stg.Ret{}, nil
		}
		return nil, &diag.SimpleMessageError{
			Pin:  r.Pin,
			Text: "return expression in function with void result",
		}
	}
	if r.Exp == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  r.Pin,
			Text: "return with no expression in function that must return something",
		}
	}

	exp, err := t.translateExp(r.Exp)
	if err != nil {
		return nil, err
	}
	err = t.checkResultExp(t.sig.Result, exp)
	if err != nil {
		return nil, err
	}

	return &stg.Ret{Exp: exp}, nil
}

func (t *Typer) translateExp(exp ast.Exp) (stg.Exp, diag.Error) {
	return t.ctx.Types.MakeInteger(exp.Span().Pin, 0), nil
}

// check that function result and expression types are compatible
func (t *Typer) checkResultExp(want *stg.Type, exp stg.Exp) diag.Error {
	if exp.Type() == want {
		return nil
	}

	return nil
}
