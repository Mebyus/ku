package ast

func (g *Printer) Fun(f Fun) {
	if f.Pub {
		g.puts("pub")
		g.nl()
	}

	g.puts("fun ")

	if f.Unsafe {
		g.puts("unsafe.")
	}

	g.fun(f)
}

func (g *Printer) Test(t Fun) {
	g.puts("test ")
	g.fun(t)
}

func (g *Printer) fun(f Fun) {
	g.puts(f.Name.Str)
	g.Signature(f.Signature)
	g.space()
	g.Block(f.Body)
}

func (g *Printer) FunStub(s FunStub) {
	g.puts("#stub")
	g.nl()
	g.puts("fun ")
	g.puts(s.Name.Str)
	g.Signature(s.Signature)
}

func (g *Printer) Signature(s Signature) {
	g.puts("(")
	g.Params(s.Params)
	g.puts(")")

	if s.Never {
		g.puts(" => #never")
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

func (g *Printer) Method(m Method) {
	if m.Pub {
		g.puts("pub")
		g.nl()
	}

	g.puts("fun ")
	g.puts("(")
	if m.Receiver.Ptr {
		g.puts("*")
	}
	g.puts(m.Receiver.Name.Str)
	g.puts(") ")

	if m.Unsafe {
		g.puts("unsafe.")
	}

	g.puts(m.Name.Str)
	g.Signature(m.Signature)
	g.space()
	g.Block(m.Body)
}
