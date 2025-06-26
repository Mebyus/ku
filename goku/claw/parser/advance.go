package parser

func (p *Parser) advance() {
	p.peek = p.next
	p.next = p.lx.Lex()
}

func (p *Parser) init() {
	p.advance()
	p.advance()
}
