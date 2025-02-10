package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) parse() diag.Error {
	for {
		if p.c.Kind == token.EOF {
			// if p.c.Data != "" {
			// 	return fmt.Errorf(p.tok.Lit)
			// }
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
	if p.c.Kind == token.Pub {
		traits.Pub = true
		p.advance() // skip "pub"
	}

	switch p.c.Kind {
	case token.Type:
		return p.topType(traits)
	case token.Fun:
		if p.n.Kind == token.LeftParen {
			return p.topMethod(traits)
		}
		return p.topFun(traits)
	case token.Const:
		return p.topConst(traits)
	case token.Var:
		return p.TopVar(traits)
	case token.Test:
		return p.Test(traits)
	case token.Stub:
		return p.FunStub(traits)
	case token.Gen:
		return p.Gen(traits)
	default:
		return p.unexpected()
	}
}

func (p *Parser) gatherProps() diag.Error {
	for {
		if p.c.Kind != token.HashCurly {
			return nil
		}

		prop, err := p.prop()
		if err != nil {
			return err
		}
		p.props = append(p.props, prop)
	}
}

func (p *Parser) prop() (ast.Prop, diag.Error) {
	panic("not implemented")
}

func (p *Parser) takeProps() *[]ast.Prop {
	if len(p.props) == 0 {
		return nil
	}
	props := p.props
	p.props = nil
	return &props
}

// Consume and then return a word.
func (p *Parser) word() ast.Word {
	word := ast.Word{
		Pin: p.c.Pin,
		Str: p.c.Data,
	}
	p.advance()
	return word
}
