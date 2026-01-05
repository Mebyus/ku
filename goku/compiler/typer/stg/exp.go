package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/bok"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
	"github.com/mebyus/ku/goku/compiler/enums/uok"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Exp node that represents an arbitrary expression.
type Exp interface {
	Type() *Type

	Span() sm.Span

	// Use only for debugging.
	String() string
}

func (s *Scope) TranslateExp(hint *Hint, exp ast.Exp) (Exp, diag.Error) {
	switch e := exp.(type) {
	case ast.Nil:
		return s.Types.MakeNil(e.Pin), nil
	case ast.Integer:
		return s.Types.MakeInteger(e.Pin, e.Val), nil
	case ast.String:
		return s.Types.MakeString(e.Pin, e.Val), nil
	case ast.True:
		return s.Types.MakeBoolean(e.Pin, true), nil
	case ast.False:
		return s.Types.MakeBoolean(e.Pin, false), nil
	case ast.Rune:
		return s.Types.MakeRune(e.Pin, uint32(e.Val)), nil
	case ast.Symbol:
		return s.translateSymbolExp(e)
	case ast.DotName:
		return s.translateDotName(hint, e)
	case ast.Paren:
		return s.TranslateExp(hint, e.Exp)
	case ast.Unary:
		return s.translateUnaryExp(hint, e)
	case ast.Binary:
		return s.translateBinaryExp(hint, e)
	case ast.Chain:
		return s.TranslateChain(hint, e)
	case ast.Call:
		return s.TranslateCall(hint, e)
	case ast.Slice:
		return s.translateSlice(hint, e)
	case ast.Pack:
		return s.translatePackExp(hint, e)
	case ast.Cast:
		return s.translateCast(hint, e)
	case ast.DerefSlice:
		return s.translateDerefSlice(hint, e)
	default:
		panic(fmt.Sprintf("unexpected %s (=%d) expression (%T)", e.Kind(), e.Kind(), e))
	}
}

func (s *Scope) translateUnaryExp(hint *Hint, u ast.Unary) (Exp, diag.Error) {
	exp, err := s.TranslateExp(hint, u.Exp)
	if err != nil {
		return nil, err
	}

	typ := exp.Type()
	if typ.IsStatic() {
		return s.evalConstUnaryExp(exp, u.Op)
	}

	k := u.Op.Kind
	switch k {
	case uok.Plus:
		if typ.Kind != tpk.Integer {
			return nil, &diag.SimpleMessageError{
				Pin:  u.Op.Pin,
				Text: fmt.Sprintf("type %s does not have + (unary plus) operation", typ),
			}
		}
		return exp, nil
	case uok.Minus:
		if typ.Kind != tpk.Integer {
			return nil, &diag.SimpleMessageError{
				Pin:  u.Op.Pin,
				Text: fmt.Sprintf("type %s does not have - (unary minus) operation", typ),
			}
		}
		return &Unary{
			Exp: exp,
			Op:  u.Op,
		}, nil
	case uok.Not:
		if typ.Kind != tpk.Boolean {
			return nil, &diag.SimpleMessageError{
				Pin:  u.Op.Pin,
				Text: fmt.Sprintf("type %s does not have ! (unary not) operation", typ),
			}
		}
		return &Unary{
			Exp: exp,
			Op:  u.Op,
		}, nil
	default:
		panic(fmt.Sprintf("unexpected %s (=%d) unary operator", k, k))
	}
}

func (s *Scope) translateDerefSlice(hint *Hint, d ast.DerefSlice) (Exp, diag.Error) {
	if d.End == nil {
		return s.translateSliceArrayRef(hint, d.Chain, d.Start)
	}

	return s.translateMakeSpan(hint, d)
}

