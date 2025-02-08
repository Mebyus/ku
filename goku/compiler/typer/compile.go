package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/source"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

// Typer is a high-level algorithm driver that gathers multiple ASTs of unit's source
// texts to produce that unit's STG.
//
// Unit typing is done in several separate phases:
//
//   - 1 (gather) - gather all AST texts of the unit
//   - 2 (index) - index unit level symbols
//   - 3 (method bind) - bind methods to corresponding receivers
//   - 4 (inspect) - determine dependency relations between unit level symbols
//   - 5 (hoist) - construct, map and rank (hoist) symbol dependency graph
//   - 6 (static eval) - eval and finalize all properties of unit level types and constants
//   - 7 (block scan) - recursively scan statements and expressions inside functions
type Typer struct {
	box Box

	Warns []diag.Error

	unit *stg.Unit

	ctx *stg.Context
}

func Compile(c *stg.Context, unit *stg.Unit, texts []*ast.Text) diag.Error {
	if c == nil {
		panic("nil context")
	}
	if unit == nil {
		panic("nil unit")
	}
	if len(texts) == 0 {
		panic("no texts")
	}

	unit.InitScopes(&c.Global)
	t := &Typer{
		ctx:  c,
		unit: unit,
	}
	t.box.init(texts)
	return t.compile(texts)
}

func (t *Typer) compile(texts []*ast.Text) diag.Error {
	err := t.addTexts(texts)
	if err != nil {
		return err
	}
	return nil
}

func (t *Typer) addTexts(texts []*ast.Text) diag.Error {
	for _, text := range texts {
		err := t.addText(text)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) addText(text *ast.Text) diag.Error {
	err := t.addImports(t.unit.Imports)
	if err != nil {
		return err
	}
	err = t.addTypes(text.Types)
	if err != nil {
		return err
	}
	err = t.addConstants(text.Constants)
	if err != nil {
		return err
	}
	err = t.addFuns(text.Functions)
	if err != nil {
		return err
	}
	err = t.addVars(text.Variables)
	if err != nil {
		return err
	}
	err = t.addMethods(text.Methods)
	if err != nil {
		return err
	}
	err = t.addTests(text.Tests)
	if err != nil {
		return err
	}

	return nil
}

func (t *Typer) addImports(imports []stg.ImportSite) diag.Error {
	for _, s := range imports {
		err := t.addImport(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) addTypes(types []ast.Type) diag.Error {
	for _, typ := range types {
		err := t.addType(typ)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) addConstants(constants []ast.TopLet) diag.Error {
	for _, c := range constants {
		err := t.addConst(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) addFuns(funs []ast.Fun) diag.Error {
	for _, fun := range funs {
		err := t.addFun(fun)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) addVars(vars []ast.TopVar) diag.Error {
	for _, v := range vars {
		err := t.addVar(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) addMethods(methods []ast.Method) diag.Error {
	for _, method := range methods {
		err := t.addMethod(method)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) addTests(tests []ast.Fun) diag.Error {
	for _, test := range tests {
		err := t.addTest(test)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Typer) addImport(s stg.ImportSite) diag.Error {
	unit := t.ctx.Map[s.Path]
	if unit == nil {
		panic("unit not found: impossible due to map construction")
	}

	name := s.Name
	pin := s.Pin
	if t.unit.Scope.Has(name) {
		return errMultDef(name, pin)
	}

	_ = t.unit.Scope.Alloc(smk.Import, name, pin)
	return nil
}

func (t *Typer) addFun(fun ast.Fun) diag.Error {
	name := fun.Name.Str
	pin := fun.Name.Pin

	if t.unit.Scope.Has(name) {
		return errMultDef(name, pin)
	}

	symbol := t.unit.Scope.Alloc(smk.Fun, name, pin)
	symbol.Aux = t.box.addFun(fun)
	return nil
}

func (t *Typer) addType(typ ast.Type) diag.Error {
	name := typ.Name.Str
	pin := typ.Name.Pin

	if t.unit.Scope.Has(name) {
		return errMultDef(name, pin)
	}

	symbol := t.unit.Scope.Alloc(smk.Type, name, pin)
	symbol.Aux = t.box.addType(typ)
	return nil
}

func (t *Typer) addConst(c ast.TopLet) diag.Error {
	name := c.Name.Str
	pin := c.Name.Pin

	if t.unit.Scope.Has(name) {
		return errMultDef(name, pin)
	}

	symbol := t.unit.Scope.Alloc(smk.Let, name, pin)
	symbol.Aux = t.box.addConst(c)
	return nil
}

func (t *Typer) addVar(v ast.TopVar) diag.Error {
	name := v.Name.Str
	pin := v.Name.Pin

	if t.unit.Scope.Has(name) {
		return errMultDef(name, pin)
	}

	symbol := t.unit.Scope.Alloc(smk.Var, name, pin)
	symbol.Aux = t.box.addVar(v)
	return nil
}

func (t *Typer) addMethod(m ast.Method) diag.Error {
	pin := m.Name.Pin
	name := m.Receiver.Name.Str + "." + m.Name.Str

	if t.unit.Scope.Has(name) {
		return errMultDef(name, pin)
	}

	symbol := t.unit.Scope.Alloc(smk.Method, name, pin)
	symbol.Aux = t.box.addMethod(m)
	return nil
}

func (t *Typer) addTest(test ast.Fun) diag.Error {
	name := test.Name.Str
	pin := test.Name.Pin

	if t.unit.TestScope.Has(name) {
		return errMultDef(name, pin)
	}

	symbol := t.unit.TestScope.Alloc(smk.Test, name, pin)
	symbol.Aux = t.box.addTest(test)
	return nil
}

func errMultDef(name string, pin source.Pin) diag.Error {
	return &diag.SimpleMessageError{
		Pin:  pin,
		Text: fmt.Sprintf("multiple definitions of symbol \"%s\"", name),
	}
}
