package stg

import (
	"slices"

	"github.com/mebyus/ku/internal/ku/ast"
)

// Alloc is a more efficient way (compared to Add) to add multiple Texts.
func (t *Typer) alloc(texts []*ast.Text) {
	t.box.alloc(texts)
}

// NodeBox is a container for gathering AST nodes from all unit texts.
type NodeBox struct {
	funs  []ast.Fun
	stubs []ast.FunStub
}

// alloc is an efficient way to add Texts in bulk.
func (b *NodeBox) alloc(texts []*ast.Text) {
	// b.texts = append(b.texts, texts...)

	funs := 0
	// vars := 0
	// tests := 0
	// types := 0
	// methods := 0
	stubs := 0
	// consts := 0
	for _, t := range texts {
		funs += len(t.Funs)
		// vars += len(t.Variables)
		// tests += len(t.Tests)
		// types += len(t.Types)
		stubs += len(t.Stubs)
		// consts += len(t.Constants)
		// methods += len(t.Methods)
	}

	b.funs = slices.Grow(b.funs, funs)
	// b.vars = slices.Grow(b.vars, vars)
	// b.tests = slices.Grow(b.tests, tests)
	// b.types = slices.Grow(b.types, types)
	// b.methods = slices.Grow(b.methods, methods)
	b.stubs = slices.Grow(b.stubs, stubs)
	// b.consts = slices.Grow(b.consts, consts)

	for _, text := range texts {
		b.gather(text)
	}
}

func (b *NodeBox) reset() {
	// b.texts = b.texts[:0]
	// b.consts = b.consts[:0]
	b.funs = b.funs[:0]
	b.stubs = b.stubs[:0]
	// b.tests = b.tests[:0]
	// b.vars = b.vars[:0]
	// b.types = b.types[:0]
	// b.methods = b.methods[:0]
}

func (b *NodeBox) gather(text *ast.Text) {
	for _, f := range text.Funs {
		b.funs = append(b.funs, f)
	}
	for _, s := range text.Stubs {
		b.stubs = append(b.stubs, s)
	}
}
