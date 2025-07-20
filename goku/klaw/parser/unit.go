package parser

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
	"github.com/mebyus/ku/goku/klaw/ast"
	"github.com/mebyus/ku/goku/klaw/token"
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
	case token.Exe:
		return p.exe()
	default:
		return nil, p.unexpected()
	}
}

func (p *Parser) imp() (ast.Dir, diag.Error) {
	p.advance() // skip "import"

	var originPin srcmap.Pin
	var originName string
	if p.peek.Kind == token.Word {
		originPin = p.peek.Pin
		originName = p.peek.Data
		p.advance() // skip origin name
	}
	o, ok := origin.Parse(originName)
	if !ok {
		return nil, &diag.SimpleMessageError{
			Pin:  originPin,
			Text: fmt.Sprintf("unexpected import origin \"%s\"", originName),
		}
	}

	if p.peek.Kind == token.String {
		return p.simpleImport(o)
	}
	if p.peek.Kind != token.LeftCurly {
		return nil, p.unexpected()
	}
	p.advance() // skip "{"

	var imports []ast.ImportString
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			break
		}

		if p.peek.Kind != token.String {
			return nil, p.unexpected()
		}
		pin := p.peek.Pin
		val := p.peek.Data
		p.advance() // skip string

		if p.peek.Kind != token.Semicolon {
			return ast.Import{}, p.unexpected()
		}
		p.advance() // skip ";"

		imports = append(imports, ast.ImportString{
			Pin: pin,
			Val: val,
		})
	}

	return ast.ImportBlock{
		Origin:  o,
		Imports: imports,
	}, nil
}

func (p *Parser) simpleImport(o origin.Origin) (ast.Import, diag.Error) {
	pin := p.peek.Pin
	val := p.peek.Data
	p.advance() // skip string

	if p.peek.Kind != token.Semicolon {
		return ast.Import{}, p.unexpected()
	}
	p.advance() // skip ";"

	return ast.Import{
		Origin: o,

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

func (p *Parser) exe() (ast.Exe, diag.Error) {
	p.advance() // skip "exe"

	block, err := p.block()
	if err != nil {
		return ast.Exe{}, err
	}

	return ast.Exe{Block: block}, nil
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
