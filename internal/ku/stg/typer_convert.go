package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/scok"
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/enums/typk"
	"github.com/mebyus/ku/internal/ku/sx"
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

	if def.Signature.Result != nil && def.Body.ExitType != ExitAlways {
		var pin sx.Pin
		if len(def.Body.Nodes) == 0 {
			pin = def.Body.Pin()
		} else {
			pin = def.Body.Nodes[len(def.Body.Nodes)-1].Pin()
		}
		t.report(pin, "not all paths result in return from this function body")
	}
}

// provides context for converting statements and expressions to STG form
type context struct {
	scope *Scope

	// if true means that expression must be evaluated at compile-time
	// and cannot contain runtime values inside
	static bool

	// true when expression is used for writing to a variable/field
	wuse bool
}

// various information about statement after it was converted and analyzed
type nodestat struct {
	// total number of exits inside this node
	exits uint32

	etyp ExitType
}

func (t *Typer) convertBlock(block *Block, b *ast.Block) {
	block.pin = b.Pin

	if len(b.Nodes) == 0 {
		return
	}

	// stats of this block
	var exits uint32
	var etyp ExitType

	nodes := make([]Statement, 0, len(b.Nodes))
	for i, n := range b.Nodes {
		s, stat := t.convertNode(&block.Scope, n)
		if s == nil {
			// skip empty statements
			continue
		}

		exits += stat.exits
		if stat.etyp == ExitAlways && i != len(b.Nodes)-1 {
			// TODO: make this a warning
			// TODO: extract pin from statement here
			t.report(0, "dead code after statement that terminates further execution")
		}

		if stat.etyp > etyp {
			// in order for this analysis hack to work
			// exit types must be ordered: never < branch < always
			etyp = stat.etyp
		}

		nodes = append(nodes, s)
	}
	t.checkSymbolUsage(block)

	if len(nodes) == 0 {
		// discard allocated nodes memory
		return
	}

	block.Nodes = nodes
	block.Exits = exits
	block.ExitType = etyp
}

func (t *Typer) checkSymbolUsage(block *Block) {
	for _, s := range block.Scope.Symbols {
		switch s.Kind {
		case symk.Param:
			continue
		}

		if s.rnum+s.wnum == 0 {
			t.report(s.Pin, fmt.Sprintf("symbol \"%s\" was declared but never used", s.Name))
			continue
		}
		if s.rnum == 0 {
			t.report(s.Pin, fmt.Sprintf("no reads from symbol \"%s\"", s.Name))
		}
	}
}

// returned statement can be nil in case of error of if it was
// optimized out (for example it is the case for empty block statement)
//
// if exit return value equals true it means that the statement exits
// function execution (or it never stops like endless loop)
func (t *Typer) convertNode(c *Scope, s ast.Statement) (Statement, nodestat) {
	switch s := s.(type) {
	case *ast.Return:
		return t.convertReturn(c, s), nodestat{exits: 1, etyp: ExitAlways}
	// case ast.Var:
	// 	return t.translateVar(s)
	case *ast.Const:
		return t.convertConst(c, s), nodestat{}
	case *ast.Create:
		return t.convertCreate(c, s), nodestat{}
	case *ast.Branch:
		return t.convertBranch(c, s)
	case *ast.LineIf:
		return t.convertLineIf(c, s)
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
	case *ast.While:
		return t.convertWhile(c, s)
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
			return nil, nodestat{}
		}
		var block Block
		block.Scope.Init(scok.Block, c)
		t.convertBlock(&block, s)
		if len(block.Nodes) == 0 {
			// non empty AST block can still result in empty block
			return nil, nodestat{}
		}

		return &block, nodestat{exits: block.Exits, etyp: block.ExitType}
	default:
		panic(fmt.Sprintf("unexpected %T statement", s))
	}
}

