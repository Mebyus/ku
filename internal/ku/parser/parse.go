package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) parse() {
	for !p.stop {
		p.top()
	}

	if p.text.IsOk() && len(p.text.Errors) != 0 {
		p.text.Status = ast.Flawed
	}
}

func (p *Parser) top() {
	switch p.peek.Kind {
	case token.Fun:
		p.topFun()
	case token.EOF:
		p.abort(ast.Ok)
	default:
		p.topError(p.peek.Pin, fmt.Sprintf("expected top-level node start, found %s token instead", p.peek.Kind))
	}
}

func (p *Parser) topFun() {
	fun := p.Fun()
	if fun != nil {
		p.text.Funs = append(p.text.Funs, fun)
	}
}
