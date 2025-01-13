package tnk

// Kind indicates AST top level node kind.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	Fun
	Let
	Var
	Type
	Test
	Method
	FunStub
)

var text = [...]string{
	empty: "<nil>",

	Fun: "fun",
	Let: "let",
	Var: "var",

	Type:   "type",
	Test:   "test",
	Method: "method",

	FunStub: "fun.stub",
}

func (k Kind) String() string {
	return text[k]
}
