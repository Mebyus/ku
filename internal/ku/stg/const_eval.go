package stg

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/ku/token"
)

func (t *Typer) evalConstDef(s *Symbol, def *ast.Const) {
	var typ *Type
	if def.Type != nil {
		typ = t.LookupType(&t.unit.Scope, def.Type)
	}

	exp := t.evalConstExp(def.Exp)
	if typ == nil {
		s.Type = exp.Type()
		return
	}

	// TODO: typecheck declared type and expression type
}

func (t *Typer) evalConstExp(exp ast.Exp) Exp {
	switch e := exp.(type) {
	case *ast.InvExp:
		return t.makeInvExp(e.Pin)
		// already reported during parsing
	case *ast.Integer:
		return t.makeInteger(e.Pin, e.Val, uint16(e.Aux))
	case *ast.True:
		return t.makeBoolean(e.Pin, true)
	case *ast.False:
		return t.makeBoolean(e.Pin, false)
	case *ast.String:
	case *ast.Rune:
	case *ast.SymExp:
		return t.evalConstSymExp(e)
	case *ast.ParenExp:
		return t.evalConstExp(e.Exp)
	case *ast.UnExp:
	case *ast.BinExp:
		return t.evalConstBinExp(e)
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}

	panic(fmt.Sprintf("TODO: %T", exp))
}

func (t *Typer) evalConstSymExp(exp *ast.SymExp) Exp {
	name := exp.Name
	pin := exp.Pin

	symbol := t.unit.Scope.Lookup(name)
	if symbol == nil {
		t.report(pin, sx.ETypUnknownSymbol, fmt.Sprintf("unknown symbol \"%s\" used as expression", name))
		return t.makeInvExp(pin)
	}

	switch symbol.Kind {
	case symk.Const:
		symbol.rnum += 1
		return symbol.Def.(*StaticValue).Exp
	default:
		t.report(pin, 0, fmt.Sprintf("symbol \"%s\" cannot be used as operand in compile-time expression", name))
		return t.makeInvExp(pin)
	}
}

func (t *Typer) evalConstBinExp(exp *ast.BinExp) Exp {
	panic("TODO")
}

// Integer represents an integer constant (directly from source or evaluated)
// which value is known at compile time.
type Integer struct {
	operand

	Val uint64

	pin sx.Pin

	typ *Type

	// holds information whether this value was obtained
	// from source, evaluated or generated
	origin uint16

	// preferred format style:
	// 	decimal, binary, octal, hexadecimal
	Style uint16

	// Equals true for negative values.
	// Should be false when value holds 0.
	Neg bool
}

// flags for constant value origin
const (
	// from source code
	//
	// should be zero value
	ovalsrc = iota

	// from constant evaluation/propagation
	ovaleva

	// generated
	ovalgen
)

// Styles for integer constant values.
const (
	IntDec = token.DecInt
	IntBin = token.BinInt
	IntOct = token.OctInt
	IntHex = token.HexInt
)

// Explicit interface implementation check.
var _ Operand = &Integer{}

func (n *Integer) Type() *Type {
	return n.typ
}

func (n *Integer) Pin() sx.Pin {
	return n.pin
}

func (*Integer) Target() bool {
	return false
}

// create static unsized integer value that comes from source.
func (t *Typer) makeInteger(pin sx.Pin, v uint64, style uint16) *Integer {
	return &Integer{
		pin:   pin,
		Val:   v,
		typ:   t.common.Types.Static.Integer,
		Style: style,
	}
}

// Boolean represents a boolean constant true/false (directly from source or evaluated)
// which value is known at compile time.
type Boolean struct {
	operand

	pin sx.Pin

	typ *Type

	Val bool
}

// Explicit interface implementation check.
var _ Operand = &Boolean{}

func (b *Boolean) Type() *Type {
	return b.typ
}

func (b *Boolean) Pin() sx.Pin {
	return b.pin
}

func (b *Boolean) Target() bool {
	return false
}

// create static boolean value.
func (t *Typer) makeBoolean(pin sx.Pin, v bool) *Boolean {
	return &Boolean{
		pin: pin,
		Val: v,
		typ: t.common.Types.Static.Boolean,
	}
}
