package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
	"github.com/mebyus/ku/goku/compiler/sm"
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
	case ast.Invoke:
		return t.translateInvoke(s)
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

func (t *Typer) translateInvoke(v ast.Invoke) (*stg.InvokeSymbol, diag.Error) {
	name := v.Call.Chain.Start.Str
	pin := v.Call.Chain.Start.Pin

	s := t.scope.Lookup(name)
	if s == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", name),
		}
	}

	var args []stg.Exp
	if len(v.Call.Args) != 0 {
		args = make([]stg.Exp, 0, len(v.Call.Args))

		for _, arg := range v.Call.Args {
			a, err := t.translateExp(arg)
			if err != nil {
				return nil, err
			}
			args = append(args, a)
		}
	}

	switch s.Kind {
	case smk.Fun:
		if len(v.Call.Chain.Parts) != 0 {
			return nil, &diag.SimpleMessageError{
				Pin:  pin,
				Text: fmt.Sprintf("\"%s\" is a function and cannot be chained", name),
			}
		}

		err := stg.CheckCall(&s.Def.(*stg.Fun).Signature, args)
		if err != nil {
			err.SetFallbackSpan(sm.Span{Pin: pin})
			return nil, err
		}

		return &stg.InvokeSymbol{
			Symbol: s,
			Args:   args,
		}, nil
	case smk.Receiver:
		panic("not implemented")
	case smk.Param, smk.Var:
		panic("not implemented")
	case smk.Method:
		panic(fmt.Sprintf("unexpected %s (=%d) symbol \"%s\" at chain start", s.Kind, s.Kind, name))
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("%s symbol \"%s\" cannot start a chain", s.Kind, name),
		}
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
	err = stg.CheckAssign(s.Type, exp)
	if err != nil {
		return nil, err
	}

	return &stg.Assign{
		Symbol: s,
		Exp:    exp,
	}, nil
}

func (t *Typer) translateVar(v ast.Var) (*stg.Var, diag.Error) {
	typ, err := t.scope.LookupType(v.Type)
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

	var exp stg.Exp
	if v.Exp != nil {
		exp, err = t.scope.TranslateExp(v.Exp)
		if err != nil {
			return nil, err
		}
		err = stg.CheckAssign(typ, exp)
		if err != nil {
			return nil, err
		}
	}

	return &stg.Var{
		Symbol: s,
		Exp:    exp,
	}, nil
}

func (t *Typer) translateIf(f ast.If) (stg.Statement, diag.Error) {
	var s stg.If
	branches := make([]*stg.Branch, 0, len(f.ElseIfs)+1)

	b, err := t.translateBranch(f.If)
	if err != nil {
		return nil, err
	}
	branches = append(branches, b)

	for _, e := range f.ElseIfs {
		b, err := t.translateBranch(e)
		if err != nil {
			return nil, err
		}
		branches = append(branches, b)
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

	k := 0 // counts runtime branches
	for i := range len(branches) {
		b := branches[i]

		if b.IsStatic() {
			if !b.IsTrue() {
				// static false branch is nop
				continue
			}

			// here k is equal to number of runtime branches that has come before
			// this branch
			if k == 0 {
				if b.IsEmpty() {
					// empty static true branch with no runtime branches before
					// transforms the whole statement into nop
					return nil, nil
				}

				// static true branch with no runtime branches before
				// transforms the whole statement into block statement
				return &b.Block, nil
			}

			if b.IsEmpty() {
				// empty static true branch with at least one runtime branch before
				// ends remaining branches with no else block
				s.Else = nil
				break
			}

			// static true branch with at least one runtime branch before
			// replaces else branch
			s.Else = &b.Block
			break
		}

		k += 1
	}

	const debug = false
	if k != len(branches) {
		if debug {
			fmt.Printf("branch %d/%d alive\n", k, len(branches))
		}

		branches = branches[:k]
		clear(branches[k:])
	}

	if len(branches) == 0 {
		if s.Else == nil {
			// all branches were eliminated at compile-time, statement becomes nop
			return nil, nil
		}

		// only else branch survived, transform it into block statement
		return s.Else, nil
	}

	s.Branches = branches
	return &s, nil
}

func (t *Typer) translateBranch(f ast.IfClause) (*stg.Branch, diag.Error) {
	var b stg.Branch

	exp, err := t.translateExp(f.Exp)
	if err != nil {
		return nil, err
	}
	typ := exp.Type()
	if typ.Kind != tpk.Boolean {
		return nil, &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("branch condition yields %s value, not a boolean", typ),
		}
	}

	b.Block.Scope.Init(sck.Branch, t.scope)
	err = t.translateBlock(&b.Block, f.Body)
	if err != nil {
		return nil, err
	}

	b.Exp = exp
	b.SetFlags()
	return &b, nil
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
	return t.scope.TranslateExp(exp)
}

// check that function result and expression types are compatible
func (t *Typer) checkResultExp(want *stg.Type, exp stg.Exp) diag.Error {
	if exp.Type() == want {
		return nil
	}

	return nil
}