func (s *Scope) translateMakeSpan(hint *Hint, d ast.DerefSlice) (Exp, diag.Error) {
	c, err := s.TranslateChain(hint, d.Chain)
	if err != nil {
		return nil, err
	}

	var elem *Type
	typ := c.Type()
	switch typ.Kind {
	case tpk.ArrayRef:
		elem = typ.Def.(ArrayRef).Type
	case tpk.ArrayPointer:
		elem = typ.Def.(ArrayPointer).Type
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  c.Span().Pin,
			Text: fmt.Sprintf("this operation is not allowed on %s value, must be array ref or array pointer", typ),
		}
	}

	var start Exp
	if d.Start != nil {
		start, err := s.TranslateExp(hint, d.Start)
		if err != nil {
			return nil, err
		}
		t := start.Type()
		if t.Kind != tpk.Integer {
			return nil, &diag.SimpleMessageError{
				Pin:  start.Span().Pin,
				Text: fmt.Sprintf("cannot use %s value as index, must be integer", t),
			}
		}
		// TODO: add static negative integer check
	}

	end, err := s.TranslateExp(hint, d.End)
	if err != nil {
		return nil, err
	}
	t := end.Type()
	if t.Kind != tpk.Integer {
		return nil, &diag.SimpleMessageError{
			Pin:  end.Span().Pin,
			Text: fmt.Sprintf("cannot use %s value as index, must be integer", t),
		}
	}

	return &MakeSpan{
		Start: start,
		End:   end,
		typ:   s.Types.getSpan(elem),
	}, nil
}

func (s *Scope) translateSliceArrayRef(hint *Hint, chain ast.Chain, start ast.Exp) (Exp, diag.Error) {
	c, err := s.TranslateChain(hint, chain)
	if err != nil {
		return nil, err
	}

	typ := c.Type()
	if typ.Kind != tpk.ArrayRef {
		return nil, &diag.SimpleMessageError{
			Pin:  start.Span().Pin,
			Text: fmt.Sprintf("this operation is not allowed on %s value, must be array ref", typ),
		}
	}

	if start == nil {
		return c, nil
	}

	index, err := s.TranslateExp(hint, start)
	if err != nil {
		return nil, err
	}

	t := index.Type()
	if t.Kind != tpk.Integer {
		return nil, &diag.SimpleMessageError{
			Pin:  index.Span().Pin,
			Text: fmt.Sprintf("cannot use %s value as index, must be integer", t),
		}
	}
	if t.IsStatic() {
		v := index.(*Integer)
		if v.Neg {
			return nil, &diag.SimpleMessageError{
				Pin:  v.Pin,
				Text: fmt.Sprintf("negative index value -%d", v.Val),
			}
		}
		if v.Val == 0 {
			return c, nil
		}
	}

	return &SliceArrayRef{
		Exp:   c,
		Index: index,
	}, nil
}

func (s *Scope) translateCast(hint *Hint, c ast.Cast) (Exp, diag.Error) {
	exp, err := s.TranslateExp(hint, c.Exp)
	if err != nil {
		return nil, err
	}

	want, err := s.LookupType(c.Type)
	if err != nil {
		return nil, err
	}

	if exp.Type() == want {
		// cast is not needed, simplify expression
		return exp, nil
	}
	err = s.Types.CheckCast(want, exp)
	if err != nil {
		return nil, err
	}

	return &Cast{
		Exp: exp,
		Pin: c.Type.Span().Pin,
		typ: want,
	}, nil
}

func (s *Scope) translatePackExp(hint *Hint, exp ast.Pack) (*Pack, diag.Error) {
	list := make([]Exp, 0, len(exp.List))
	for _, e := range exp.List {
		x, err := s.TranslateExp(hint, e)
		if err != nil {
			return nil, err
		}
		list = append(list, x)
	}

	return s.Types.MakePack(list), nil
}

func (s *Scope) translateSlice(hint *Hint, slice ast.Slice) (Exp, diag.Error) {
	exp, err := s.TranslateChain(hint, slice.Chain)
	if err != nil {
		return nil, err
	}

	var start Exp
	if slice.Start != nil {
		start, err = s.TranslateExp(hint, slice.Start)
		if err != nil {
			return nil, err
		}

		t := start.Type()
		if t.Kind != tpk.Integer {
			return nil, &diag.SimpleMessageError{
				Pin:  start.Span().Pin,
				Text: fmt.Sprintf("cannot use %s value as index, must be integer", t),
			}
		}
		if t.IsStatic() {
			v := start.(*Integer)
			if v.Neg {
				return nil, &diag.SimpleMessageError{
					Pin:  v.Pin,
					Text: fmt.Sprintf("negative index value -%d", v.Val),
				}
			}
		}
	}

	var end Exp
	if slice.End != nil {
		end, err = s.TranslateExp(hint, slice.End)
		if err != nil {
			return nil, err
		}

		t := end.Type()
		if t.Kind != tpk.Integer {
			return nil, &diag.SimpleMessageError{
				Pin:  end.Span().Pin,
				Text: fmt.Sprintf("cannot use %s value as index, must be integer", t),
			}
		}
		if t.IsStatic() {
			v := end.(*Integer)
			if v.Neg {
				return nil, &diag.SimpleMessageError{
					Pin:  v.Pin,
					Text: fmt.Sprintf("negative index value -%d", v.Val),
				}
			}
		}
	}

	typ := exp.Type()
	switch typ.Kind {
	case tpk.Array:
		panic("not implemented")
	case tpk.Span:
		return &SpanSlice{
			Exp:   exp,
			Start: start,
			End:   end,
		}, nil
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("type %s does not support slice expression", typ),
		}
	}
}

