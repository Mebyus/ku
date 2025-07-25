package ast

import (
	"slices"

	"github.com/mebyus/ku/goku/compiler/enums/tnk"
)

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
	OrderIndex []NodeIndex

	// List of top custom type definition nodes.
	Types []Type

	// List of top constant definition nodes.
	Constants []TopConst

	// List of top variable definition nodes.
	Variables []TopVar

	// List of top function definition nodes.
	Functions []Fun

	// List of top alias nodes.
	Aliases []TopAlias

	// List of unit test functions.
	Tests []TestFun

	// List of top function declaration nodes.
	FunStubs []FunStub

	// List of bag implementation registrations.
	RegBags []RegBag

	// List of method nodes.
	Methods []Method

	// List of top level static asserts.
	Musts []StaticMust

	// List of generic definition nodes.
	Generics []Gen

	// List of generic binds.
	GenBinds []GenBind

	// Optional build block. Always comes before imports.
	Build *Build
}

type NodeIndex struct {
	Index uint32
	Kind  tnk.Kind
}

func New() *Text {
	return &Text{}
}

func (t *Text) AddType(typ Type) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.Type,
		Index: uint32(len(t.Types)),
	})
	t.Types = append(t.Types, typ)
}

func (t *Text) AddVar(v TopVar) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.Var,
		Index: uint32(len(t.Variables)),
	})
	t.Variables = append(t.Variables, v)
}

func (t *Text) AddConst(l TopConst) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.Const,
		Index: uint32(len(t.Constants)),
	})
	t.Constants = append(t.Constants, l)
}

func (t *Text) AddAlias(a TopAlias) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.Alias,
		Index: uint32(len(t.Aliases)),
	})
	t.Aliases = append(t.Aliases, a)
}

func (t *Text) AddFun(f Fun) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.Fun,
		Index: uint32(len(t.Functions)),
	})
	t.Functions = append(t.Functions, f)
}

func (t *Text) AddTest(f TestFun) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.Test,
		Index: uint32(len(t.Tests)),
	})
	t.Tests = append(t.Tests, f)
}

func (t *Text) AddMethod(m Method) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.Method,
		Index: uint32(len(t.Methods)),
	})
	t.Methods = append(t.Methods, m)
}

func (t *Text) AddStub(s FunStub) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.FunStub,
		Index: uint32(len(t.FunStubs)),
	})
	t.FunStubs = append(t.FunStubs, s)
}

func (t *Text) AddRegBag(r RegBag) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.RegBag,
		Index: uint32(len(t.RegBags)),
	})
	t.RegBags = append(t.RegBags, r)
}

func (t *Text) AddMust(m StaticMust) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.Must,
		Index: uint32(len(t.Musts)),
	})
	t.Musts = append(t.Musts, m)
}

func (t *Text) AddGen(g Gen) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.Gen,
		Index: uint32(len(t.Generics)),
	})
	t.Generics = append(t.Generics, g)
}

func (t *Text) AddGenBind(b GenBind) {
	t.OrderIndex = append(t.OrderIndex, NodeIndex{
		Kind:  tnk.GenBind,
		Index: uint32(len(t.GenBinds)),
	})
	t.GenBinds = append(t.GenBinds, b)
}

func (t *Text) AppendTestNames(list []string) []string {
	if len(t.Tests) == 0 {
		return list
	}
	list = slices.Grow(list, len(t.Tests))
	for _, tt := range t.Tests {
		list = append(list, tt.Name.Str)
	}
	return list
}
