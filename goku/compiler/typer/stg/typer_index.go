package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

func (t *Typer) indexImports() {
	for _, s := range t.unit.Imports {
		err := t.indexImport(s)
		if err != nil {
			t.report(err)
		}
	}
}

func (t *Typer) indexConsts() {
	for i, c := range t.box.consts {
		err := t.indexConst(i, c)
		if err != nil {
			t.report(err)
		}
	}
}

func (t *Typer) indexTypes() {
	for i, typ := range t.box.types {
		err := t.indexType(i, typ)
		if err != nil {
			t.report(err)
		}
	}
}

func (t *Typer) indexVars() {
	for i, v := range t.box.vars {
		err := t.indexVar(i, v)
		if err != nil {
			t.report(err)
		}
	}
}

func (t *Typer) indexFuns() {
	for i, f := range t.box.funs {
		err := t.indexFun(i, f)
		if err != nil {
			t.report(err)
		}
	}
}

func (t *Typer) indexMethods() {
	for i, m := range t.box.methods {
		err := t.indexMethod(i, m)
		if err != nil {
			t.report(err)
		}
	}
}

func (t *Typer) indexImport(s sm.ImportSite) diag.Error {
	unit := t.com.Map[s.Path]
	if unit == nil {
		panic(fmt.Sprintf("unit \"%s\" not found: impossible due to map construction", s.Path))
	}

	name := s.Name
	pin := s.Pin
	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		return newMultDefError(name, pin, symbol.Pin)
	}

	symbol = t.unit.Scope.New(smk.Import, name, pin)
	symbol.Def = &SymDefUnit{Unit: unit}
	return nil
}

func (t *Typer) indexConst(i int, c ast.TopConst) diag.Error {
	name := c.Name.Str
	pin := c.Name.Pin

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		return newMultDefError(name, pin, symbol.Pin)
	}

	symbol = t.unit.Scope.New(smk.Const, name, pin)
	symbol.Aux = uint32(i)
	if c.Pub {
		symbol.Flags |= SymbolPublic
	}
	t.consts = append(t.consts, symbol)
	return nil
}

func (t *Typer) indexType(i int, typ ast.Type) diag.Error {
	name := typ.Name.Str
	pin := typ.Name.Pin

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		return newMultDefError(name, pin, symbol.Pin)
	}

	symbol = t.unit.Scope.New(smk.Type, name, pin)
	symbol.Aux = uint32(i)
	if typ.Pub {
		symbol.Flags |= SymbolPublic
	}
	t.types = append(t.types, symbol)
	return nil
}

func (t *Typer) indexVar(i int, v ast.TopVar) diag.Error {
	name := v.Name.Str
	pin := v.Name.Pin

	if v.Pub {
		return &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("variable \"%s\" declared as public", name),
		}
	}

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		return newMultDefError(name, pin, symbol.Pin)
	}

	symbol = t.unit.Scope.New(smk.Var, name, pin)
	symbol.Aux = uint32(i)
	t.vars = append(t.vars, symbol)
	return nil
}

func (t *Typer) indexFun(i int, f ast.Fun) diag.Error {
	name := f.Name.Str
	pin := f.Name.Pin

	if f.Unsafe {
		// name = "unsafe." + name
		panic("not implemented")
	}
	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		return newMultDefError(name, pin, symbol.Pin)
	}

	symbol = t.unit.Scope.New(smk.Fun, name, pin)
	symbol.Aux = uint32(i)
	if f.Pub {
		symbol.Flags |= SymbolPublic
	}
	if f.Export {
		symbol.Flags |= SymbolExport
		t.unit.Export = append(t.unit.Export, symbol)
	}
	t.funs = append(t.funs, symbol)
	return nil
}

func (t *Typer) indexMethod(i int, m ast.Method) diag.Error {
	pin := m.Name.Pin
	name := m.Receiver.Name.Str + "." + m.Name.Str

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		return newMultDefError(name, pin, symbol.Pin)
	}

	symbol = t.unit.Scope.New(smk.Method, name, pin)
	symbol.Aux = uint32(i)
	if m.Pub {
		symbol.Flags |= SymbolPublic
	}

	name = m.Receiver.Name.Str
	pin = m.Receiver.Name.Pin

	r := t.unit.Scope.Get(name)
	if r == nil {
		return &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("method receiver \"%s\" refers to undefined symbol", name),
		}
	}
	if r.Kind != smk.Type {
		return &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("method receiver \"%s\" refers to %s symbol (instead of custom type)", name, r.Kind),
		}
	}

	// bind method to its receiver
	t.mr[r] = append(t.mr[r], symbol)
	t.methods = append(t.methods, symbol)
	return nil
}

func newMultDefError(name string, pin, prev sm.Pin) diag.Error {
	return &diag.SimpleMessageError{
		Pin:  pin,
		Text: fmt.Sprintf("multiple definitions of symbol \"%s\"", name),
	}
}