func (s *Scope) TranslateCall(hint *Hint, exp ast.Call) (Exp, diag.Error) {
	chain, err := s.TranslateChain(hint, exp.Chain)
	if err != nil {
		return nil, err
	}

	var call *Call
	var sig *Signature
	var args []Exp
	switch c := chain.(type) {
	case *SymExp:
		symbol := c.Symbol
		args = make([]Exp, 0, len(exp.Args))

		switch symbol.Kind {
		case smk.Fun:
			call = &Call{
				Pin:    exp.Span().Pin,
				Symbol: symbol,
			}
			sig = &symbol.Def.(*Fun).Signature
		default:
			panic(fmt.Sprintf("unexpected %s (=%d) symbol \"%s\"", symbol.Kind, symbol.Kind, symbol.Name))
		}
	case *BoundMethod:
		symbol := c.Symbol
		fun := symbol.Def.(*Fun)
		args = make([]Exp, 0, len(exp.Args))
		args = append(args, c.Receiver)

		call = &Call{
			Pin:    exp.Span().Pin,
			Symbol: symbol,
		}
		sig = &fun.Signature
	default:
		panic(fmt.Sprintf("unexpected (%T) chain expression in call", c))
	}

	if len(exp.Args) != len(sig.Params) {
		return nil, &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("call requires %d argument(s), but got %d", len(sig.Params), len(exp.Args)),
		}
	}

	for i, arg := range exp.Args {
		var h *Hint
		typ := sig.Params[i]
		if typ.Kind == tpk.Custom && typ.Def.(*Custom).Type.Kind == tpk.Enum {
			h = &Hint{Enum: typ.Def.(*Custom).Type.Def.(*Enum)}
		} else {
			h = hint
		}

		a, err := s.TranslateExp(h, arg)
		if err != nil {
			return nil, err
		}
		args = append(args, a)
	}

	err = CheckCall(sig, args)
	if err != nil {
		err.SetFallbackSpan(chain.Span())
		return nil, err
	}

	call.Args = args
	return call, nil
}

