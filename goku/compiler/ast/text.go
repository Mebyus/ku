package ast

import "github.com/mebyus/ku/goku/compiler/enums/tnk"

// Text smallest piece of processed source code inside a unit. In most
// cases this represents a file with source code. Exceptions may include
// source code generated at compile time or text that comes from a string
// during automated testing.
//
// <Text> = { <ImportBlock> } { <TopNode> }
//
// All top nodes inside text are listed in order they appear in source code.
type Text struct {
	ImportBlocks []ImportBlock

	// List of (Kind, Index) pairs for all top level nodes.
	// Elements in this list are in the same order as they appear in source text.
	//
	// Kind determines to which slice of nodes Index applies.
	OrderIndex []TopNodeIndex

	// List of top custom type definition nodes.
	Types []Type

	// List of top constant definition nodes.
	Constants []TopLet

	// List of top variable definition nodes.
	Variables []TopVar

	// List of top function definition nodes.
	Functions []Fun

	// List of unit test functions.
	Tests []Fun

	// List of top function declaration nodes.
	FunStubs []FunStub

	// List of method nodes.
	Methods []Method
}

type TopNodeIndex struct {
	Index uint32
	Kind  tnk.Kind
}

func New() *Text {
	return &Text{}
}

func (t *Text) AddType(typ Type) {
	t.OrderIndex = append(t.OrderIndex, TopNodeIndex{
		Kind:  tnk.Type,
		Index: uint32(len(t.Types)),
	})
	t.Types = append(t.Types, typ)
}

func (t *Text) AddVar(v Var) {
	t.OrderIndex = append(t.OrderIndex, TopNodeIndex{
		Kind:  tnk.Var,
		Index: uint32(len(t.Variables)),
	})
	t.Variables = append(t.Variables, TopVar{Var: v})
}

func (t *Text) AddLet(l Let) {
	t.OrderIndex = append(t.OrderIndex, TopNodeIndex{
		Kind:  tnk.Let,
		Index: uint32(len(t.Constants)),
	})
	t.Constants = append(t.Constants, TopLet{Let: l})
}

func (t *Text) AddFun(f Fun) {
	t.OrderIndex = append(t.OrderIndex, TopNodeIndex{
		Kind:  tnk.Fun,
		Index: uint32(len(t.Functions)),
	})
	t.Functions = append(t.Functions, f)
}
