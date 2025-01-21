package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) For() (ast.Statement, diag.Error) {
	if p.n.Kind == token.LeftCurly {
		return p.Loop()
	}

	// p.advance() // skip "for"
	// if p.tok.Kind == token.Word && (p.next.Kind == token.In || p.next.Kind == token.Colon) {
	// 	return p.forIn()
	// }
	// return p.forIf()
	panic("not implemented")
}

func (p *Parser) Loop() (ast.Loop, diag.Error) {
	p.advance() // skip "for"

	body, err := p.Block()
	if err != nil {
		return ast.Loop{}, err
	}

	// TODO: maybe generate warning
	// if len(body.Nodes) == 0 {
	// 	return ast.Loop{}, fmt.Errorf("%s for loop without condition cannot have empty body", pos.Short())
	// }

	return ast.Loop{Body: body}, nil
}