func (s *Scope) TranslateChain(hint *Hint, exp ast.Chain) (Exp, diag.Error) {
	name := exp.Start.Str
	pin := exp.Start.Pin

	start := s.Lookup(name)
	if start == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", name),
		}
	}

	var e Exp
	switch start.Kind {
	case smk.Receiver, smk.Var, smk.Param:
		e = &SymExp{
			Pin:    pin,
			Symbol: start,
		}
	case smk.Fun:
		if len(exp.Parts) != 0 {
			return nil, &diag.SimpleMessageError{
				Pin:  pin,
				Text: fmt.Sprintf("function \"%s\" cannot start a chain", name),
			}
		}

		e = &SymExp{
			Pin:    pin,
			Symbol: start,
		}
	case smk.Import:
		if len(exp.Parts) != 1 {
			return nil, &diag.SimpleMessageError{
				Pin:  pin,
				Text: fmt.Sprintf("multiple chain parts after import symbol \"%s\"", name),
			}
		}

		p, ok := exp.Parts[0].(ast.Select)
		if !ok {
			return nil, &diag.SimpleMessageError{
				Pin:  pin,
				Text: fmt.Sprintf("bad operation (%T) on import symbol \"%s\"", exp.Parts[0], name),
			}
		}

		unit := start.Def.(SymDefUnit).Unit
		sname := p.Name.Str
		symbol := unit.Scope.Lookup(sname)
		if symbol == nil {
			return nil, &diag.SimpleMessageError{
				Pin:  p.Name.Pin,
				Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", sname),
			}
		}
		if symbol.Kind != smk.Fun {
			return nil, &diag.SimpleMessageError{
				Pin:  p.Name.Pin,
				Text: fmt.Sprintf("name \"%s\" refers to %s, not a function", sname, symbol.Kind),
			}
		}
		if !symbol.IsPublic() {
			return nil, &diag.SimpleMessageError{
				Pin:  p.Name.Pin,
				Text: fmt.Sprintf("symbol \"%s\" is not public", sname),
			}
		}

		return &SymExp{
			Symbol: symbol,
			Pin:    p.Name.Pin,
		}, nil
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("%s symbol \"%s\" cannot start a chain", start.Kind, name),
		}
	}

	for _, p := range exp.Parts {
		var err diag.Error
		e, err = s.applyChainPart(hint, e, p)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (s *Scope) applyChainPart(hint *Hint, exp Exp, part ast.Part) (Exp, diag.Error) {
	switch p := part.(type) {
	case ast.Select:
		return s.applySelectPart(exp, p)
	case ast.DerefIndex:
		return s.applyDerefIndexPart(hint, exp, p)
	case ast.Index:
		return s.applyIndexPart(hint, exp, p)
	case ast.BagSelect:
		return s.applyBagSelectPart(exp, p)
	default:
		panic(fmt.Sprintf("unexpected %s (=%d) chain part (%T)", p.Kind(), p.Kind(), p))
	}
}

func (s *Scope) applyBagSelectPart(exp Exp, part ast.BagSelect) (Exp, diag.Error) {
	typ := exp.Type()

	switch typ.Kind {
	case tpk.Custom:
		if typ.Def.(*Custom).Type.Kind == tpk.Union {
			return s.applySelectFieldToUnion(exp, part)
		}
		return nil, &diag.SimpleMessageError{
			Pin:  part.Name.Pin,
			Text: fmt.Sprintf("cannot apply bag select to %s type which is a %s", typ, typ.Def.(*Custom).Type.Kind),
		}
	// case tpk.Pointer, tpk.Ref:
	// 	return s.applySelectToPointer(exp, part)
	default:
		panic(fmt.Sprintf("unexpected %s type", typ))
	}

}

func (s *Scope) applySelectFieldToUnion(exp Exp, part ast.BagSelect) (Exp, diag.Error) {
	name := part.Name.Str
	pin := part.Name.Pin

	typ := exp.Type()
	c := typ.Def.(*Custom).Type.Def.(*Union)
	f := c.getField(name)
	if f == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("type %s does not have field %s", typ, name),
		}
	}

	return &SelectField{
		Exp:   exp,
		Pin:   pin,
		Field: f,
	}, nil
}

func (s *Scope) applyIndexPart(hint *Hint, exp Exp, part ast.Index) (Exp, diag.Error) {
	index, err := s.TranslateExp(hint, part.Exp)
	if err != nil {
		return nil, err
	}

	xtyp := index.Type()
	if xtyp.Kind != tpk.Integer {
		return nil, &diag.SimpleMessageError{
			Pin:  index.Span().Pin,
			Text: fmt.Sprintf("index must have integer type, got %s", xtyp),
		}
	}

	if xtyp.IsStatic() {
		n := index.(*Integer)
		if n.Neg {
			return nil, &diag.SimpleMessageError{
				Pin:  n.Pin,
				Text: fmt.Sprintf("negative index value -%d", n.Val),
			}
		}
	}

	typ := exp.Type()
	switch typ.Kind {
	case tpk.Array:
		panic("not implemented")
	case tpk.Span:
		return &SpanIndex{
			Exp:   exp,
			Index: index,
			typ:   typ.Def.(Span).Type,
		}, nil
	case tpk.CapBuf:
		panic("not implemented")
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  index.Span().Pin,
			Text: fmt.Sprintf("type %s does not support index expression", typ),
		}
	}
}

