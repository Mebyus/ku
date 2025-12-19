package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/enums/stk"
	"github.com/mebyus/ku/goku/compiler/enums/tnk"
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/sm"
)

// Node common interface for any node inside a tree.
type Node interface {
	Span() sm.Span
	String() string
}

// Exp node that represents an arbitrary expression.
type Exp interface {
	Node

	Kind() exk.Kind

	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_exp()
}

type Operand interface {
	Exp

	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_operand()
}

// Part bag for chain parts.
//
// Note that Part does not have _exp() and therefore should not be an expression.
type Part interface {
	Node

	Kind() exk.Kind

	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_part()
}

// Top node that represents top level node.
type Top interface {
	Node

	Kind() tnk.Kind
}

// TypeSpec node that represents type specifier of any kind.
type TypeSpec interface {
	Node

	Kind() tsk.Kind
}

// Statement node that represents statement of any kind.
type Statement interface {
	Node

	Kind() stk.Kind
}

// Word represents a single word token usage inside a tree.
// For expressions use Symbol instead.
type Word struct {
	// String that constitues the word.
	Str string

	Pin sm.Pin
}

func (w Word) Span() sm.Span {
	return sm.Span{Pin: w.Pin, Len: uint32(len(w.Str))}
}

// Traits container object for passing around node attributes and properties.
type Traits struct {
	// List of node's properties.
	Props *[]Prop

	// True for public top level nodes.
	Pub bool

	// True for unsafe functions and methods.
	Unsafe bool
}

type Prop struct {
	Exp  Exp
	Name string
	Pin  sm.Pin
}

// Embed this to quickly implement _exp() discriminator from Exp interface.
// Do not use it for anything else.
type nodeExp struct{}

func (nodeExp) _exp() {}

// Embed this to quickly implement _exp() and _operand() discriminators from Operand interface.
// Do not use it for anything else.
type nodeOperand struct{}

func (nodeOperand) _exp() {}

func (nodeOperand) _operand() {}

// Embed this to quickly implement _part() discriminator from Part interface.
// Do not use it for anything else.
type nodePart struct{}

func (nodePart) _part() {}
