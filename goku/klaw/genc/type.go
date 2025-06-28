package genc

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/ast"
)

func (g *Gen) Type(typ ast.Type) {
	g.puts("typedef ")

	if len(typ.Bags) != 0 {
		panic("not supported")
	}

	f, ok := typ.Spec.(ast.FunType)
	if ok {
		g.typedefFunType(typ.Name.Str, f)
		g.semi()
		return
	}

	g.TypeSpec(typ.Spec)
	g.space()
	g.puts(typ.Name.Str)
	g.semi()
}

func (g *Gen) TypeSpec(spec ast.TypeSpec) {
	switch s := spec.(type) {
	case nil:
		panic("nil type specifier")
	case ast.TypeName:
		g.TypeName(s)
	case ast.TypeFullName:
		panic("not supported")
	case ast.Tuple:
		panic("not supported")
	case ast.Form:
		panic("not supported")
	case ast.Chunk:
		g.Chunk(s)
	case ast.Array:
		panic("not implemented")
	case ast.Struct:
		g.Struct(s)
	case ast.Trivial:
		panic("not supported")
	case ast.Pointer:
		g.Pointer(s)
	case ast.AnyPointer:
		g.AnyPointer(s)
	case ast.ArrayPointer:
		g.ArrayPointer(s)
	case ast.AnyType:
		panic("not implemented")
	case ast.Enum:
		panic("not implemented")
	case ast.Bag:
		panic("not implemented")
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type specifier (%T)", s.Kind(), s.Kind(), s))
	}
}

func (g *Gen) TypeName(t ast.TypeName) {
	g.puts(t.Name.Str)
}

func (g *Gen) AnyPointer(p ast.AnyPointer) {
	g.puts("void*")
}

func (g *Gen) Pointer(p ast.Pointer) {
	g.TypeSpec(p.Type)
	g.putb('*')
}

func (g *Gen) ArrayPointer(p ast.ArrayPointer) {
	g.TypeSpec(p.Type)
	g.putb('*')
}

func (g *Gen) Chunk(c ast.Chunk) {
	t, ok := c.Type.(ast.TypeName)
	if !ok {
		panic("not implemented")
	}
	g.puts("span_")
	g.puts(t.Name.Str)
}

func (g *Gen) Struct(s ast.Struct) {
	g.puts("struct ")
	if len(s.Fields) == 0 {
		g.puts("{}")
		return
	}

	g.puts("{")
	g.nl()
	g.inc()

	for _, f := range s.Fields {
		g.indent()
		g.Field(f)
		g.puts(";")
		g.nl()
	}

	g.dec()
	g.puts("}")
}

func (g *Gen) Field(f ast.Field) {
	g.NameDef(f.Name.Str, f.Type)
}

// NameDef formats declaration of variable, constant or function parameter
// according to C syntax. Namely it handles C arrays declarations.
func (g *Gen) NameDef(name string, spec ast.TypeSpec) {
	a, ok := spec.(ast.Array)
	if ok {
		g.TypeSpec(a.Type)
		g.space()
		g.puts(name)
		g.putb('[')
		g.Exp(a.Size)
		g.putb(']')
		return
	}

	g.TypeSpec(spec)
	g.space()
	g.puts(name)
}

func (g *Gen) typedefFunType(name string, f ast.FunType) {
	if f.Signature.Result == nil {
		g.puts("void")
	} else {
		g.TypeSpec(f.Signature.Result)
	}

	g.puts(" (*")
	g.puts(name)
	g.puts(")")

	if len(f.Params) == 0 {
		g.puts("(void)")
		return
	}

	g.puts("(")
	g.TypeSpec(f.Params[0].Type)
	for _, p := range f.Params[1:] {
		g.puts(", ")
		g.TypeSpec(p.Type)
	}
	g.puts(")")
}
