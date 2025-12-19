package ast

import "fmt"

func (g *Printer) Type(typ Type) {
	g.puts("type ")
	g.puts(typ.Name.Str)

	if len(typ.Bags) != 0 {
		g.puts(" in ")
		g.bagList(typ.Bags)
	}

	g.space()
	g.TypeSpec(typ.Spec)
}

func (g *Printer) bagList(names []Word) {
	if len(names) == 0 {
		g.puts("()")
		return
	}

	g.puts("(")
	g.puts(names[0].Str)
	for _, name := range names[1:] {
		g.puts(", ")
		g.puts(name.Str)
	}
	g.puts(")")
}

func (g *Printer) TypeSpec(typ TypeSpec) {
	switch t := typ.(type) {
	case nil:
		panic("nil type specifier")
	case TypeName:
		g.TypeName(t)
	case TypeFullName:
		g.TypeFullName(t)
	case Tuple:
		g.Tuple(t)
	case Form:
		g.Form(t)
	case Chunk:
		g.Chunk(t)
	case Array:
		g.Array(t)
	case Struct:
		g.Struct(t)
	case Void:
		g.Void(t)
	case Pointer:
		g.Pointer(t)
	case VoidPointer:
		g.VoidPointer(t)
	case ArrayPointer:
		g.ArrayPointer(t)
	case AnyType:
		g.AnyType(t)
	case Enum:
		g.Enum(t)
	case Bag:
		g.Bag(t)
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", t.Kind(), t.Kind(), t))
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

func (g *Printer) Void(t Void) {
	g.puts("void")
}

func (g *Printer) AnyType(t AnyType) {
	g.puts("type")
}

func (g *Printer) Pointer(p Pointer) {
	g.puts("*")
	g.TypeSpec(p.Type)
}

func (g *Printer) Ref(r Ref) {
	g.puts("&")
	g.TypeSpec(r.Type)
}

func (g *Printer) VoidPointer(p VoidPointer) {
	g.puts("*void")
}

func (g *Printer) ArrayPointer(p ArrayPointer) {
	g.puts("[*]")
	g.TypeSpec(p.Type)
}

func (g *Printer) ArrayRef(r ArrayRef) {
	g.puts("[&]")
	g.TypeSpec(r.Type)
}

func (g *Printer) Chunk(c Chunk) {
	g.puts("[]")
	g.TypeSpec(c.Type)
}

func (g *Printer) Array(a Array) {
	g.puts("[")
	g.Exp(a.Size)
	g.puts("]")
	g.TypeSpec(a.Type)
}

func (g *Printer) TypeName(t TypeName) {
	g.puts(t.Name.Str)
}

func (g *Printer) TypeFullName(t TypeFullName) {
	g.puts(t.Import.Str)
	g.puts(".")
	g.puts(t.Name.Str)
}

func (g *Printer) FunType(f FunType) {
	g.puts("fun ")
	g.Signature(f.Signature)
}

func (g *Printer) Struct(s Struct) {
	g.puts("struct ")
	g.fieldsCurly(s.Fields)
}

func (g *Printer) Union(u Union) {
	g.puts("union ")
	g.fieldsCurly(u.Fields)
}

func (g *Printer) fieldsCurly(fields []Field) {
	if len(fields) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.inc()

	for _, f := range fields {
		g.indent()
		g.Field(f)
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

func (g *Printer) Form(form Form) {
	if len(form.Fields) == 0 {
		g.puts("()")
		return
	}

	g.puts("(")
	g.Field(form.Fields[0])
	for _, f := range form.Fields[1:] {
		g.puts(", ")
		g.Field(f)
	}
	g.puts(")")
}

func (g *Printer) Bag(b Bag) {
	g.puts("bag {")
	if len(b.Funs) == 0 {
		g.puts("}")
		return
	}

	g.inc()
	for _, f := range b.Funs {
		g.nl()
		g.indent()
		g.puts(f.Name.Str)
		g.Signature(f.Signature)
		g.semi()
	}
	g.dec()
	g.nl()
	g.puts("}")
}