func (s *Scope) applyDerefIndexPart(hint *Hint, exp Exp, part ast.DerefIndex) (*DerefIndex, diag.Error) {
	index, err := s.TranslateExp(hint, part.Exp)
	if err != nil {
		return nil, err
	}

	xtyp := index.Type()
	if xtyp.Kind != tpk.Integer {
		return nil, &diag.SimpleMessageError{
			Pin:  index.Span().Pin,
			Text: fmt.Sprintf("index must have integer type, got %s", xtyp),
		}
	}

	if xtyp.IsStatic() {
		n := index.(*Integer)
		if n.Neg {
			return nil, &diag.SimpleMessageError{
				Pin:  n.Pin,
				Text: fmt.Sprintf("negative index value -%d", n.Val),
			}
		}
	}

	typ := exp.Type()
	var rtyp *Type
	switch typ.Kind {
	case tpk.ArrayPointer:
		return nil, &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("indexed expression has array pointer %s type, which could have nil value", typ),
		}
	case tpk.ArrayRef:
		rtyp = typ.Def.(ArrayRef).Type
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("cannot use deref index on %s type", typ),
		}
	}

	return &DerefIndex{
		Exp:   exp,
		Index: index,
		typ:   rtyp,
	}, nil
}

func (s *Scope) applySelectPart(exp Exp, part ast.Select) (Exp, diag.Error) {
	typ := exp.Type()

	switch typ.Kind {
	case tpk.Span:
		return s.applySelectToSpan(exp, part)
	case tpk.Pointer, tpk.Ref:
		return s.applySelectToPointer(exp, part)
	case tpk.Custom:
		if typ.Def.(*Custom).Type.Kind == tpk.Struct {
			return s.applySelectToStruct(exp, part)
		}
		return nil, &diag.SimpleMessageError{
			Pin:  exp.Span().Pin,
			Text: fmt.Sprintf("cannot apply select part to %s type", typ),
		}
	default:
		panic(fmt.Sprintf("unexpected %s type", typ))
	}
}

func (s *Scope) applySelectToStruct(exp Exp, part ast.Select) (Exp, diag.Error) {
	name := part.Name.Str
	pin := part.Name.Pin

	typ := exp.Type()
	c := typ.Def.(*Custom)
	m := c.getMethod(name)
	if m != nil {
		panic("not implemented")
	}

	f := c.Type.Def.(*Struct).getField(name)
	if f == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("type %s does not have field or method named \"%s\"", typ, name),
		}
	}

	return &SelectField{
		Exp:   exp,
		Pin:   pin,
		Field: f,
	}, nil
}

func (s *Scope) applySelectToSpan(exp Exp, part ast.Select) (Exp, diag.Error) {
	name := part.Name.Str
	pin := part.Name.Pin

	switch name {
	case "len":
		return &SelectSpanLen{
			Exp: exp,
			Pin: pin,
			typ: s.Types.Known.Uint,
		}, nil
	case "ptr":
		return &SelectSpanPtr{
			Exp: exp,
			Pin: pin,
			typ: s.Types.getArrayPointer(exp.Type().Def.(Span).Type),
		}, nil
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("span does not have \"%s\" field", name),
		}
	}
}

func (s *Scope) applySelectToPointer(exp Exp, part ast.Select) (Exp, diag.Error) {
	typ := exp.Type().getDerefType()
	name := part.Name.Str
	pin := part.Name.Pin

	switch typ.Kind {
	case tpk.Custom:
		c := typ.Def.(*Custom)
		m := c.getMethod(name)
		if m != nil {
			return bindMethodToPointerReceiver(exp, m)
		}

		if c.Type.Kind != tpk.Struct {
			panic("not implemented")
		}

		s := c.Type.Def.(*Struct)
		f := s.getField(name)
		if f == nil {
			return nil, &diag.SimpleMessageError{
				Pin:  pin,
				Text: fmt.Sprintf("type %s does not have field or method named \"%s\"", typ, name),
			}
		}

		return &DerefSelectField{
			Pin:   pin,
			Exp:   exp,
			Field: f,
		}, nil
	default:
		panic(fmt.Sprintf("unexpected %s type", typ))
	}
}

func bindMethodToPointerReceiver(r Exp, m *Symbol) (*BoundMethod, diag.Error) {
	got := r.Type()

	want := m.Def.(*Fun).Receiver
	if want != got {
		switch {
		case want.Kind == tpk.Ref && got.Kind == tpk.Pointer:
			return nil, &diag.SimpleMessageError{
				Pin:  r.Span().Pin,
				Text: fmt.Sprintf("method \"%s\" has ref receiver type, but being used with pointer value", m.Name),
			}
		case want.Kind != tpk.Ref && want.Kind != tpk.Pointer:
			return nil, &diag.SimpleMessageError{
				Pin:  r.Span().Pin,
				Text: fmt.Sprintf("method \"%s\" has value receiver type, but being used with pointer value (sugar dereferencing is not allowed)", m.Name),
			}
		case want.Kind == tpk.Pointer && got.Kind == tpk.Ref:
			// auto conversion is allowed
		default:
			panic(fmt.Sprintf("unexpected combination of types %s and %s", want, got))
		}
	}

	return &BoundMethod{
		Receiver: r,
		Symbol:   m,
	}, nil
}

