package ast

import "fmt"

func (g *Printer) Type(typ Type) {
	g.puts("type ")
	g.puts(typ.Name.Str)
	g.space()
	g.TypeSpec(typ.Spec)
}

func (g *Printer) TypeSpec(spec TypeSpec) {
	switch s := spec.(type) {
	case TypeName:
		g.TypeName(s)
	case TypeFullName:
		g.TypeFullName(s)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier", s.Kind(), s.Kind()))
	}
}

func (g *Printer) TypeName(name TypeName) {
	g.puts(name.Name.Str)
}

func (g *Printer) TypeFullName(name TypeFullName) {
	g.puts(name.Import.Str)
	g.puts(".")
	g.puts(name.Name.Str)
}
