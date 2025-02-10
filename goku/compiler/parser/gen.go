package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Gen(traits ast.Traits) diag.Error {
	p.advance() // skip "gen"

	if p.c.Kind != token.Word {
		return p.unexpected()
	}
	name := p.word()

	if p.c.Kind != token.LeftParen {
		return p.unexpected()
	}

	if p.n.Kind == token.Ellipsis {
		return p.genBind(traits, name)
	}

	params, err := p.Params()
	if err != nil {
		return err
	}

	var control *ast.Static
	if p.c.Kind == token.HashCurly {
		block, err := p.Block()
		if err != nil {
			return err
		}
		control = &ast.Static{
			Pin:   block.Pin,
			Nodes: block.Nodes,
		}
	}

	p.text.AddGen(ast.Gen{
		Name:    name,
		Params:  params,
		Control: control,
	})
	return nil
}

func (p *Parser) genBind(traits ast.Traits, name ast.Word) diag.Error {
	p.advance() // skip "("
	p.advance() // skip "..."

	if p.c.Kind != token.RightParen {
		return p.unexpected()
	}
	p.advance() // skip ")"

	body, err := p.GenBlock()
	if err != nil {
		return err
	}

	p.text.AddGenBind(ast.GenBind{
		Name: name,
		Body: body,
	})
	return nil
}

func (p *Parser) GenBlock() (ast.GenBlock, diag.Error) {
	if p.c.Kind != token.LeftCurly {
		return ast.GenBlock{}, p.unexpected()
	}
	p.advance() // skip "{"

	var block ast.GenBlock
	for {
		if p.c.Kind == token.RightCurly {
			p.advance() // skip "}"
			return block, nil
		}

		err := p.genNode(&block)
		if err != nil {
			return ast.GenBlock{}, err
		}
	}
}

func (p *Parser) genNode(b *ast.GenBlock) diag.Error {
	switch p.c.Kind {
	case token.Type:
		return p.genType(b)
	case token.Fun:
		if p.n.Kind == token.LeftParen {
			return p.genMethod(b)
		}
		return p.genFun(b)
	case token.Const:
		return p.genConst(b)
	case token.Let:
		return p.genAlias(b)
	default:
		return p.unexpected()
	}
}

func (p *Parser) genFun(b *ast.GenBlock) diag.Error {
	f, err := p.Fun(ast.Traits{}) // TODO: parse traits
	if err != nil {
		return err
	}
	b.AddFun(f)
	return nil
}

func (p *Parser) genMethod(b *ast.GenBlock) diag.Error {
	m, err := p.Method(ast.Traits{}) // TODO: parse traits
	if err != nil {
		return err
	}
	b.AddMethod(m)
	return nil
}

func (p *Parser) genConst(b *ast.GenBlock) diag.Error {
	l, err := p.Const()
	if err != nil {
		return err
	}

	b.AddConst(ast.TopConst{
		Const:  l,
		Traits: ast.Traits{}, // TODO: parse traits
	})
	return nil
}

func (p *Parser) genType(b *ast.GenBlock) diag.Error {
	t, err := p.Type(ast.Traits{}) // TODO: parse traits
	if err != nil {
		return err
	}
	b.AddType(t)
	return nil
}

func (p *Parser) genAlias(b *ast.GenBlock) diag.Error {
	a, err := p.Alias()
	if err != nil {
		return err
	}

	b.AddAlias(ast.TopAlias{
		Alias:  a,
		Traits: ast.Traits{}, // TODO: parse traits
	})
	return nil
}
