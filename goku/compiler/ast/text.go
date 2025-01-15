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
	// Kind determines to which slice of nodes applies Index.
	TopList []TopNodeIndex

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
