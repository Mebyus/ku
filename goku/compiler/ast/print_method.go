package ast

func (g *Printer) Method(m Method) {
	g.puts("fun ")
	g.puts("(")
	if m.Receiver.Ptr {
		g.puts("*")
	}
	g.puts(m.Receiver.Name.Str)
	g.puts(") ")
	g.puts(m.Name.Str)
	g.Signature(m.Signature)
	g.Block(m.Body)
}
