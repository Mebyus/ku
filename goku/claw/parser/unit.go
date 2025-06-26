package parser

import (
	"github.com/mebyus/ku/goku/claw/ast"
	"github.com/mebyus/ku/goku/claw/token"
	"github.com/mebyus/ku/goku/compiler/diag"
)

func (p *Parser) unitParse() diag.Error {
	for {
		if p.peek.Kind == token.EOF {
			return nil
		}

		err := p.unitTop()
		if err != nil {
			return err
		}
	}
}

func (p *Parser) unitTop() diag.Error {
	dir, err := p.dir()
	if err != nil {
		return err
	}

	p.unit.Dirs = append(p.unit.Dirs, dir)
	return nil
}

func (p *Parser) dir() (ast.Dir, diag.Error) {
	switch p.peek.Kind {
	case token.Import:
		return p.imp()
	case token.Include:
		return p.include()
	case token.Test:
		return p.test()
	default:
		return nil, p.unexpected()
	}
}

func (p *Parser) imp() (ast.Import, diag.Error) {
	p.advance() // skip "import"

	if p.peek.Kind != token.String {
		return ast.Import{}, p.unexpected()
	}
	pin := p.peek.Pin
	val := p.peek.Data
	p.advance() // skip string

	if p.peek.Kind != token.Semicolon {
		return ast.Import{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Import{
		Val: val,
		Pin: pin,
	}, nil
}

func (p *Parser) include() (ast.Include, diag.Error) {
	p.advance() // skip "include"

	if p.peek.Kind != token.String {
		return ast.Include{}, p.unexpected()
	}
	pin := p.peek.Pin
	val := p.peek.Data
	p.advance() // skip string

	if p.peek.Kind != token.Semicolon {
		return ast.Include{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Include{
		Val: val,
		Pin: pin,
	}, nil
}

func (p *Parser) test() (ast.Test, diag.Error) {
	p.advance() // skip "test"

	block, err := p.block()
	if err != nil {
		return ast.Test{}, err
	}

	return ast.Test{Block: block}, nil
}

func (p *Parser) block() (ast.Block, diag.Error) {
	if p.peek.Kind != token.LeftCurly {
		return ast.Block{}, p.unexpected()
	}
	pin := p.peek.Pin
	p.advance() // skip "{"

	var dirs []ast.Dir
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			return ast.Block{
				Pin:  pin,
				Dirs: dirs,
			}, nil
		}

		dir, err := p.dir()
		if err != nil {
			return ast.Block{}, err
		}
		dirs = append(dirs, dir)
	}
}
