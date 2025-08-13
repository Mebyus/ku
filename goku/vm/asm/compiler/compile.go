package compiler

import (
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/ir"
)

func Compile(text *ast.Text) (*ir.Program, diag.Error) {
	return &ir.Program{}, nil
}
