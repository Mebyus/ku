package genc

import (
	"fmt"
	"strings"

	"github.com/mebyus/ku/goku/compiler/ast"
)

func (g *Gen) Type(typ ast.Type) {
	if len(typ.Bags) != 0 {
		panic("not supported")
	}

	switch s := typ.Spec.(type) {
	case ast.FunType:
		g.typedefFunType(typ.Name.Str, s)
		return
	case ast.Enum:
		g.typedefEnumType(typ.Name.Str, s)
		return
	case ast.Bag:
		g.typedefBagType(typ.Name.Str, s)
		return
	}

	g.puts("typedef ")
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
	case ast.Union:
		g.Union(s)
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
		panic("not supported")
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

func (g *Gen) typedefEnumType(name string, enum ast.Enum) {
	g.puts("typedef ")
	g.puts(enum.Base.Name.Str)
	g.space()
	g.puts(name)
	g.semi()

	if len(enum.Entries) == 0 {
		return
	}

	g.nl()
	g.puts("enum {")
	g.nl()
	g.inc()
	for _, e := range enum.Entries {
		g.indent()
		g.puts(e.Name.Str)
		g.puts(" = ")
		g.Exp(e.Exp)
		g.putb(',')
		g.nl()
	}
	g.dec()
	g.puts("};")
}

func (g *Gen) Struct(s ast.Struct) {
	g.puts("struct ")
	g.fieldsCurly(s.Fields)
}

func (g *Gen) Union(u ast.Union) {
	g.puts("union ")
	g.fieldsCurly(u.Fields)
}

func (g *Gen) fieldsCurly(fields []ast.Field) {
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
	g.puts("typedef ")

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
	g.puts(");")
}

func (g *Gen) typedefBagFun(f ast.BagFun) {
	g.puts("typedef ")

	s := f.Signature
	if s.Result == nil {
		g.puts("void")
	} else {
		g.TypeSpec(s.Result)
	}

	g.puts(" (*BagFun")
	g.puts(strings.Title(f.Name.Str)) // TODO: make siplified transform for uppercase letter for function name
	g.puts(")")

	g.puts("(uint")
	for _, p := range f.Signature.Params {
		g.puts(", ")
		g.TypeSpec(p.Type)
	}
	g.puts(");")
}

func (g *Gen) typedefBagType(name string, b ast.Bag) {
	if len(b.Funs) == 0 {
		panic("empty bag type")
	}

	for _, f := range b.Funs {
		g.typedefBagFun(f)
	}

	g.nl()
	g.puts("typedef struct {")
	g.nl()
	g.inc()

	for _, f := range b.Funs {
		g.indent()
		g.puts("BagFun")
		g.puts(strings.Title(f.Name.Str)) // TODO: maybe we should optimize function type name generation
		g.space()
		g.puts(f.Name.Str)
		g.semi()
		g.nl()
	}

	g.dec()
	g.puts("} Bag")
	g.puts(name)
	g.puts("Tab;")

	g.nl()
	g.puts("typedef struct {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("uint obj;")
	g.nl()

	g.indent()
	g.puts("Bag")
	g.puts(name)
	g.puts("Tab* tab;")
	g.nl()

	g.dec()
	g.puts("} ")
	g.puts(name)
	g.semi()
}
