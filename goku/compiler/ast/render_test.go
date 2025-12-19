package ast

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mebyus/ku/goku/compiler/enums/aok"
	"github.com/mebyus/ku/goku/compiler/enums/bok"
	"github.com/mebyus/ku/goku/compiler/enums/uok"
	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
)

type RenderTestCase struct {
	Text *Text
	File string
}

func text1() *Text {
	t := New()
	t.AddVar(TopVar{
		Var: Var{
			Name: word("foo"),
			Type: u32,
			Exp:  dec(42),
		},
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
	t.AddConst(TopConst{
		Const: Const{
			Name: word("x"),
			Type: u64,
			Exp:  hex(0x10),
		},
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
			constdef("i", u32, dec(0)),
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
			assign(
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
			Kind: ReceiverPtr,
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
			Kind: ReceiverVal,
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

func text7() *Text {
	t := New()
	t.AddType(Type{
		Name: word("CustomError"),
		Spec: enum("erval",
			ee("Error1", dec(1)),
			ee("Error2", dec(2)),
			ee("Error3", nil),
		),
	})
	t.AddFun(Fun{
		Name: word("count"),
		Signature: Signature{
			Params: nil,
			Result: s32,
		},
		Body: block(
			vardef("i", s32, dec(0)),
			loop(
				ifbr(
					ifcl(bin(bok.Greater, sym("i"), dec(10)),
						block(
							ret(sym("i"))),
					),
					nil,
				),
				assignAdd(sym("i"), dec(1)),
			),
			ret(unary(uok.Minus, hex(0x3))),
		),
	})
	return t
}

func text8() *Text {
	t := New()
	t.ImportBlocks = []ImportBlock{
		importBlock(origin.Std,
			imp("mem", "mem"),
			imp("json", "json"),
		),
		importBlock(origin.Pkg,
			imp("pk", "another/person/package"),
		),
		importBlock(origin.Loc,
			imp("bar", "foo/bar"),
			imp("hello", "example/hello"),
		),
	}
	return t
}

func text9() *Text {
	t := New()
	t.AddFun(Fun{
		Name: word("return_one"),
		Signature: Signature{
			Result: u32,
		},
		Body: block(
			vardef("x", u32, Dirty{}),
			assign(sym("x"), dec(1)),
			ret(sym("x")),
		),
	})
	return t
}

func text10() *Text {
	t := New()
	t.AddTest(TestFun{
		Name: word("inc"),
		Body: block(
			ret(nil),
		),
	})
	return t
}

func text11() *Text {
	t := New()
	t.AddStub(FunStub{
		Name: word("example"),
		Signature: Signature{
			Params: []Param{
				param("a", u32),
				param("b", u32),
			},
			Result: tuple(u32, typename("bool")),
		},
	})
	return t
}

func text12() *Text {
	t := New()
	t.Build = build(
		ifbr(
			ifcl(eq(chain("g", sel("target"), sel("os")), dotname("WINDOWS")),
				block(
					invoke(chain("g", sel("skip"))),
				)),
			nil,
		),
	)
	t.ImportBlocks = []ImportBlock{
		importBlock(origin.Std,
			imp("win", "os/windows"),
		),
	}
	return t
}

func text13() *Text {
	t := New()
	t.AddFun(Fun{
		Name: word("copy"),
		Traits: Traits{
			Pub:    true,
			Unsafe: true,
		},
		Signature: Signature{
			Params: []Param{
				param("dst", arrptr(u8)),
				param("src", arrptr(u8)),
				param("n", typename("uint")),
			},
		},
		Body: block(
			Stub{},
		),
	})
	return t
}

func text14() *Text {
	t := New()
	t.ImportBlocks = []ImportBlock{
		importBlock(origin.Std,
			imp("mem", "mem"),
		),
	}
	t.AddFun(Fun{
		Name: word("use_copy"),
		Signature: Signature{
			Params: []Param{
				param("a", chunk(u64)),
			},
			Result: u64,
		},
		Body: block(
			ifbr(
				ifcl(eq(chain("a", sel("len")), dec(0)),
					block(Never{})),
				nil),
			vardef("i", u32, nil),
			invoke(unsafechain("copy"), sym("a"), Nil{}),
			invoke(chain("mem", unsafe("copy")), sym("a"), Nil{}),
			ret(chain("a", index(sym("i")))),
		),
	})
	return t
}

func text15() *Text {
	t := New()
	t.AddFun(Fun{
		Name: word("example_debug"),
		Signature: Signature{
			Params: []Param{
				param("a", str),
			},
		},
		Body: block(
			debug(
				invoke(chain("print"), slit("debug: ")),
				invoke(chain("print"), sym("a")),
				ifbr(
					ifcl(great(chain("a", sel("len")), dec(0)),
						block(ret(nil))),
					nil,
				),
			),
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
		{
			File: "00007.ku",
			Text: text7(),
		},
		{
			File: "00008.ku",
			Text: text8(),
		},
		{
			File: "00009.ku",
			Text: text9(),
		},
		{
			File: "00010.ku",
			Text: text10(),
		},
		{
			File: "00011.ku",
			Text: text11(),
		},
		{
			File: "00012.ku",
			Text: text12(),
		},
		{
			File: "00013.ku",
			Text: text13(),
		},
		{
			File: "00014.ku",
			Text: text14(),
		},
		{
			File: "00015.ku",
			Text: text15(),
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

func arrptr(typ TypeSpec) ArrayPointer {
	return ArrayPointer{Type: typ}
}

func chunk(typ TypeSpec) Span {
	return Span{Type: typ}
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

func loop(nodes ...Statement) Loop {
	return Loop{Body: block(nodes...)}
}

func ifcl(exp Exp, body Block) IfClause {
	return IfClause{
		Exp:  exp,
		Body: body,
	}
}

func ifbr(ifclause IfClause, elseBody *Block, elseifs ...IfClause) If {
	return If{
		If:      ifclause,
		ElseIfs: elseifs,
		Else:    elseBody,
	}
}

func ret(exp Exp) Ret {
	return Ret{Exp: exp}
}

func constdef(name string, typ TypeSpec, exp Exp) Const {
	return Const{
		Name: word(name),
		Type: typ,
		Exp:  exp,
	}
}

func vardef(name string, typ TypeSpec, exp Exp) Var {
	return Var{
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

func unary(op uok.Kind, exp Exp) Unary {
	return Unary{
		Op:  UnaryOp{Kind: op},
		Exp: exp,
	}
}

func imp(name string, s string) Import {
	return Import{
		Name:   word(name),
		String: ImportString{Str: s},
	}
}

func importBlock(o origin.Origin, imports ...Import) ImportBlock {
	return ImportBlock{
		Imports: imports,
		Origin:  o,
	}
}

func build(nodes ...Statement) *Build {
	return &Build{Body: block(nodes...)}
}

func debug(nodes ...Statement) Debug {
	return Debug{Block: block(nodes...)}
}

func add(a Exp, b Exp) Binary {
	return bin(bok.Add, a, b)
}

func eq(a Exp, b Exp) Binary {
	return bin(bok.Equal, a, b)
}

func great(a Exp, b Exp) Binary {
	return bin(bok.Greater, a, b)
}

func dotname(name string) DotName {
	return DotName{Name: name}
}

func chain(start string, parts ...Part) Chain {
	return Chain{
		Start: word(start),
		Parts: parts,
	}
}

func unsafechain(name string, parts ...Part) Chain {
	s := []Part{Unsafe{Name: name}}
	return chain("", append(s, parts...)...)
}

func unsafe(name string) Unsafe {
	return Unsafe{Name: name}
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

func assign(target Exp, value Exp) Assign {
	return Assign{
		Op:     AssignOp{Kind: aok.Simple},
		Target: target,
		Value:  value,
	}
}

func assignAdd(target Exp, value Exp) Assign {
	return Assign{
		Op:     AssignOp{Kind: aok.Add},
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

var triv = Void{}
