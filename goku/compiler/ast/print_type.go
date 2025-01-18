package ast

import "fmt"

func (g *Printer) Type(typ Type) {
	g.puts("type ")
	g.puts(typ.Name.Str)
	g.puts(" => ")
	g.TypeSpec(typ.Spec)
}

func (g *Printer) TypeSpec(typ TypeSpec) {
	switch t := typ.(type) {
	case TypeName:
		g.TypeName(t)
	case TypeFullName:
		g.TypeFullName(t)
	case Tuple:
		g.Tuple(t)
	case Chunk:
		g.Chunk(t)
	case Struct:
		g.Struct(t)
	case Trivial:
		g.Trivial(t)
	case Pointer:
		g.Pointer(t)
	case Enum:
		g.Enum(t)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier", t.Kind(), t.Kind()))
	}
}

func (g *Printer) Enum(e Enum) {
	g.puts(e.Base.Name.Str)
	if len(e.Entries) == 0 {
		g.puts(" {}")
		return
	}

	g.puts(" {")
	g.nl()
	g.inc()

	for _, entry := range e.Entries {
		g.indent()
		g.EnumEntry(entry)
		g.puts(",")
		g.nl()
	}

	g.dec()
	g.puts("}")
}

func (g *Printer) EnumEntry(entry EnumEntry) {
	g.puts(entry.Name.Str)
	if entry.Exp == nil {
		return
	}

	g.puts(" = ")
	g.Exp(entry.Exp)
}

func (g *Printer) Trivial(t Trivial) {
	g.puts("()")
}

func (g *Printer) Pointer(p Pointer) {
	g.puts("*")
	g.TypeSpec(p.Type)
}

func (g *Printer) Chunk(c Chunk) {
	g.puts("[]")
	g.TypeSpec(c.Type)
}

func (g *Printer) TypeName(t TypeName) {
	g.puts(t.Name.Str)
}

func (g *Printer) TypeFullName(t TypeFullName) {
	g.puts(t.Import.Str)
	g.puts(".")
	g.puts(t.Name.Str)
}

func (g *Printer) Struct(s Struct) {
	if len(s.Fields) == 0 {
		g.puts("struct {}")
		return
	}

	g.puts("struct {")
	g.nl()
	g.inc()

	for _, f := range s.Fields {
		g.indent()
		g.Field(f)
		g.puts(",")
		g.nl()
	}

	g.dec()
	g.puts("}")
}

func (g *Printer) Field(f Field) {
	g.puts(f.Name.Str)
	g.puts(": ")
	g.TypeSpec(f.Type)
}

func (g *Printer) Tuple(tuple Tuple) {
	if len(tuple.Types) == 0 {
		g.puts("()")
		return
	}

	g.puts("(")
	g.TypeSpec(tuple.Types[0])
	for _, t := range tuple.Types[1:] {
		g.puts(", ")
		g.TypeSpec(t)
	}
	g.puts(")")
}