func (s *Scope) translateDotName(hint *Hint, d ast.DotName) (Exp, diag.Error) {
	name := d.Name
	pin := d.Pin

	entry := hint.lookupDotName(name)
	if entry == nil {
		if hint.Enum == nil {
			return nil, &diag.SimpleMessageError{
				Pin:  pin,
				Text: fmt.Sprintf("expression context does not have enum type hint to resolve \".%s\"", name),
			}
		}

		// TODO: supply enum type name to error
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("enum does not have \"%s\" entry", name),
		}
	}

	return entry.Value.WithPin(pin), nil
}

func (s *Scope) translateSymbolExp(sym ast.Symbol) (Exp, diag.Error) {
	name := sym.Name
	pin := sym.Pin

	symbol := s.Lookup(name)
	if symbol == nil {
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("name \"%s\" refers to undefined symbol", name),
		}
	}

	switch symbol.Kind {
	case smk.Const:
		return symbol.Def.(StaticValue).Exp, nil
	case smk.Var, smk.Loop, smk.Param:
		return &SymExp{Pin: pin, Symbol: symbol}, nil
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  pin,
			Text: fmt.Sprintf("symbol \"%s\" refers to %s, which cannot used as operand or expression", name, symbol.Kind),
		}
	}
}

func (s *Scope) translateBinaryExp(hint *Hint, exp ast.Binary) (Exp, diag.Error) {
	a, err := s.TranslateExp(hint, exp.A)
	if err != nil {
		return nil, err
	}
	b, err := s.TranslateExp(hint, exp.B)
	if err != nil {
		return nil, err
	}

	ta := a.Type()
	tb := b.Type()
	if ta.IsStatic() && tb.IsStatic() {
		return s.evalConstBinaryExp(a, b, exp.Op)
	}

	// boolean simplification when one side of expression is static
	if ta.Kind == tpk.Boolean && ta.IsStatic() {
		// tb is not static here
		if tb.Kind != tpk.Boolean {
			return nil, &diag.SimpleMessageError{
				Pin:  b.Span().Pin,
				Text: fmt.Sprintf("incompatible types in binary expression bool and %s", tb),
			}
		}

		v := a.(*Boolean).Val
		op := exp.Op
		switch op.Kind {
		case bok.And:
			if !v {
				return a, nil
			}
			return b, nil
		case bok.Or:
			if v {
				return a, nil
			}
			return b, nil
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("operation %s is not defined for boolean type", op.Kind),
			}
		}
	}
	if tb.Kind == tpk.Boolean && tb.IsStatic() {
		// ta is not static here
		if ta.Kind != tpk.Boolean {
			return nil, &diag.SimpleMessageError{
				Pin:  a.Span().Pin,
				Text: fmt.Sprintf("incompatible types in binary expression %s and bool", ta),
			}
		}

		v := b.(*Boolean).Val
		op := exp.Op
		switch op.Kind {
		case bok.And:
			if v {
				return a, nil
			}
			return s.Types.MakeBoolExp(a.Span().Pin, a, false), nil
		case bok.Or:
			if !v {
				return a, nil
			}
			return s.Types.MakeBoolExp(a.Span().Pin, a, true), nil
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("operation %s is not defined for boolean type", op.Kind),
			}
		}
	}

	typ, err := s.Types.deduceBinaryExpType(a, b, exp.Op)
	if err != nil {
		return nil, err
	}

	return &Binary{
		typ: typ,
		Op:  exp.Op,
		A:   a,
		B:   b,
	}, nil
}

