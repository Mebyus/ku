package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/aok"
	"github.com/mebyus/ku/goku/compiler/enums/exk"
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
	case ast.Const:
		return t.translateConst(s)
	case ast.If:
		return t.translateIf(s)
	case ast.Match:
		return t.translateMatch(s)
	case ast.Assign:
		return t.translateAssign(s)
	case ast.Invoke:
		return t.translateInvoke(s)
	case ast.DeferCall:
		return t.translateDeferCall(s)
	case ast.Loop:
		return t.translateLoop(s)
	case ast.While:
		return t.translateWhile(s)
	case ast.ForRange:
		return t.translateForRange(s)
	case ast.Must:
		return t.translateMust(s)
	case ast.Panic:
		return t.translatePanic(s)
	case ast.Stub:
		return &stg.Stub{Pin: s.Pin}, nil
	case ast.Never:
		return &stg.Never{Pin: s.Pin}, nil
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

func (t *Typer) translatePanic(p ast.Panic) (stg.Statement, diag.Error) {
	exp, err := t.scope.TranslateExp(&stg.Hint{}, p.Exp)
	if err != nil {
		return nil, err
	}

	typ := exp.Type()
	if typ.Kind != tpk.String {
		panic(fmt.Sprintf("not implemented for %s type", typ))
	}

	return &stg.Panic{
		Exp: exp,
		Pin: p.Pin,
	}, nil
}

func (t *Typer) translateAssign(a ast.Assign) (stg.Statement, diag.Error) {
	k := a.Op.Kind
	switch k {
	case aok.Simple:
		return t.translateSimpleAssign(a)
	case aok.Walrus:
		return t.translateWalrusAssign(a)
	case aok.Add, aok.Sub, aok.Mul, aok.Div, aok.And, aok.Or, aok.Rem, aok.LeftShift, aok.RightShift:
		return t.translateOpAssign(a)
	default:
		panic(fmt.Sprintf("unexpected %s (=%d) assign operator", k, k))
	}
}

func (t *Typer) translateSimpleAssign(a ast.Assign) (stg.Statement, diag.Error) {
	target, err := t.scope.TranslateExp(&stg.Hint{}, a.Target)
	if err != nil {
		return nil, err
	}

	var hint stg.Hint
	typ := target.Type()
	if typ.Kind == tpk.Custom && typ.Def.(*stg.Custom).Type.Kind == tpk.Enum {
		hint.Enum = typ.Def.(*stg.Custom).Type.Def.(*stg.Enum)
	}

	exp, err := t.scope.TranslateExp(&hint, a.Value)
	if err != nil {
		return nil, err
	}
	err = stg.CheckAssign(target.Type(), exp)
	if err != nil {
		return nil, err
	}

	// TODO: probably need to move this check to stg.CheckAssign
	e, ok := target.(*stg.SymExp)
	if ok {
		s := e.Symbol
		if s.Kind != smk.Var && s.Kind != smk.Param {
			return nil, &diag.SimpleMessageError{
				Pin:  e.Pin,
				Text: fmt.Sprintf("cannot assign to %s symbol \"%s\"", s.Kind, s.Name),
			}
		}
	}

	return &stg.Assign{
		Target: target,
		Exp:    exp,
	}, nil
}

func (t *Typer) translateOpAssign(a ast.Assign) (stg.Statement, diag.Error) {
	_, ok := a.Target.(ast.Pack)
	if ok {
		return nil, &diag.SimpleMessageError{
			Pin:  a.Target.Span().Pin,
			Text: fmt.Sprintf("cannot use multiple assignment with %s operator", a.Op.Kind),
		}
	}

	target, err := t.scope.TranslateExp(&stg.Hint{}, a.Target)
	if err != nil {
		return nil, err
	}

	typ := target.Type()
	if typ.Kind != tpk.Integer {
		return nil, &diag.SimpleMessageError{
			Pin:  a.Target.Span().Pin,
			Text: fmt.Sprintf("cannot use assignment %s operator on %s type", a.Op.Kind, typ),
		}
	}

	exp, err := t.scope.TranslateExp(&stg.Hint{}, a.Value)
	if err != nil {
		return nil, err
	}
	err = stg.CheckAssign(target.Type(), exp)
	if err != nil {
		return nil, err
	}

	// TODO: need separate statement type for this
	return nil, nil
}

