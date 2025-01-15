package parser

func (p *Parser) advance() {
	p.c = p.n
	p.n = p.s.Lex()
}

func (p *Parser) init() {
	p.advance()
	p.advance()
}
