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
		Type: typename("u32"),
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
				param("a", typename("s32")),
				param("b", typename("s32")),
			},
			Result: typename("s32"),
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
		Type: typename("uint"),
		Exp:  hex(0x10),
	})
	t.AddFun(Fun{
		Name: word("foo_bar"),
		Signature: Signature{
			Params: []Param{
				param("a", typename("str")),
			},
			Never: true,
		},
		Body: block(),
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

func call() {}