func (t *Typer) convertCreate(s *Scope, c *ast.Create) Statement {
	name := c.Name
	pin := c.Pin

	symbol := t.unit.Scope.Get(name)
	if symbol != nil {
		t.report(pin, fmt.Sprintf("symbol named \"%s\" was already declared in this block", name))
		return nil
	}

	exp := t.convertExp(&context{scope: s}, c.Exp)
	typ := exp.Type()

	symbol = s.New(symk.Fixed, name, pin)
	symbol.Type = typ

	return &CreateFixed{
		Symbol: symbol,
		Exp:    exp,
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

func (t *Typer) convertWhile(s *Scope, while *ast.While) (Statement, nodestat) {
	exp := t.convertExp(&context{scope: s}, while.Exp)
	typ := exp.Type()
	if typ.Kind != typk.Invalid && typ.Kind != typk.Boolean {
		// invalid expression was already reported earlier during expression conversion
		// no sense reporting it twice
		t.report(exp.Pin(), fmt.Sprintf("value of %s type used in branch condition (must be boolean type instead)", typ))
	}

	w := While{Exp: exp, pin: while.Pin}
	w.Body.Scope.Init(scok.Loop, s)
	t.convertBlock(&w.Body, &while.Body)
	exits := w.Body.Exits

	var etyp ExitType
	switch w.Body.ExitType {
	case ExitAlways, ExitBranch:
		etyp = ExitBranch
	}

	return &w, nodestat{exits: exits, etyp: etyp}
}

func (t *Typer) convertLineIf(s *Scope, f *ast.LineIf) (Statement, nodestat) {
	exp := t.convertExp(&context{scope: s}, f.Exp)
	typ := exp.Type()
	if typ.Kind != typk.Invalid && typ.Kind != typk.Boolean {
		// invalid expression was already reported earlier during expression conversion
		// no sense reporting it twice
		t.report(exp.Pin(), fmt.Sprintf("value of %s type used in branch condition (must be boolean type instead)", typ))
	}

	m := LineIf{Exp: exp, pin: f.Pin}
	m.Then, _ = t.convertNode(s, f.Then)
	return &m, nodestat{exits: 1, etyp: ExitBranch}
}

func (t *Typer) convertBranch(s *Scope, f *ast.Branch) (Statement, nodestat) {
	exp := t.convertExp(&context{scope: s}, f.Exp)
	typ := exp.Type()
	if typ.Kind != typk.Invalid && typ.Kind != typk.Boolean {
		// invalid expression was already reported earlier during expression conversion
		// no sense reporting it twice
		t.report(exp.Pin(), fmt.Sprintf("value of %s type used in branch condition (must be boolean type instead)", typ))
	}

	m := Branch{Exp: exp, pin: f.Pin}
	m.Body.Scope.Init(scok.Branch, s)
	t.convertBlock(&m.Body, &f.Body)
	exits := m.Body.Exits

	if f.Else == nil || len(f.Else.Nodes) == 0 {
		var etyp ExitType
		switch m.Body.ExitType {
		case ExitAlways, ExitBranch:
			etyp = ExitBranch
		}
		return &m, nodestat{exits: exits, etyp: etyp}
	}

	var block Block
	block.Scope.Init(scok.Branch, s)
	t.convertBlock(&block, f.Else)
	if len(block.Nodes) != 0 {
		exits += block.Exits
		m.Else = &block
	}

	var etyp ExitType
	if m.Body.ExitType == ExitAlways && block.ExitType == ExitAlways {
		etyp = ExitAlways
	} else if m.Body.ExitType == ExitNever && block.ExitType == ExitNever {
		etyp = ExitNever
	} else {
		etyp = ExitBranch
	}

	return &m, nodestat{exits: exits, etyp: etyp}
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
		pin: pin,
		Exp: exp,
	}
}

func (t *Typer) convertExp(c *context, exp ast.Exp) Exp {
	switch e := exp.(type) {
	// case ast.Nil:
	// return s.Types.MakeNil(e.Pin), nil
	case *ast.Integer:
		return t.makeInteger(e.Pin, e.Val)
	case *ast.String:
		return t.makeString(e.Pin, e.Val)
	case *ast.True:
		return t.makeBoolean(e.Pin, true)
	case *ast.False:
		return t.makeBoolean(e.Pin, false)
	case *ast.SymExp:
		return t.convertSymExp(c, e)
	case *ast.SymZeroExp:
		return t.convertSymZeroExp(c, e)
	case *ast.Chain:
		return t.convertChain(c, e)
	case *ast.BinExp:
		return t.convertBinExp(c, e)
	case *ast.ParenExp:
		return t.convertExp(c, e.Exp)
	case *ast.InvExp:
		return t.makeInvExp(e.Pin)
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}

func (t *Typer) convertSymZeroExp(c *context, exp *ast.SymZeroExp) Exp {
	name := exp.Name
	pin := exp.Pin

	symbol := c.scope.Lookup(name)
	if symbol == nil {
		t.report(pin, fmt.Sprintf("unknown symbol \"%s\" used as expression", name))
		return t.makeInvExp(pin)
	}

	switch symbol.Kind {
	case symk.Type:
		typ := symbol.Def.(*Type)
		return &ZeroValue{pin: pin, typ: typ}
	default:
		panic(fmt.Sprintf("unexpected %s symbol", symbol.Kind))
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

	if c.wuse {
		symbol.wnum += 1
	} else {
		symbol.rnum += 1
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

func (t *Typer) convertChain(c *context, chain *ast.Chain) Exp {
	name := chain.Name
	pin := chain.Pin

	symbol := c.scope.Lookup(name)
	if symbol == nil {
		t.report(pin, fmt.Sprintf("unknown symbol \"%s\" used as expression", name))
		return t.makeInvExp(pin)
	}

	if c.wuse {
		symbol.wnum += 1
	} else {
		symbol.rnum += 1
	}

	var exp Exp
	exp = &SymExp{
		pin:    pin,
		typ:    symbol.Type,
		Symbol: symbol,
	}
	for _, p := range chain.Parts {
		exp = t.applyChainPart(c, exp, p)
	}
	return exp
}

func (t *Typer) applyChainPart(c *context, exp Exp, part ast.Part) Exp {
	switch p := part.(type) {
	case *ast.Select:
		return t.applySelect(c, exp, p)
	default:
		panic(fmt.Sprintf("unexpected %T part", p))
	}
}

func (t *Typer) applySelect(c *context, exp Exp, s *ast.Select) Exp {
	typ := exp.Type()
	switch typ.Kind {
	case typk.Span:
		switch s.Name {
		case "num":
			return &SpanNum{
				pin: s.Pin,
				Exp: exp,
				typ: t.common.Types.Known.Uint,
			}
		case "ptr":
			panic("stub")
		default:
			t.report(s.Pin, fmt.Sprintf("span type %s does not have \"%s\" field", typ, s.Name))
			return &InvExp{
				pin: s.Pin,
				typ: t.common.Types.Invalid,
			}
		}
	case typk.Invalid:
		return exp
	default:
		t.report(s.Pin, fmt.Sprintf("cannot use select operation on %s type", typ))
		return &InvExp{
			pin: s.Pin,
			typ: t.common.Types.Invalid,
		}
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
