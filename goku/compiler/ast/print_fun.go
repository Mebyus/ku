package ast

func (g *Printer) Fun(f Fun) {
	g.puts("fun ")
	g.puts(f.Name.Str)
	g.Signature(f.Signature)
	g.space()
	g.Block(f.Body)
}

func (g *Printer) Test(t Fun) {

}

func (g *Printer) FunStub(s FunStub) {

}

func (g *Printer) Signature(s Signature) {
	g.puts("(")
	g.Params(s.Params)
	g.puts(")")

	if s.Never {
		g.puts(" => never")
		return
	}

	if s.Result == nil {
		return
	}

	g.puts(" => ")
	g.TypeSpec(s.Result)
}

func (g *Printer) Params(params []Param) {
	if len(params) == 0 {
		g.puts("()")
		return
	}

	g.Param(params[0])
	for _, p := range params[1:] {
		g.puts(", ")
		g.Param(p)
	}
}

func (g *Printer) Param(p Param) {
	g.puts(p.Name.Str)
	g.puts(": ")
	g.TypeSpec(p.Type)
}
