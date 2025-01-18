package ast

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mebyus/ku/goku/compiler/enums/aok"
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
	t.AddType(Type{
		Name: word("FooError"),
		Spec: enum("u32",
			ee("BIG", dec(1)),
			ee("SMALL", dec(2)),
		),
	})
	t.AddFun(Fun{
		Name: word("modify"),
		Signature: Signature{
			Params: []Param{
				param("foo", ptr(typename("Foo"))),
			},
		},
		Body: block(
			asimp(
				chain("foo", sel("bar")),
				slit("hello, foo!"),
			),
		),
	})
	return t
}

func text6() *Text {
	t := New()
	t.AddMethod(Method{
		Receiver: Receiver{
			Name: word("Foo"),
			Ptr:  true,
		},
		Name: word("write"),
		Signature: Signature{
			Params: []Param{
				param("data", chunk(u8)),
			},
			Result: tuple(u64, typename("error")),
		},
		Body: block(
			ret(pack(chain("data", sel("len")), Nil{})),
		),
	})
	t.AddMethod(Method{
		Receiver: Receiver{
			Name: word("Foo"),
		},
		Name: word("string"),
		Signature: Signature{
			Params: []Param{
				param("al", full("mem", "Allocator")),
			},
			Result: str,
		},
		Body: block(
			ret(slit("")),
		),
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
		{
			File: "00006.ku",
			Text: text6(),
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

func full(imp string, name string) TypeFullName {
	return TypeFullName{
		Import: word(imp),
		Name:   word(name),
	}
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

func ee(name string, exp Exp) EnumEntry {
	return EnumEntry{
		Name: word(name),
		Exp:  exp,
	}
}

func enum(base string, entries ...EnumEntry) Enum {
	return Enum{
		Base:    typename(base),
		Entries: entries,
	}
}

func ptr(typ TypeSpec) Pointer {
	return Pointer{Type: typ}
}

func chunk(typ TypeSpec) Chunk {
	return Chunk{Type: typ}
}

func slit(s string) String {
	return String{Val: s}
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

func index(exp Exp) Index {
	return Index{Exp: exp}
}

func sel(name string) Select {
	return Select{Name: word(name)}
}

func invoke(chain Chain, args ...Exp) Invoke {
	return Invoke{
		Call: Call{
			Chain: chain,
			Args:  args,
		},
	}
}

func asimp(target Exp, value Exp) Assign {
	return Assign{
		Op:     AssignOp{Kind: aok.Simple},
		Target: target,
		Value:  value,
	}
}

func tuple(typs ...TypeSpec) Tuple {
	return Tuple{Types: typs}
}

func pack(exps ...Exp) Pack {
	if len(exps) < 2 {
		panic("not enough elements to form a pack")
	}
	return Pack{List: exps}
}

var (
	u8  = typename("u8")
	u32 = typename("u32")
	u64 = typename("u64")

	s32 = typename("s32")

	str = typename("str")
)

var triv = Trivial{}
