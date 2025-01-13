package ast

import (
	"github.com/mebyus/ku/goku/compiler/enums/exk"
	"github.com/mebyus/ku/goku/compiler/enums/tnk"
	"github.com/mebyus/ku/goku/compiler/enums/tsk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Node common interface for any node inside a tree.
type Node interface {
	Span() source.Span
	String() string
}

// Exp node that represents expression.
type Exp interface {
	Node

	Kind() exk.Kind
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

// Word represents a single word token usage inside a tree.
// For expressions use Symbol instead.
type Word struct {
	// String that constitues the word.
	Str string

	Pin source.Pin
}

func (w Word) Span() source.Span {
	return source.Span{Pin: w.Pin, Len: uint32(len(w.Str))}
}

// Trait container object for passing around node attributes and properties.
type Traits struct {
	// List of node's properties.
	Props *[]Prop

	// True for public top level nodes.
	Pub bool
}

type Prop struct{}
