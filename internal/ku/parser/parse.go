package parser

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/ku/token"
)

func (p *Parser) parse() {
	p.imports()

	for p.status == 0 {
		p.top()
	}

	switch p.status {
	case ast.Stop:
		// normal exit after parsing whole text
		if p.nerr != 0 {
			p.text.Status = ast.Flawed
		}
		return
	}
}

func (p *Parser) imports() {
	for p.status == 0 {
		switch p.peek.Kind {
		case token.Import:
			p.iblock()
		default:
			return
		}
	}
}

func (p *Parser) top() {
	if p.peek.Kind == token.Pub {
		p.advance() // skip "pub"
	}

	switch p.peek.Kind {
	case token.Const:
		p.topConst()
	case token.Fun:
		p.topFun()
	case token.Type:
		p.topType()
	case token.EOF:
		p.stop(ast.Stop)
	default:
		p.topError(p.peek.Pin, fmt.Sprintf("expected top-level node start, found %s token instead", &p.peek))
	}
}

// parse import block
func (p *Parser) iblock() {
	p.advance() // skip "import"

	var oname string // origin name
	var opin sx.Pin  // origin pin
	if p.peek.Kind == token.Word {
		oname = p.peek.Data
		opin = p.peek.Pin
		p.advance() // skip origin word
	}

	var origin sx.Origin
	switch oname {
	case "", "loc":
		origin = sx.Loc
	case "std":
		origin = sx.Std
	case "pkg":
		origin = sx.Pkg
	default:
		p.report(opin, 0, fmt.Sprintf("unknown origin \"%s\"", oname))
		origin = sx.Bad
	}

	if p.peek.Kind != token.LeftParen {
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected \"(\" for import block start, found %s token instead", &p.peek))
		return
	}
	p.advance() // skip "("

	for {
		if p.peek.Kind == token.RightParen {
			p.advance() // skip ")"
			return
		}

		s := p.imp(origin)
		if s != 0 {
			return
		}
	}

}

// parse single import
func (p *Parser) imp(origin sx.Origin) ss {
	var name string
	var pin sx.Pin

	switch p.peek.Kind {
	case token.Word:
		name = p.peek.Data
		pin = p.peek.Pin
		p.advance() // skip import name
	case token.RightArrow:
		p.report(p.peek.Pin, 0, "missing import name")
	default:
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected import name, found %s token instead", &p.peek))
	}

	if p.peek.Kind != token.RightArrow {
		p.report(p.peek.Pin, 0, "missing \"->\" before import string")
	} else {
		p.advance() // skip "->"
	}

	var s string // import string
	var ipin sx.Pin
	switch p.peek.Kind {
	case token.String:
		s = p.peek.Data
		ipin = p.peek.Pin
		p.advance() // skip import string

		var verr string
		s, verr = sx.ValidImport(s)
		if s == "" && verr == "" {
			verr = "bad import string"
		}
		if verr != "" {
			p.report(ipin, 0, verr)
		}
	case token.Semicolon:
		p.report(p.peek.Pin, 0, "missing import string")
	default:
		p.report(p.peek.Pin, 0, fmt.Sprintf("expected import string, found %s token instead", &p.peek))
	}

	if p.peek.Kind != token.Semicolon {
		p.report(p.peek.Pin, 0, "missing \";\" after import")
	} else {
		p.advance() // skip ";"
	}

	var path sx.Path
	if s != "" {
		path = sx.MakePath(origin, s)
	}

	p.text.AddImport(ast.Import{
		Path:   path,
		Name:   name,
		Pin:    pin,
		ImpPin: ipin,
	})
	return 0
}
