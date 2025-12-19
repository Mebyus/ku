package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) parse() diag.Error {
	for {
		if p.peek.Kind == token.EOF {
			return nil
		}

		err := p.top()
		if err != nil {
			return err
		}
	}
}

func (p *Parser) top() diag.Error {
	err := p.gatherProps()
	if err != nil {
		return err
	}
	traits := ast.Traits{Props: p.takeProps()}
	if p.peek.Kind == token.Pub {
		traits.Pub = true
		p.advance() // skip "pub"
	}

	switch p.peek.Kind {
	case token.Type:
		return p.topType(traits)
	case token.Fun:
		if p.next.Kind == token.LeftParen {
			return p.topMethod(traits)
		}
		return p.topFun(traits)
	case token.Const:
		return p.topConst(traits)
	case token.Var:
		return p.topVar(traits)
	case token.Let:
		return p.topAlias(traits)
	case token.Test:
		return p.TestFun(traits)
	case token.Stub:
		return p.FunStub(traits)
	case token.StaticMust:
		return p.topStaticMust()
	case token.Gen:
		return p.Gen(traits)
	case token.Bag:
		return p.RegBag()
	default:
		return p.unexpected()
	}
}

// Consume and then return a word.
func (p *Parser) word() ast.Word {
	word := ast.Word{
		Pin: p.peek.Pin,
		Str: p.peek.Data,
	}
	p.advance()
	return word
}