// type checks binary expression and returns its resulting type
func (x *TypeIndex) deduceBinaryExpType(a, b Exp, op BinOp) (*Type, diag.Error) {
	ta := a.Type()
	tb := b.Type()
	if ta == tb {
		// TODO: check operator
		//
		// Boths types are equal and therefore cannot be static since
		// this function is only used when at most one of a or b is static.
		return x.checkBinaryForType(a.Type(), op)
	}

	if ta.IsStatic() {
		return x.deduceBinaryExpTypeA(a, b, op)
	}

	if tb.IsStatic() {
		return x.deduceBinaryExpTypeB(a, b, op)
	}

	switch op.Kind {
	case bok.LeftShift, bok.RightShift:
		if ta.Kind == tpk.Integer && tb.Kind == tpk.Integer {
			return ta, nil
		}
	}

	return nil, &diag.SimpleMessageError{
		Pin:  op.Pin,
		Text: fmt.Sprintf("type %s and %s are incompatible for binary operation", ta, tb),
	}
}

// type checks binary expression and returns its resulting type
//
// a has static type, b has runtime type
func (x *TypeIndex) deduceBinaryExpTypeA(a, b Exp, op BinOp) (*Type, diag.Error) {
	ta := a.Type()
	tb := b.Type()

	switch tb.Kind {
	case tpk.Integer:
		switch ta.Kind {
		case tpk.Integer:
			if ta.Size == 0 {
				return x.getBinaryForIntegerType(tb, op)
			}

			panic("sized static integers not implemented")
		case tpk.Rune:
			if tb.IsSigned() {
				return nil, &diag.SimpleMessageError{
					Pin:  op.Pin,
					Text: "binary operation on rune and signed integer",
				}
			}

			v := a.(*Rune).Val
			switch tb.Size {
			case 1:
				if v > 0xFF {
					return nil, &diag.SimpleMessageError{
						Pin:  op.Pin,
						Text: fmt.Sprintf("rune value 0x%X cannot fit into u8", v),
					}
				}
			case 2:
				if v > 0xFFFF {
					return nil, &diag.SimpleMessageError{
						Pin:  op.Pin,
						Text: fmt.Sprintf("rune value 0x%X cannot fit into u8", v),
					}
				}
			}
			return x.getBinaryForIntegerType(tb, op)
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("binary operation on incompatible types %s and %s", ta, tb),
			}
		}
	case tpk.Pointer, tpk.ArrayPointer:
		if ta.Kind != tpk.Nil {
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("binary operation on incompatible types %s and %s", ta, tb),
			}
		}

		switch op.Kind {
		case bok.Equal, bok.NotEqual:
			return x.Known.Bool, nil
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("operation %s is not defined for pointer type", op.Kind),
			}
		}
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  op.Pin,
			Text: fmt.Sprintf("type %s does not have binary operations", tb),
		}
	}
}

// type checks binary expression and returns its resulting type
//
// a has runtime type, b has static type
func (x *TypeIndex) deduceBinaryExpTypeB(a, b Exp, op BinOp) (*Type, diag.Error) {
	return x.deduceBinaryExpTypeA(b, a, op)
}

// returns resulting type of binary operation on integer types
func (x *TypeIndex) getBinaryForIntegerType(typ *Type, op BinOp) (*Type, diag.Error) {
	switch op.Kind {
	case bok.Add, bok.Sub, bok.Mul, bok.Mod, bok.Div, bok.BitAnd, bok.BitOr, bok.Xor, bok.LeftShift, bok.RightShift:
		return typ, nil
	case bok.Equal, bok.NotEqual, bok.Greater, bok.GreaterOrEqual, bok.Less, bok.LessOrEqual:
		return x.Known.Bool, nil
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  op.Pin,
			Text: fmt.Sprintf("operation %s is not defined for integer type", op.Kind),
		}
	}
}

// Returns resulting type for binary operation when both operands' types are equal
// and not static.
func (x *TypeIndex) checkBinaryForType(typ *Type, op BinOp) (*Type, diag.Error) {
	switch typ.Kind {
	case tpk.Integer:
		return x.getBinaryForIntegerType(typ, op)
	case tpk.String:
		switch op.Kind {
		case bok.Equal, bok.NotEqual:
			return x.Known.Bool, nil
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("operation %s is not defined for runtime string type", op.Kind),
			}
		}
	case tpk.Boolean:
		switch op.Kind {
		case bok.Equal, bok.NotEqual, bok.And, bok.Or:
			return x.Known.Bool, nil
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  op.Pin,
				Text: fmt.Sprintf("operation %s is not defined for boolean type", op.Kind),
			}
		}
	default:
		return nil, &diag.SimpleMessageError{
			Pin:  op.Pin,
			Text: fmt.Sprintf("type %s does not have binary operations", typ),
		}
	}
}
