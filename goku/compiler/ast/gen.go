package ast

import "github.com/mebyus/ku/goku/compiler/enums/tnk"

// Gen represents generic (meta code) definition.
//
// Formal definition:
//
//	Gen => "gen" Name "(" ParamList ")" [ Control ]
//	Name => word
//	Control => Static
type Gen struct {
	Name Word

	Params []Param

	Control *Static
}

// GenBind represents generic (meta code) block bind.
//
// Formal definition:
//
//	gen Name "(" "..." ")" GenBlock
type GenBind struct {
	Body GenBlock
	Name Word
}

// GenBlock represents generic (meta code) block.
//
// Formal definition:
//
//	GenBlock => "{" { Const | Type | Method } "}"
type GenBlock struct {
	// List of (Kind, Index) pairs for block level nodes.
	// Elements in this list are in the same order as they appear in source text.
	//
	// Kind determines to which slice of nodes Index applies.
	OrderIndex []NodeIndex

	// List of type definition nodes.
	Types []Type

	// List of constant definition nodes.
	Constants []TopConst

	// List of function definition nodes.
	Functions []Fun

	// List of method nodes.
	Methods []Method

	// List of alias nodes.
	Aliases []TopAlias

	// List of lookup nodes.
	Lookups []Lookup
}

func (b *GenBlock) AddFun(f Fun) {
	b.OrderIndex = append(b.OrderIndex, NodeIndex{
		Kind:  tnk.Fun,
		Index: uint32(len(b.Functions)),
	})
	b.Functions = append(b.Functions, f)
}

func (b *GenBlock) AddMethod(m Method) {
	b.OrderIndex = append(b.OrderIndex, NodeIndex{
		Kind:  tnk.Method,
		Index: uint32(len(b.Methods)),
	})
	b.Methods = append(b.Methods, m)
}

func (b *GenBlock) AddType(typ Type) {
	b.OrderIndex = append(b.OrderIndex, NodeIndex{
		Kind:  tnk.Type,
		Index: uint32(len(b.Types)),
	})
	b.Types = append(b.Types, typ)
}

func (b *GenBlock) AddConst(l TopConst) {
	b.OrderIndex = append(b.OrderIndex, NodeIndex{
		Kind:  tnk.Const,
		Index: uint32(len(b.Constants)),
	})
	b.Constants = append(b.Constants, l)
}

func (b *GenBlock) AddAlias(a TopAlias) {
	b.OrderIndex = append(b.OrderIndex, NodeIndex{
		Kind:  tnk.Alias,
		Index: uint32(len(b.Aliases)),
	})
	b.Aliases = append(b.Aliases, a)
}

func (b *GenBlock) AddLookup(l Lookup) {
	b.OrderIndex = append(b.OrderIndex, NodeIndex{
		Kind:  tnk.Lookup,
		Index: uint32(len(b.Lookups)),
	})
	b.Lookups = append(b.Lookups, l)
}
