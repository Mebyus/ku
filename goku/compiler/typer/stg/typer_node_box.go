package stg

import (
	"slices"

	"github.com/mebyus/ku/goku/compiler/ast"
)

// NodeBox is a container for gathering AST nodes from all unit texts.
type NodeBox struct {
	texts []*ast.Text

	// List of top custom type definition nodes.
	types []ast.Type

	// List of top constant definition nodes.
	consts []ast.TopConst

	// List of top variable definition nodes.
	vars []ast.TopVar

	// List of top function definition nodes.
	funs []ast.Fun

	// List of unit test function nodes.
	tests []ast.TestFun

	// List of top function stub nodes.
	stubs []ast.FunStub

	// List of method nodes.
	methods []ast.Method
}

// alloc is an efficient way to add Texts in bulk.
func (b *NodeBox) alloc(texts []*ast.Text) {
	b.texts = append(b.texts, texts...)

	funs := 0
	vars := 0
	tests := 0
	types := 0
	methods := 0
	stubs := 0
	consts := 0
	for _, t := range texts {
		funs += len(t.Functions)
		vars += len(t.Variables)
		tests += len(t.Tests)
		types += len(t.Types)
		stubs += len(t.FunStubs)
		consts += len(t.Constants)
		methods += len(t.Methods)
	}

	b.funs = slices.Grow(b.funs, funs)
	b.vars = slices.Grow(b.vars, vars)
	b.tests = slices.Grow(b.tests, tests)
	b.types = slices.Grow(b.types, types)
	b.methods = slices.Grow(b.methods, methods)
	b.stubs = slices.Grow(b.stubs, stubs)
	b.consts = slices.Grow(b.consts, consts)

	for _, text := range texts {
		b.gather(text)
	}
}

func (b *NodeBox) reset() {
	b.texts = b.texts[:0]
	b.consts = b.consts[:0]
	b.funs = b.funs[:0]
	b.stubs = b.stubs[:0]
	b.tests = b.tests[:0]
	b.vars = b.vars[:0]
	b.types = b.types[:0]
	b.methods = b.methods[:0]
}

func (b *NodeBox) addText(text *ast.Text) {
	b.texts = append(b.texts, text)
	b.gather(text)
}

func (b *NodeBox) gather(text *ast.Text) {
	for _, c := range text.Constants {
		b.consts = append(b.consts, c)
	}
	for _, t := range text.Types {
		b.types = append(b.types, t)
	}
	for _, v := range text.Variables {
		b.vars = append(b.vars, v)
	}
	for _, f := range text.Functions {
		b.funs = append(b.funs, f)
	}
	for _, t := range text.Tests {
		b.tests = append(b.tests, t)
	}
	for _, m := range text.Methods {
		b.methods = append(b.methods, m)
	}
	for _, s := range text.FunStubs {
		b.stubs = append(b.stubs, s)
	}
}