func (t *Typer) translateWalrusAssign(a ast.Assign) (stg.Statement, diag.Error) {
	exp, err := t.scope.TranslateExp(&stg.Hint{}, a.Value)
	if err != nil {
		return nil, err
	}

	if exp.Type().Kind == tpk.Tuple {
		if a.Target.Kind() != exk.Pack {
			return nil, &diag.SimpleMessageError{
				Pin:  a.Op.Pin,
				Text: "assignment of multiple values to a single target",
			}
		}
		return t.translatePackAssignOrDefine(a.Target.(ast.Pack).List, exp)
	}

	if a.Target.Kind() == exk.Pack {
		return nil, &diag.SimpleMessageError{
			Pin:  a.Op.Pin,
			Text: "assignment of single value to multiple targets",
		}
	}

	target, err := t.defineOrAssign(a.Target, exp.Type())
	if err != nil {
		return nil, err
	}

	return &stg.Assign{
		Target: target,
		Exp:    exp,
	}, nil
}

func (t *Typer) translatePackAssignOrDefine(targets []ast.Exp, exp stg.Exp) (stg.Statement, diag.Error) {
	tuple := exp.Type().Def.(stg.Tuple)
	types := tuple.Types

	if len(targets) != len(types) {
		return nil, &diag.SimpleMessageError{
			Pin:  targets[0].Span().Pin,
			Text: fmt.Sprintf("mismatched number of assign targets (%d) and values (%d)", len(targets), len(types)),
		}
	}

	list := make([]stg.Exp, 0, len(targets))
	for i := range len(targets) {
		target, err := t.defineOrAssign(targets[i], types[i])
		if err != nil {
			return nil, err
		}
		list = append(list, target)
	}

	return &stg.Assign{
		Exp:    exp,
		Target: t.ctx.Types.MakePack(list),
	}, nil
}

func (t *Typer) defineOrAssign(target ast.Exp, typ *stg.Type) (stg.Exp, diag.Error) {
	switch e := target.(type) {
	case ast.Symbol:
		name := e.Name
		pin := e.Pin

		symbol := t.scope.Lookup(name)
		if symbol == nil {
			symbol = t.scope.Alloc(smk.Var, name, pin)
			symbol.Type = typ
		}

		return &stg.SymExp{
			Pin:    pin,
			Symbol: symbol,
		}, nil
	default:
		exp, err := t.scope.TranslateExp(&stg.Hint{}, e)
		if err != nil {
			return exp, nil
		}

		// TODO: typecheck
		// err = stg.CheckAssign(exp.Type(), )
		return exp, nil
	}
}

func (t *Typer) translateDeferCall(c ast.DeferCall) (*stg.DeferCall, diag.Error) {
	call, err := t.scope.TranslateCall(&stg.Hint{}, c.Call)
	if err != nil {
		return nil, err
	}

	return &stg.DeferCall{Call: call}, nil

}

func (t *Typer) translateInvoke(v ast.Invoke) (*stg.Invoke, diag.Error) {
	call, err := t.scope.TranslateCall(&stg.Hint{}, v.Call)
	if err != nil {
		return nil, err
	}

	return &stg.Invoke{Call: call}, nil
}

func (t *Typer) translateMust(m ast.Must) (stg.Statement, diag.Error) {
	exp, err := t.scope.TranslateExp(&stg.Hint{}, m.Exp)
	if err != nil {
		return nil, err
	}

	typ := exp.Type()
	if typ.Kind != tpk.Boolean {
		return nil, &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("assert condition yields %s value, not a boolean", typ),
		}
	}

	if !typ.IsStatic() {
		return &stg.Must{Exp: exp}, nil
	}

	if exp.(*stg.Boolean).Val {
		// assert is true at compile time, transform to empty statement
		return nil, nil
	}

	return nil, &diag.SimpleMessageError{
		Pin:  exp.Span().Pin,
		Text: "compile-time assert is false",
	}
}

