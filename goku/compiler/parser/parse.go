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

	switch p.c.Kind {
	case token.Type:
		return p.Type(traits)
	case token.Fun:
		return p.Fun(traits)
	case token.Let:
		return p.TopLet(traits)
	case token.Var:
		return p.TopVar(traits)
	case token.Test:
		return p.Test(traits)
	case token.Stub:
		return p.Stub(traits)
	// case token.Pub:
	// 	traits.Pub = true
	// 	return p.topPub(traits)
	default:
		return p.unexpected()
	}
}

func (p *Parser) gatherProps() diag.Error {
	for {
		if p.c.Kind != token.PropStart {
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
