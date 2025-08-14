package asm

import (
	"io"

	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/vm/asm/compiler"
	"github.com/mebyus/ku/goku/vm/asm/parser"
	"github.com/mebyus/ku/goku/vm/kvx"
)

func Compile(code io.Reader) (*kvx.Program, error) {
	t, err := srcmap.NewTextFromReader(code)
	if err != nil {
		return nil, err
	}
	return CompileText(t)
}

func CompileText(text *srcmap.Text) (*kvx.Program, error) {
	t, err := parser.Parse(text)
	if err != nil {
		return nil, err
	}
	p, err := compiler.Compile(t)
	if err != nil {
		return nil, err
	}
	prog := Assemble(p)
	return prog, nil
}
