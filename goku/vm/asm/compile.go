package asm

import (
	"io"

	"github.com/mebyus/ku/goku/compiler/sm"
	"github.com/mebyus/ku/goku/vm/asm/compiler"
	"github.com/mebyus/ku/goku/vm/asm/parser"
	"github.com/mebyus/ku/goku/vm/kvx"
)

func Compile(code io.Reader) (*kvx.Program, error) {
	t, err := sm.NewTextFromReader(code)
	if err != nil {
		return nil, err
	}
	return CompileText(t)
}

func CompileText(text *sm.Text) (*kvx.Program, error) {
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
