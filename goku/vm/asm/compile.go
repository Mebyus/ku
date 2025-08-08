package asm

import (
	"io"

	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/vm/kvx"
)

func Compile(code io.Reader) (*kvx.Program, error) {
	return &kvx.Program{}, nil
}

func CompileText(text *srcmap.Text) (*kvx.Program, error) {
	return &kvx.Program{}, nil
}
