package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
)

func (t *Typer) inspectStatement(stm ast.Statement) diag.Error {
	switch s := stm.(type) {
	case ast.Gonext, ast.Break, ast.Never, ast.Stub, ast.Panic:
		return nil
	case ast.Var:
		return t.inspectVarStatement(s)
	case ast.Const:
		return t.inspectConstStatement(s)
	case ast.If:
		return t.inspectIfStatement(s)
	case ast.Ret:
		return t.inspectRetStatement(s)
	case ast.Must:
		return t.inspectMustStatement(s)
	case ast.StaticMust:
		return t.inspectStaticMust(s)
	case ast.Assign:
		return t.inspectAssignStatement(s)
	case ast.Invoke:
		return t.inspectInvoke(s)
	case ast.ForRange:
		return t.inspectForRange(s)
	case ast.While:
		return t.inspectWhile(s)
	case ast.Block:
		return t.inspectBlock(s)
	case ast.Loop:
		return t.inspectLoop(s)
	case ast.Match:
		return t.inspectMatch(s)
	case ast.DeferCall:
		return t.inspectInvoke(ast.Invoke{Call: s.Call})
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) statement (%T)", s.Kind(), s.Kind(), s))
	}
}

func (t *Typer) inspectMatch(m ast.Match) diag.Error {
	err := t.inspectExp(m.Exp)
	if err != nil {
		return err
	}

	for _, c := range m.Cases {
		for _, e := range c.List {
			err := t.inspectExp(e)
			if err != nil {
				return err
			}
		}
		err := t.inspectBlock(c.Body)
		if err != nil {
			return err
		}
	}

	if m.Else == nil {
		return nil
	}

	return t.inspectBlock(*m.Else)
}

func (t *Typer) inspectLoop(l ast.Loop) diag.Error {
	return t.inspectBlock(l.Body)
}

func (t *Typer) inspectWhile(w ast.While) diag.Error {
	err := t.inspectExp(w.Exp)
	if err != nil {
		return err
	}

	return t.inspectBlock(w.Body)
}

func (t *Typer) inspectForRange(f ast.ForRange) diag.Error {
	if f.Type != nil {
		err := t.inspectVarType(f.Type)
		if err != nil {
			return err
		}
	}

	if f.Start != nil {
		err := t.inspectExp(f.Start)
		if err != nil {
			return err
		}
	}

	err := t.inspectExp(f.End)
	if err != nil {
		return err
	}

	return t.inspectBlock(f.Body)
}

func (t *Typer) inspectInvoke(v ast.Invoke) diag.Error {
	return t.inspectCallExp(v.Call)
}

func (t *Typer) inspectAssignStatement(a ast.Assign) diag.Error {
	err := t.inspectExp(a.Target)
	if err != nil {
		return err
	}
	return t.inspectExp(a.Value)
}

func (t *Typer) inspectConstStatement(c ast.Const) diag.Error {
	if c.Type != nil {
		err := t.inspectVarType(c.Type)
		if err != nil {
			return err
		}
	}

	return t.inspectVarInitExp(c.Exp)
}

func (t *Typer) inspectStaticMust(m ast.StaticMust) diag.Error {
	return t.inspectExp(m.Exp)
}

func (t *Typer) inspectMustStatement(m ast.Must) diag.Error {
	return t.inspectExp(m.Exp)
}

func (t *Typer) inspectRetStatement(r ast.Ret) diag.Error {
	if r.Exp == nil {
		return nil
	}

	return t.inspectExp(r.Exp)
}

func (t *Typer) inspectIfStatement(f ast.If) diag.Error {
	err := t.inspectExp(f.If.Exp)
	if err != nil {
		return err
	}
	err = t.inspectBlock(f.If.Body)
	if err != nil {
		return err
	}

	for _, ef := range f.ElseIfs {
		err = t.inspectExp(ef.Exp)
		if err != nil {
			return err
		}
		err = t.inspectBlock(ef.Body)
		if err != nil {
			return err
		}
	}

	if f.Else == nil {
		return nil
	}

	return t.inspectBlock(*f.Else)
}

func (t *Typer) inspectVarStatement(stm ast.Var) diag.Error {
	err := t.inspectVarType(stm.Type)
	if err != nil {
		return err
	}

	return t.inspectVarInitExp(stm.Exp)
}

func (t *Typer) inspectVarType(spec ast.TypeSpec) diag.Error {
	return t.inspectType(spec)
}

func (t *Typer) inspectVarInitExp(exp ast.Exp) diag.Error {
	if exp == nil {
		// variable can have empty init expression
		return nil
	}

	return t.inspectExp(exp)
}

