package stk

// Kind indicates statement kind.
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	Block
	Assign
	Asm
	Ret
	Const
	Alias
	Var
	If
	Invoke
	Loop
	While
	ForRange
	Match
	MatchBool
	Goto
	Never
	Stub
	Lookup
	Debug

	Break
	Gonext

	Panic
	Must
	Test
	StaticMust
	StaticIf

	// DeferCall function or method call.
	DeferCall

	maxKind
)

var text = [...]string{
	empty: "<nil>",

	Block:     "block",
	Assign:    "assign",
	Ret:       "ret",
	Asm:       "asm",
	Const:     "const",
	Alias:     "alias",
	Var:       "var",
	If:        "if",
	Invoke:    "invoke",
	Loop:      "loop",
	While:     "while",
	ForRange:  "for.range",
	Match:     "match",
	MatchBool: "match.bool",
	DeferCall: "defer call",
	Goto:      "goto",
	Never:     "never",
	Stub:      "stub",
	Debug:     "debug",
	Lookup:    "lookup",
	Break:     "break",
	Gonext:    "gonext",

	Must:       "must",
	Test:       "test",
	Panic:      "panic",
	StaticMust: "must.static",
	StaticIf:   "if.static",
}

func (k Kind) String() string {
	return text[k]
}
