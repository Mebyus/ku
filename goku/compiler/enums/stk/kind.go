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
	Ret
	Const
	Alias
	Var
	If
	Invoke
	Loop
	While
	ForRange
	ForRange2
	Match
	MatchBool
	Jump
	Never
	Stub
	Lookup
	Debug

	JumpNext
	JumpOut

	Panic
	Must
	Test
	StaticMust

	// Defer function or method call.
	Defer

	maxKind
)

var text = [...]string{
	empty: "<nil>",

	Block:     "block",
	Assign:    "assign",
	Ret:       "ret",
	Const:     "const",
	Alias:     "alias",
	Var:       "var",
	If:        "if",
	Invoke:    "invoke",
	Loop:      "loop",
	While:     "while",
	ForRange:  "for.range",
	ForRange2: "for.range2",
	Match:     "match",
	MatchBool: "match.bool",
	Defer:     "defer",
	Jump:      "jump",
	Never:     "never",
	Stub:      "stub",
	Debug:     "debug",
	Lookup:    "lookup",

	JumpNext: "jump.next",
	JumpOut:  "jump.out",

	Must:       "must",
	Test:       "test",
	Panic:      "panic",
	StaticMust: "must.static",
}

func (k Kind) String() string {
	return text[k]
}
