package tnk

// Kind indicates AST top level node kind.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	Var
	Const
	Type
	Alias

	Test
	Fun
	Method
	FunStub

	Gen
	GenBind
	Lookup
)

var text = [...]string{
	empty: "<nil>",

	Fun:   "fun",
	Var:   "var",
	Const: "const",
	Type:  "type",
	Alias: "alias",

	Test:    "test",
	Method:  "method",
	FunStub: "fun.stub",

	Gen:     "gen",
	GenBind: "gen.bind",
	Lookup:  "lookup",
}

func (k Kind) String() string {
	return text[k]
}
