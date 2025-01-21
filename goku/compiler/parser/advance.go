package parser

func (p *Parser) advance() {
	p.c = p.n
	p.n = p.lx.Lex()
}

func (p *Parser) init() {
	p.advance()
	p.advance()
}
