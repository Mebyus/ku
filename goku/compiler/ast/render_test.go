package ast

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mebyus/ku/goku/compiler/enums/bok"
)

type RenderTestCase struct {
	Text *Text
	File string
}

func text1() *Text {
	t := New()
	t.AddVar(Var{
		Name: word("foo"),
		Type: u32,
		Exp:  dec(42),
	})
	return t
}

func text2() *Text {
	t := New()
	t.AddFun(Fun{
		Name: word("add"),
		Signature: Signature{
			Params: []Param{
				param("a", s32),
				param("b", s32),
			},
			Result: s32,
		},
		Body: block(
			ret(add(sym("a"), sym("b"))),
		),
	})
	return t
}

func text3() *Text {
	t := New()
	t.AddLet(Let{
		Name: word("x"),
		Type: u64,
		Exp:  hex(0x10),
	})
	t.AddFun(Fun{
		Name: word("foo_bar"),
		Signature: Signature{
			Params: []Param{
				param("a", str),
			},
			Never: true,
		},
		Body: block(
			invoke(chain("print"), sym("a")),
			invoke(chain("crash")),
		),
	})
	return t
}

func text4() *Text {
	t := New()
	t.AddFun(Fun{
		Name: word("at_zero"),
		Signature: Signature{
			Params: []Param{
				param("a", chunk(u64)),
			},
			Result: u64,
		},
		Body: block(
			let("i", u32, dec(0)),
			ret(chain("a", index(sym("i")))),
		),
	})
	return t
}

func text5() *Text {
	t := New()
	t.AddType(Type{
		Name: word("Foo"),
		Spec: typestruct(field("bar", str)),
	})
	t.AddType(Type{
		Name: word("Bar"),
		Spec: str,
	})
	t.AddType(Type{
		Name: word("Hello"),
		Spec: triv,
	})
	// t.AddType(Type{})
	t.AddFun(Fun{
		Name: word("modify"),
		Signature: Signature{
			Params: []Param{
				param("foo", ptr(typename("Foo"))),
			},
		},
	})
	return t
}

func prepareRenderTestCases() []RenderTestCase {
	return []RenderTestCase{
		{
			File: "00000.ku",
			Text: New(),
		},
		{
			File: "00001.ku",
			Text: text1(),
		},
		{
			File: "00002.ku",
			Text: text2(),
		},
		{
			File: "00003.ku",
			Text: text3(),
		},
		{
			File: "00004.ku",
			Text: text4(),
		},
		{
			File: "00005.ku",
			Text: text5(),
		},
	}
}

func TestRender(t *testing.T) {
	tests := prepareRenderTestCases()
	for _, tt := range tests {
		t.Run(tt.File, func(t *testing.T) {
			wantBytes, err := os.ReadFile(filepath.Join("testdata", tt.File))
			if err != nil {
				t.Error(err)
				return
			}

			var p Printer
			p.Text(tt.Text)
			gotBytes := p.Bytes()
			if !bytes.Equal(gotBytes, wantBytes) {
				t.Errorf("\n========================================\n%s\n========================================\n%s\n========================================", gotBytes, wantBytes)
			}
		})
	}
}

func word(s string) Word {
	return Word{Str: s}
}

func typename(s string) TypeName {
	return TypeName{Name: word(s)}
}

func field(name string, typ TypeSpec) Field {
	return Field{
		Name: word(name),
		Type: typ,
	}
}

func typestruct(fields ...Field) Struct {
	return Struct{Fields: fields}
}

func ptr(typ TypeSpec) Pointer {
	return Pointer{Type: typ}
}

func chunk(typ TypeSpec) Chunk {
	return Chunk{Type: typ}
}

func dec(n uint64) Integer {
	return Integer{Val: n, Aux: uint32(IntDec)}
}

func hex(n uint64) Integer {
	return Integer{Val: n, Aux: uint32(IntHex)}
}

func param(name string, typ TypeSpec) Param {
	return Param{
		Name: word(name),
		Type: typ,
	}
}

func block(nodes ...Statement) Block {
	return Block{Nodes: nodes}
}

func ret(exp Exp) Ret {
	return Ret{Exp: exp}
}

func let(name string, typ TypeSpec, exp Exp) Let {
	return Let{
		Name: word(name),
		Type: typ,
		Exp:  exp,
	}
}

func sym(name string) Symbol {
	return Symbol{Name: name}
}

func bin(op bok.Kind, a Exp, b Exp) Binary {
	return Binary{
		Op: BinOp{Kind: op},
		A:  a,
		B:  b,
	}
}

func add(a Exp, b Exp) Binary {
	return bin(bok.Add, a, b)
}

func chain(start string, parts ...Part) Chain {
	return Chain{
		Start: word(start),
		Parts: parts,
	}
}

func index(exp Exp) Part {
	return Index{Exp: exp}
}

func invoke(chain Chain, args ...Exp) Invoke {
	return Invoke{
		Call: Call{
			Chain: chain,
			Args:  args,
		},
	}
}

var (
	u32 = typename("u32")
	u64 = typename("u64")

	s32 = typename("s32")

	str = typename("str")
)

var triv = Trivial{}
