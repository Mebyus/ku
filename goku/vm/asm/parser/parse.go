package parser

import (
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/vm/asm/ast"
)

func Parse(text *srcmap.Text) (*ast.Text, diag.Error) {
	return &ast.Text{}, nil
}
