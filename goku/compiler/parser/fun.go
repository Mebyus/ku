package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) topFun(traits ast.Traits) diag.Error {
	f, err := p.Fun(traits)
	if err != nil {
		return err
	}
	p.text.AddFun(f)
	return nil
}

func (p *Parser) Fun(traits ast.Traits) (ast.Fun, diag.Error) {
	p.advance() // skip "fun"

	err := p.unsafe(&traits)
	if err != nil {
		return ast.Fun{}, err
	}

	f, err := p.fun()
	if err != nil {
		return ast.Fun{}, err
	}

	f.Traits = traits
	return f, nil
}

func (p *Parser) TestFun(traits ast.Traits) diag.Error {
	p.advance() // skip "test"

	t, err := p.testFun()
	if err != nil {
		return err
	}

	// TODO: maybe we need traits on tests?
	// t.Traits = traits
	p.text.AddTest(t)
	return nil
}

func (p *Parser) testFun() (ast.TestFun, diag.Error) {
	if p.peek.Kind != token.Word {
		return ast.TestFun{}, p.unexpected()
	}
	name := p.word()

	body, err := p.Block()
	if err != nil {
		return ast.TestFun{}, err
	}

	return ast.TestFun{
		Name: name,
		Body: body,
	}, nil
}

func (p *Parser) FunStub(traits ast.Traits) diag.Error {
	p.advance() // skip "#stub"

	if p.peek.Kind != token.Fun {
		return p.unexpected()
	}
	p.advance() // skip "fun"

	if p.peek.Kind != token.Word {
		return p.unexpected()
	}
	name := p.word()

	signature, err := p.signature()
	if err != nil {
		return err
	}

	p.text.AddStub(ast.FunStub{
		Name:      name,
		Signature: signature,
		Traits:    traits,
	})
	return nil
}

func (p *Parser) fun() (ast.Fun, diag.Error) {
	if p.peek.Kind != token.Word {
		return ast.Fun{}, p.unexpected()
	}
	name := p.word()

	signature, err := p.signature()
	if err != nil {
		return ast.Fun{}, err
	}

	if p.peek.Kind != token.LeftCurly {
		return ast.Fun{}, p.unexpected()
	}

	body, err := p.Block()
	if err != nil {
		return ast.Fun{}, err
	}

	return ast.Fun{
		Name:      name,
		Signature: signature,
		Body:      body,
	}, nil
}

func (p *Parser) signature() (ast.Signature, diag.Error) {
	params, err := p.Params()
	if err != nil {
		return ast.Signature{}, err
	}

	if p.peek.Kind != token.RightArrow {
		return ast.Signature{Params: params}, nil
	}

	p.advance() // skip "=>"

	if p.peek.Kind == token.Never {
		p.advance() // skip "#never"
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
	if p.peek.Kind != token.LeftParen {
		return nil, p.unexpected()
	}
	p.advance() // skip "("

	var params []ast.Param
	for {
		if p.peek.Kind == token.RightParen {
			p.advance() // skip ")"
			return params, nil
		}

		param, err := p.Param()
		if err != nil {
			return nil, err
		}
		params = append(params, param)

		if p.peek.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.peek.Kind == token.RightParen {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected()
		}
	}
}

func (p *Parser) Param() (ast.Param, diag.Error) {
	if p.peek.Kind != token.Word {
		return ast.Param{}, p.unexpected()
	}
	name := p.word()

	if p.peek.Kind != token.Colon {
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

// check for unsafe trait before function or method name
func (p *Parser) unsafe(traits *ast.Traits) diag.Error {
	if p.peek.Kind != token.Unsafe {
		return nil
	}

	p.advance() // skip "unsafe"

	if p.peek.Kind != token.Period {
		return p.unexpected()
	}
	p.advance() // skip "."

	traits.Unsafe = true
	return nil
}