func (t *Typer) translateForRange(r ast.ForRange) (stg.Statement, diag.Error) {
	var f stg.ForRange
	var start stg.Exp
	var err diag.Error

	if r.Start != nil {
		start, err = t.scope.TranslateExp(&stg.Hint{}, r.Start)
		if err != nil {
			return nil, err
		}
		typ := start.Type()
		if typ.Kind != tpk.Integer {
			return nil, &diag.SimpleMessageError{
				Pin:  start.Span().Pin,
				Text: fmt.Sprintf("range start expression has %s value, not integer", typ),
			}
		}
	}

	end, err := t.scope.TranslateExp(&stg.Hint{}, r.End)
	if err != nil {
		return nil, err
	}
	typ := end.Type()
	if typ.Kind != tpk.Integer {
		return nil, &diag.SimpleMessageError{
			Pin:  end.Span().Pin,
			Text: fmt.Sprintf("range end expression has %s value, not integer", typ),
		}
	}

	var vt *stg.Type
	if r.Type != nil {
		vt, err = t.scope.LookupType(r.Type)
		if err != nil {
			return nil, err
		}
		if vt.Kind != tpk.Integer {
			return nil, &diag.SimpleMessageError{
				Pin:  r.Type.Span().Pin,
				Text: fmt.Sprintf("loop variable type %s, is not an integer", vt),
			}
		}
	} else {
		vt = t.ctx.Types.Known.Uint
	}

	f.Body.Scope.Init(sck.Loop, t.scope)
	name := r.Name.Str
	symbol := f.Body.Scope.Alloc(smk.Loop, name, r.Name.Pin)
	symbol.Type = vt

	err = t.translateBlock(&f.Body, r.Body)
	if err != nil {
		return nil, err
	}
	if len(f.Body.Nodes) == 0 {
		// transform for range loop with empty body into empty statement
		return nil, nil
	}

	f.Start = start
	f.End = end
	f.Var = symbol
	return &f, nil
}

func (t *Typer) translateLoop(l ast.Loop) (stg.Statement, diag.Error) {
	var loop stg.Loop
	loop.Body.Scope.Init(sck.Loop, t.scope)
	err := t.translateBlock(&loop.Body, l.Body)
	if err != nil {
		return nil, err
	}
	if len(loop.Body.Nodes) == 0 {
		return nil, &diag.SimpleMessageError{
			Pin:  l.Body.Pin,
			Text: "loop has empty body, spinloops are forbidden",
		}
	}

	return &loop, nil
}

func (t *Typer) translateWhile(w ast.While) (stg.Statement, diag.Error) {
	exp, err := t.scope.TranslateExp(&stg.Hint{}, w.Exp)
	if err != nil {
		return nil, err
	}

	typ := exp.Type()
	if typ.Kind != tpk.Boolean {
		return nil, &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("loop condition yields %s value, not a boolean", typ),
		}
	}

	if typ.IsStatic() {
		if exp.(*stg.Boolean).Val {
			// condition is always true, transform to Loop statement
			var loop stg.Loop
			loop.Body.Scope.Init(sck.Loop, t.scope)
			err = t.translateBlock(&loop.Body, w.Body)
			if err != nil {
				return nil, err
			}
			if len(loop.Body.Nodes) == 0 {
				return nil, &diag.SimpleMessageError{
					Pin:  w.Body.Pin,
					Text: "loop has empty body and its condition is always true, spinloops are forbidden",
				}
			}

			return &loop, nil
		} else {
			// condition is always false, check loop body and transform to empty statement
			var body stg.Block
			body.Scope.Init(sck.Loop, t.scope)
			err = t.translateBlock(&body, w.Body)
			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	}

	while := stg.While{Exp: exp}
	while.Body.Scope.Init(sck.Loop, t.scope)
	err = t.translateBlock(&while.Body, w.Body)
	if err != nil {
		return nil, err
	}

	return &while, nil
}

