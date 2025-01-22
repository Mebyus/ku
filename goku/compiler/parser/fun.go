package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) Fun(traits ast.Traits) diag.Error {
	if p.n.Kind == token.LeftParen {
		return p.Method(traits)
	}

	p.advance() // skip "fun"

	if p.c.Kind != token.Word {
		return p.unexpected()
	}
	name := p.word()

	signature, err := p.signature()
	if err != nil {
		return err
	}

	if p.c.Kind != token.LeftCurly {
		// d := ast.FunStub{
		// 	Signature: signature,
		// 	Name:      name,
		// 	Traits:    traits,
		// }
		// p.atom.Nodes = append(p.atom.Nodes, ast.TopIndex{Kind: ast.NodeStub, Index: uint32(len(p.atom.Decs))})
		// p.atom.Decs = append(p.atom.Decs, d)
		return p.unexpected()
	}

	body, err := p.Block()
	if err != nil {
		return err
	}
	p.text.AddFun(ast.Fun{
		Signature: signature,
		Name:      name,
		Body:      body,
		Traits:    traits,
	})
	return nil
}

func (p *Parser) signature() (ast.Signature, diag.Error) {
	params, err := p.Params()
	if err != nil {
		return ast.Signature{}, err
	}

	if p.c.Kind != token.RightArrow {
		return ast.Signature{Params: params}, nil
	}

	p.advance() // skip "=>"

	if p.c.Kind == token.Never {
		p.advance() // skip "never"
		return ast.Signature{
			Params: params,
			Never:  true,
		}, nil
	}

	result, err := p.ResultTypeSpec()
	if err != nil {
		return ast.Signature{}, err
	}

	return ast.Signature{
		Params: params,
		Result: result,
	}, nil
}

func (p *Parser) Params() ([]ast.Param, diag.Error) {
	if p.c.Kind != token.LeftParen {
		return nil, p.unexpected()
	}
	p.advance() // skip "("

	var params []ast.Param
	for {
		if p.c.Kind == token.RightParen {
			p.advance() // skip ")"
			return params, nil
		}

		param, err := p.Param()
		if err != nil {
			return nil, err
		}
		params = append(params, param)

		if p.c.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.c.Kind == token.RightParen {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}

func (p *Parser) Param() (ast.Param, diag.Error) {
	if p.c.Kind != token.Word {
		return ast.Param{}, p.unexpected()
	}
	name := p.word()

	if p.c.Kind != token.Colon {
		return ast.Param{}, p.unexpected()
	}
	p.advance() // consume ":"

	typ, err := p.TypeSpec()
	if err != nil {
		return ast.Param{}, err
	}

	return ast.Param{
		Name: name,
		Type: typ,
	}, nil
}
