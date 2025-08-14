package parser

import (
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/tokens"
)

func Parse(text *srcmap.Text) (*ast.Text, diag.Error) {
	p := FromText(text)
	return p.Text()
}

func (p *Parser) parse() diag.Error {
	for {
		if p.peek.Kind == tokens.EOF {
			return nil
		}

		err := p.top()
		if err != nil {
			return err
		}
	}
}

func (p *Parser) top() diag.Error {
	switch p.peek.Kind {
	case tokens.Fun:
		return p.topFun()
	case tokens.Data:
		panic("not implemented")
	case tokens.Entry:
		panic("not implemented")
	default:
		return p.unexpected()
	}
}
