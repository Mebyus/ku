package parser

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/sm"
	"github.com/mebyus/ku/goku/compiler/token"
)

func (p *Parser) ImportBlocks() ([]ast.ImportBlock, diag.Error) {
	var blocks []ast.ImportBlock
	for {
		if p.peek.Kind == token.Import {
			block, err := p.ImportBlock()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, block)
		} else {
			p.text.ImportBlocks = blocks
			return blocks, nil
		}
	}
}

func (p *Parser) ImportBlock() (ast.ImportBlock, diag.Error) {
	p.advance() // skip "import"

	origin, err := p.Origin()
	if err != nil {
		return ast.ImportBlock{}, err
	}

	if p.peek.Kind != token.LeftCurly {
		return ast.ImportBlock{}, p.unexpected()
	}
	p.advance() // skip "{"

	var imports []ast.Import
	for {
		if p.peek.Kind == token.RightCurly {
			p.advance() // skip "}"
			return ast.ImportBlock{
				Imports: imports,
				Origin:  origin,
			}, nil
		}

		i, err := p.Import()
		if err != nil {
			return ast.ImportBlock{}, err
		}
		imports = append(imports, i)
	}
}

func (p *Parser) Origin() (sm.Origin, diag.Error) {
	if p.peek.Kind != token.Word {
		return sm.Loc, nil
	}
	name := p.word()

	origin, ok := sm.ParseOrigin(name.Str)
	if !ok {
		return 0, &diag.UnknownOriginError{Name: name}
	}
	return origin, nil
}

func (p *Parser) Import() (ast.Import, diag.Error) {
	if p.peek.Kind != token.Word {
		return ast.Import{}, p.unexpected()
	}
	name := p.word()

	if p.peek.Kind != token.RightArrow {
		return ast.Import{}, p.unexpected()
	}
	p.advance() // skip "=>"

	if p.peek.Kind != token.String {
		return ast.Import{}, p.unexpected()
	}
	s := p.peek
	p.advance() // skip import string

	err := sm.CheckImportString(s.Data)
	if err != nil {
		return ast.Import{}, &diag.SimpleMessageError{
			Text: err.Error(),
			Pin:  s.Pin,
		}
	}

	str := ast.ImportString{
		Pin: s.Pin,
		Val: s.Data,
	}

	return ast.Import{
		Name:   name,
		String: str,
	}, nil
}