func (t *Typer) translateConst(c ast.Const) (stg.Statement, diag.Error) {
	exp, err := t.scope.TranslateExp(&stg.Hint{}, c.Exp)
	if err != nil {
		return nil, err
	}

	var typ *stg.Type
	if c.Type != nil {
		typ, err = t.scope.LookupType(c.Type)
		if err != nil {
			return nil, err
		}
		err = stg.CheckAssign(typ, exp)
		if err != nil {
			return nil, err
		}
	} else {
		typ = exp.Type()
	}

	name := c.Name.Str
	pin := c.Name.Pin

	if t.scope.Has(name) {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("symbol \"%s\" was already declared in this block", name),
		}
	}
	s := t.scope.Alloc(smk.Const, name, pin)
	s.Type = typ
	s.Def = stg.StaticValue{Exp: exp}

	// Constant definition is compile-time only construct and thus
	// translated to empty statement.
	return nil, nil
}

func (t *Typer) translateVar(v ast.Var) (*stg.Var, diag.Error) {
	typ, err := t.scope.LookupType(v.Type)
	if err != nil {
		return nil, err
	}

	var exp stg.Exp
	if v.Exp != nil {
		exp, err = t.scope.TranslateExp(&stg.Hint{}, v.Exp)
		if err != nil {
			return nil, err
		}
		err = stg.CheckAssign(typ, exp)
		if err != nil {
			return nil, err
		}
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
		Exp:    exp,
	}, nil
}

func (t *Typer) translateMatch(m ast.Match) (stg.Statement, diag.Error) {
	exp, err := t.scope.TranslateExp(&stg.Hint{}, m.Exp)
	if err != nil {
		return nil, err
	}

	typ := exp.Type()
	if typ.Kind == tpk.Integer || (typ.Kind == tpk.Custom && typ.Def.(*stg.Custom).Type.Kind == tpk.Enum) {
		return t.translateMatchInteger(exp, m)
	}

	panic(fmt.Sprintf("not implemented for %s type", typ))
}

func (t *Typer) translateMatchInteger(exp stg.Exp, m ast.Match) (stg.Statement, diag.Error) {
	var cases []*stg.MatchCase
	if len(m.Cases) != 0 {
		typ := exp.Type()
		cases = make([]*stg.MatchCase, 0, len(m.Cases))
		for _, mc := range m.Cases {
			c, err := t.translateMatchCase(typ, mc)
			if err != nil {
				return nil, err
			}
			cases = append(cases, c)
		}
	}

	var elseBlock *stg.Block
	if m.Else != nil && len(m.Else.Nodes) != 0 {
		var block stg.Block
		block.Scope.Init(sck.Case, t.scope)
		err := t.translateBlock(&block, *m.Else)
		if err != nil {
			return nil, err
		}
		if len(block.Nodes) != 0 {
			elseBlock = &block
		}
	}

	return &stg.MatchInteger{
		Exp:   exp,
		Cases: cases,
		Else:  elseBlock,
	}, nil
}

func (t *Typer) translateMatchCase(want *stg.Type, mc ast.MatchCase) (*stg.MatchCase, diag.Error) {
	var c stg.MatchCase
	c.List = make([]stg.Exp, 0, len(mc.List))
	for _, exp := range mc.List {
		e, err := t.scope.TranslateExp(&stg.Hint{}, exp)
		if err != nil {
			return nil, err
		}

		err = stg.CheckAssign(want, e)
		if err != nil {
			return nil, err
		}

		c.List = append(c.List, e)
	}

	c.Body.Scope.Init(sck.Case, t.scope)
	err := t.translateBlock(&c.Body, mc.Body)
	if err != nil {
		return nil, err
	}

	return &c, nil
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

	exp, err := t.scope.TranslateExp(&stg.Hint{}, f.Exp)
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

	exp, err := t.scope.TranslateExp(&stg.Hint{}, r.Exp)
	if err != nil {
		return nil, err
	}
	err = t.checkResultExp(t.sig.Result, exp)
	if err != nil {
		return nil, err
	}

	return &stg.Ret{Exp: exp}, nil
}

// check that function result and expression types are compatible
func (t *Typer) checkResultExp(want *stg.Type, exp stg.Exp) diag.Error {
	if exp.Type() == want {
		return nil
	}

	return nil
}