func (t *Typer) inspectExp(exp ast.Exp) diag.Error {
	switch e := exp.(type) {
	case ast.Integer, ast.String, ast.Rune, ast.True, ast.False, ast.Nil, ast.Dirty, ast.ErrorId, ast.EnumMacro, ast.DotName:
		return nil
	case ast.Symbol:
		return t.linkExpSymbol(e)
	case ast.Unary:
		return t.inspectUnaryExp(e)
	case ast.Binary:
		return t.inspectBinaryExp(e)
	case ast.Paren:
		return t.inspectParenExp(e)
	case ast.Chain:
		return t.inspectChainExp(e)
	case ast.Cast:
		return t.inspectCastExp(e)
	case ast.Tint:
		return t.inspectTintExp(e)
	case ast.Call:
		return t.inspectCallExp(e)
	case ast.Slice:
		return t.inspectSliceExp(e)
	case ast.GetRef:
		return t.inspectRefExp(e)
	case ast.Object:
		return t.inspectObject(e)
	case ast.List:
		return t.inspectListExp(e)
	case ast.Pack:
		return t.inspectPack(e)
	case ast.Size:
		return t.inspectSizeExp(e)
	case ast.DerefSlice:
		return t.inspectDerefSlice(e)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}
}

func (t *Typer) linkExpSymbol(symbol ast.Symbol) diag.Error {
	name := symbol.Name
	s := t.unit.Scope.Get(name)
	if s == nil {
		// probably local or global symbol
		// we are interested only in unit level symbols here
		// if it does not exist then block scan phase will catch it
		return nil
	}

	if s.Kind == smk.Import || s.Kind == smk.Type {
		// TODO: this can be shadowed name, need to refactor inspection approach
		//
		// return &diag.SimpleMessageError{
		// 	Pin:  symbol.Pin,
		// 	Text: fmt.Sprintf("%s symbol \"%s\" used as operand in expression", s.Kind, name),
		// }
	}

	t.ins.link(s)
	return nil
}

func (t *Typer) inspectDerefSlice(s ast.DerefSlice) diag.Error {
	err := t.inspectChainExp(s.Chain)
	if err != nil {
		return err
	}

	if s.Start != nil {
		err = t.inspectExp(s.Start)
		if err != nil {
			return err
		}
	}

	if s.End != nil {
		err = t.inspectExp(s.End)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Typer) inspectPack(p ast.Pack) diag.Error {
	for _, e := range p.List {
		err := t.inspectExp(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) inspectListExp(l ast.List) diag.Error {
	for _, e := range l.Exps {
		err := t.inspectExp(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) inspectSizeExp(s ast.Size) diag.Error {
	return t.inspectType(s.Exp)
}

func (t *Typer) inspectTintExp(e ast.Tint) diag.Error {
	err := t.inspectVarType(e.Type)
	if err != nil {
		return err
	}
	return t.inspectExp(e.Exp)
}

func (t *Typer) inspectObject(o ast.Object) diag.Error {
	for _, f := range o.Fields {
		err := t.inspectExp(f.Exp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) inspectRefExp(r ast.GetRef) diag.Error {
	return t.inspectChainExp(r.Chain)
}

func (t *Typer) inspectParenExp(p ast.Paren) diag.Error {
	return t.inspectExp(p.Exp)
}

func (t *Typer) inspectUnaryExp(u ast.Unary) diag.Error {
	return t.inspectExp(u.Exp)
}

func (t *Typer) inspectSliceExp(s ast.Slice) diag.Error {
	err := t.inspectChainExp(s.Chain)
	if err != nil {
		return err
	}

	if s.Start != nil {
		err = t.inspectExp(s.Start)
		if err != nil {
			return err
		}
	}

	if s.End != nil {
		err = t.inspectExp(s.End)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Typer) inspectCallExp(c ast.Call) diag.Error {
	err := t.inspectChainExp(c.Chain)
	if err != nil {
		return err
	}

	for _, arg := range c.Args {
		err = t.inspectExp(arg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) inspectCastExp(c ast.Cast) diag.Error {
	err := t.inspectVarType(c.Type)
	if err != nil {
		return err
	}
	return t.inspectExp(c.Exp)
}

func (t *Typer) inspectBinaryExp(exp ast.Binary) diag.Error {
	err := t.inspectExp(exp.A)
	if err != nil {
		return err
	}
	return t.inspectExp(exp.B)
}

func (t *Typer) inspectChainExp(chain ast.Chain) diag.Error {
	name := chain.Start.Str
	if name == "" {
		// handle unsafe prefix
		if len(chain.Parts) == 0 {
			panic("empty chain")
		}
		p := chain.Parts[0]
		u, ok := p.(ast.Unsafe)
		if !ok {
			// TODO: should redesign unsafe storage in AST
			fmt.Printf("WARN: probably an error in parser")
			return nil
		}
		name = "unsafe." + u.Name
		s := t.unit.Scope.Get(name)
		if s == nil {
			return &diag.SimpleMessageError{
				Pin:  u.Pin,
				Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", name),
			}
		}

		t.ins.link(s)
		return nil
	}

	s := t.unit.Scope.Get(name)
	if s == nil {
		// probably local or global symbol
		return nil
	}

	t.ins.link(s)
	return nil
}
