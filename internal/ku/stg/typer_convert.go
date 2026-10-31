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
		t.report(pin, 0, "not all paths result in return from this function body")
	}
}

// provides context for converting statements and expressions to STG form
type context struct {
	// fallback location for reporting errors
	//
	// used when down the call chain it is not clear where error is
	// originating
	pin sx.Pin

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
			t.report(0, 0, "dead code after statement that terminates further execution")
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
			t.report(s.Pin, 0, fmt.Sprintf("symbol \"%s\" was declared but never used", s.Name))
			continue
		}
		if s.rnum == 0 {
			t.report(s.Pin, 0, fmt.Sprintf("no reads from symbol \"%s\"", s.Name))
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
	case *ast.Var:
		return t.convertVar(c, s), nodestat{}
	case *ast.IfBlock:
		return t.convertIfBlock(c, s)
	case *ast.Branch:
		return t.convertBranch(c, s)
	case *ast.IfFlow:
		return t.convertIfFlow(c, s)
	// case ast.Match:
	// 	return t.translateMatch(s)
	case *ast.AssignSymbol:
		return t.convertAssignSymbol(c, s)
	case *ast.AssignDiscard:
		return t.convertAssignDiscard(c, s)
	case *ast.Assign:
		return t.convertAssign(c, s)
	case *ast.InvokeSymbol:
		return t.convertInvokeSymbol(c, s)
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

func (t *Typer) convertVar(s *Scope, c *ast.Var) Statement {
	name := c.Name
	pin := c.Pin

	symbol := s.Get(name)
	if symbol != nil {
		t.report(pin, 0, fmt.Sprintf("symbol named \"%s\" was already declared in this block", name))
		return nil
	}

	var vt *Type // variable declared type
	if c.Type != nil {
		vt = t.LookupType(s, c.Type)
	}

	var exp Exp
	if c.Exp == nil {
		if vt == nil {
			panic("no variable type")
		}
		exp = &Zero{pin: pin, typ: vt}
	} else {
		exp = t.convertExp(&context{scope: s}, c.Exp)
		typ := exp.Type()

		if vt == nil {
			vt = typ
		} else {
			t.checkAssignType(vt, exp)
		}
	}

	symbol = s.New(symk.Var, name, pin)
	symbol.Type = vt

	return &InitVar{
		Symbol: symbol,
		Exp:    exp,
	}
}

func (t *Typer) convertCreate(s *Scope, c *ast.Create) Statement {
	name := c.Name
	pin := c.Pin

	symbol := s.Get(name)
	if symbol != nil {
		t.report(pin, 0, fmt.Sprintf("symbol named \"%s\" was already declared in this block", name))
		return nil
	}

	exp := t.convertExp(&context{scope: s}, c.Exp)
	typ := exp.Type()
	if typ.IsStaticInteger() {
		typ = t.common.Types.Known.Sint
	}

	symbol = s.New(symk.Car, name, pin)
	symbol.Type = typ

	return &InitVar{
		Symbol: symbol,
		Exp:    exp,
	}
}

func (t *Typer) convertConst(s *Scope, c *ast.Const) Statement {
	name := c.Name
	pin := c.Pin

	symbol := s.Get(name)
	if symbol != nil {
		t.report(pin, 0, fmt.Sprintf("symbol named \"%s\" was already declared in this block", name))
		return nil
	}

	var typ *Type
	if c.Type != nil {
		typ = t.LookupType(s, c.Type)
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

func (t *Typer) convertAssign(s *Scope, c *ast.Assign) (Statement, nodestat) {
	target := t.convertExp(&context{scope: s, wuse: true}, c.Target)
	if !target.Target() {
		t.report(target.Pin(), 0, "cannot assign to expression of this form")
	}
	exp := t.convertExp(&context{scope: s}, c.Exp)
	t.checkAssignType(target.Type(), exp)

	return &Assign{
		Target: target,
		Exp:    exp,
	}, nodestat{}
}

func (t *Typer) convertAssignSymbol(s *Scope, c *ast.AssignSymbol) (Statement, nodestat) {
	name := c.Name
	pin := c.Pin

	symbol := s.Lookup(name)
	if symbol == nil {
		t.report(pin, 0, fmt.Sprintf("unknown symbol \"%s\" used as assignment target", name))
		return nil, nodestat{}
	}
	if symbol.Kind != symk.Var {
		t.report(pin, 0, fmt.Sprintf("symbol \"%s\" is not a variable and thus cannot be target for assignment", name))
		return nil, nodestat{}
	}

	exp := t.convertExp(&context{scope: s}, c.Exp)
	t.checkAssignType(symbol.Type, exp)

	return &AssignVar{
		Symbol: symbol,
		Exp:    exp,
	}, nodestat{}
}

func (t *Typer) convertAssignDiscard(s *Scope, c *ast.AssignDiscard) (Statement, nodestat) {
	exp := t.convertExp(&context{scope: s}, c.Exp)
	if exp.Type() == nil {
		t.report(exp.Pin(), 0, "cannot discard expression which yields nothing")
	}

	return &DiscardExp{
		pin: c.Pin,
		Exp: exp,
	}, nodestat{}
}

func (t *Typer) convertWhile(s *Scope, while *ast.While) (Statement, nodestat) {
	exp := t.convertExp(&context{scope: s}, while.Exp)
	typ := exp.Type()
	if typ.Kind != typk.Invalid && typ.Kind != typk.Boolean {
		// invalid expression was already reported earlier during expression conversion
		// no sense reporting it twice
		t.report(exp.Pin(), 0, fmt.Sprintf("value of %s type used in branch condition (must be boolean type instead)", typ))
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

func (t *Typer) convertIfFlow(s *Scope, f *ast.IfFlow) (Statement, nodestat) {
	exp := t.convertExp(&context{scope: s}, f.Exp)
	typ := exp.Type()
	if !typ.IsInvalid() && typ.Kind != typk.Boolean {
		// invalid expression was already reported earlier during expression conversion
		// no sense reporting it twice
		t.report(exp.Pin(), 0, fmt.Sprintf("value of %s type used in branch condition (must be boolean type instead)", typ))
	}

	node, _ := t.convertNode(s, f.Then)

	switch n := node.(type) {
	case *Return:
		return &IfRet{
			Exp: exp,
			Ret: n.Exp,
			pin: f.Pin,
		}, nodestat{exits: 1, etyp: ExitBranch}
	default:
		t.report(n.Pin(), 0, "short if (without block) can only be used for jump-like statements")
		return nil, nodestat{}
	}
}

func (t *Typer) convertIfBlock(s *Scope, f *ast.IfBlock) (Statement, nodestat) {
	exp := t.convertExp(&context{scope: s}, f.Exp)
	typ := exp.Type()
	if typ.Kind != typk.Invalid && typ.Kind != typk.Boolean {
		// invalid expression was already reported earlier during expression conversion
		// no sense reporting it twice
		t.report(exp.Pin(), 0, fmt.Sprintf("value of %s type used in branch condition (must be boolean type instead)", typ))
	}

	m := IfBlock{Exp: exp, pin: f.Pin}
	m.Body.Scope.Init(scok.Branch, s)
	t.convertBlock(&m.Body, &f.Body)
	exits := m.Body.Exits

	var etyp ExitType
	switch m.Body.ExitType {
	case ExitAlways, ExitBranch:
		etyp = ExitBranch
	}
	return &m, nodestat{exits: exits, etyp: etyp}
}

func (t *Typer) convertBranch(s *Scope, f *ast.Branch) (Statement, nodestat) {
	exp := t.convertExp(&context{scope: s}, f.Exp)
	typ := exp.Type()
	if typ.Kind != typk.Invalid && typ.Kind != typk.Boolean {
		// invalid expression was already reported earlier during expression conversion
		// no sense reporting it twice
		t.report(exp.Pin(), 0, fmt.Sprintf("value of %s type used in branch condition (must be boolean type instead)", typ))
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
		t.report(pin, 0, "return in function which never returns according to declaration")
	} else {
		if t.sig.Result == nil {
			if r.Exp != nil {
				t.report(pin, 0, "non-empty return expression in function with declared void result")
			}
		} else {
			if r.Exp == nil {
				t.report(pin, 0, "return with no expression in function that must return something")
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
		return t.makeInteger(e.Pin, e.Val, uint16(e.Aux))
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
	case *ast.SliceAll:
		return t.convertSliceAll(c, e)
	case *ast.SymCall:
		return t.convertSymCall(c, e)
	case *ast.SymAddr:
		return t.convertSymAddr(c, e)
	case *ast.SymDeref:
		return t.convertSymDeref(c, e)
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

func (t *Typer) convertSymAddr(c *context, exp *ast.SymAddr) Exp {
	name := exp.Name
	pin := exp.Pin

	symbol := c.scope.Lookup(name)
	if symbol == nil {
		t.report(pin, 0, fmt.Sprintf("call to unknown symbol \"%s\"", name))
		return t.makeInvExp(pin)
	}
	switch symbol.Kind {
	case symk.Var, symk.Car, symk.Param:
		// continue execution
	default:
		t.report(pin, 0, fmt.Sprintf("cannot take address of \"%s\" which is a %s", name, symbol.Kind))
		return t.makeInvExp(pin)
	}

	typ := t.common.Types.getRef(symbol.Type)
	return &SymAddr{
		pin:    pin,
		typ:    typ,
		Symbol: symbol,
	}
}

func (t *Typer) convertSymDeref(c *context, exp *ast.SymDeref) Exp {
	name := exp.Name
	pin := exp.Pin

	symbol := c.scope.Lookup(name)
	if symbol == nil {
		t.report(pin, 0, fmt.Sprintf("call to unknown symbol \"%s\"", name))
		return t.makeInvExp(pin)
	}

	var typ *Type
	if symbol.Type.Kind != typk.Ref {
		t.report(pin, 0, fmt.Sprintf("symbol \"%s\" is not a reference type and thus cannot be dereferenced", name))
		typ = t.common.Types.Invalid
	} else {
		typ = symbol.Type.Def.(*Ref).Type
	}

	if c.wuse {
		symbol.wnum += 1
	}
	symbol.rnum += 1

	return &SymDeref{
		pin:    pin,
		typ:    typ,
		Symbol: symbol,
	}
}

func (t *Typer) convertInvokeSymbol(c *Scope, s *ast.InvokeSymbol) (Statement, nodestat) {
	name := s.Name
	pin := s.Pin

	symbol := c.Lookup(name)
	if symbol == nil {
		t.report(pin, 0, fmt.Sprintf("call to unknown symbol \"%s\"", name))
		return nil, nodestat{}
	}
	if symbol.Kind != symk.Fun {
		t.report(pin, 0, fmt.Sprintf("cannot call \"%s\" which is not a function", name))
		return nil, nodestat{}
	}

	def := symbol.Def.(*FunDef)
	if def.Result != nil {
		t.report(pin, 0, fmt.Sprintf("implicit discard of function \"%s\" call result", name))
	}

	args := t.convertCallArgs(&context{scope: c, pin: pin}, name, def.Inputs, s.Args)
	return &InvokeFun{
		pin:    pin,
		Args:   args,
		Symbol: symbol,
	}, nodestat{}
}

func (t *Typer) convertCallArgs(c *context, name string, inputs []*Type, list []ast.Exp) []Exp {
	n := len(list)
	args := make([]Exp, 0, len(list))
	for i := range len(list) {
		arg := t.convertExp(c, list[i])
		args = append(args, arg)
	}

	if len(list) < len(inputs) {
		t.report(c.pin, 0, fmt.Sprintf("not enough arguments to call \"%s\" (got %d, want %d)", name, len(list), len(list)))
		return nil
	}
	if len(list) > len(inputs) {
		t.report(c.pin, 0, fmt.Sprintf("too many arguments in call to \"%s\" (got %d, want %d)", name, len(list), len(list)))
		return nil
	}

	// true: n == len(list) == len(inputs)
	for i := range n {
		want := inputs[i]
		t.checkArgType(want, args[i])
	}

	return args
}

func (t *Typer) convertSymCall(c *context, exp *ast.SymCall) Exp {
	name := exp.Name
	pin := exp.Pin

	symbol := c.scope.Lookup(name)
	if symbol == nil {
		t.report(pin, 0, fmt.Sprintf("call to unknown symbol \"%s\"", name))
		return t.makeInvExp(pin)
	}
	if symbol.Kind != symk.Fun {
		t.report(pin, 0, fmt.Sprintf("cannot call \"%s\" which is not a function", name))
		return t.makeInvExp(pin)
	}

	inv := false // should produce invalid expression after typechecking errors
	def := symbol.Def.(*FunDef)
	if def.Result == nil {
		if def.Never {
			t.report(pin, 0, fmt.Sprintf("call to \"%s\" never returns and thus cannot be used as expression", name))
		} else {
			t.report(pin, 0, fmt.Sprintf("call to \"%s\" returns nothing and thus cannot be used as expression", name))
		}
		inv = true
	}

	args := t.convertCallArgs(&context{scope: c.scope, pin: pin}, name, def.Inputs, exp.Args)
	if inv {
		return t.makeInvExp(pin)
	}

	return &FunCall{
		pin:    pin,
		typ:    def.Result,
		Args:   args,
		Symbol: symbol,
	}
}

func (t *Typer) convertSymZeroExp(c *context, exp *ast.SymZeroExp) Exp {
	name := exp.Name
	pin := exp.Pin

	symbol := c.scope.Lookup(name)
	if symbol == nil {
		t.report(pin, 0, fmt.Sprintf("unknown symbol \"%s\" used as expression", name))
		return t.makeInvExp(pin)
	}

	switch symbol.Kind {
	case symk.Type:
		typ := symbol.Def.(*Type)
		return &Zero{pin: pin, typ: typ}
	default:
		panic(fmt.Sprintf("unexpected %s symbol", symbol.Kind))
	}
}

func (t *Typer) convertSymExp(c *context, exp *ast.SymExp) Exp {
	name := exp.Name
	pin := exp.Pin

	symbol := c.scope.Lookup(name)
	if symbol == nil {
		t.report(pin, 0, fmt.Sprintf("unknown symbol \"%s\" used as expression", name))
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
		symk.Var,
		symk.Car,
		// symk.Loop,
		symk.Param:

		if c.static {
			t.report(pin, 0, fmt.Sprintf("runtime value symbol \"%s\" used in compile-time expression", name))
			return t.makeInvExp(pin)
		}

		return &SymExp{
			pin:    pin,
			typ:    symbol.Type,
			Symbol: symbol,
		}
	default:
		t.report(pin, 0, fmt.Sprintf("symbol \"%s\" cannot be used as operand or expression", name))
		return t.makeInvExp(pin)
	}
}

func (t *Typer) convertSliceAll(c *context, s *ast.SliceAll) Exp {
	pin := s.Pin

	exp := t.convertChain(c, &s.Chain)
	typ := exp.Type()
	if typ.IsInvalid() {
		return t.makeInvExp(pin)
	}

	switch typ.Kind {
	case typk.Array:
		def := typ.Def.(*Array)
		return &SliceArray{
			pin:  pin,
			Exp:  exp,
			A:    t.makeInteger(pin, 0, 0),
			B:    t.makeInteger(pin, uint64(def.Num), 0),
			typ:  t.common.Types.getSpan(def.Type),
			safe: true,
		}
	case typk.Span:
		// slice all applied to span just gives the same span
		return exp
	case typk.String:
		// slice all applied to string just gives the same string
		return exp
	default:
		t.report(pin, 0, fmt.Sprintf("cannot slice %s expression", typ))
		return t.makeInvExp(pin)
	}
}

func (t *Typer) convertChain(c *context, chain *ast.Chain) Exp {
	name := chain.Name
	pin := chain.Pin

	symbol := c.scope.Lookup(name)
	if symbol == nil {
		t.report(pin, sx.UnknownSymbol, fmt.Sprintf("unknown symbol \"%s\" used as expression", name))
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
	case *ast.Index:
		return t.applyIndex(c, exp, p)
	default:
		panic(fmt.Sprintf("unexpected %T part", p))
	}
}

func (t *Typer) applyIndex(c *context, exp Exp, x *ast.Index) Exp {
	index := t.convertExp(c, x.Exp)
	xtyp := index.Type()

	if xtyp.IsInvalid() {
		index = t.makeInvExp(index.Pin())
	} else if xtyp.Kind != typk.Integer {
		t.report(index.Pin(), 0, fmt.Sprintf("cannot use %s expression as index", xtyp))
		index = t.makeInvExp(index.Pin())
	}

	if xtyp.IsStaticInteger() {
		n := exp.(*Integer)
		if n.Neg && n.Val != 0 {
			t.report(index.Pin(), 0, fmt.Sprintf("negative number (-%d) used as index", n.Val))
			index = t.makeInvExp(index.Pin())
		}
	}

	typ := exp.Type()
	switch typ.Kind {
	case typk.Span:
		return &IndexSpan{
			pin:   x.Pin,
			typ:   typ.Def.(*Span).Type,
			Exp:   exp,
			Index: index,
		}
	case typk.Array:
		return &IndexSpan{
			pin:   x.Pin,
			typ:   typ.Def.(*Array).Type,
			Exp:   exp,
			Index: index,
		}
	case typk.Invalid:
		return exp
	default:
		t.report(x.Pin, 0, fmt.Sprintf("cannot index into %s expression", typ))
		return t.makeInvExp(x.Pin)
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
			t.report(s.Pin, 0, fmt.Sprintf("span type %s does not have \"%s\" field", typ, s.Name))
			return &InvExp{
				pin: s.Pin,
				typ: t.common.Types.Invalid,
			}
		}
	case typk.Invalid:
		return exp
	default:
		t.report(s.Pin, 0, fmt.Sprintf("cannot use select operation on %s type", typ))
		return t.makeInvExp(s.Pin)
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
