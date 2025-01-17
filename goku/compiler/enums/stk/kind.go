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
	Let
	Var
	If
	Invoke
	Loop
	While
	ForRange
	Match
	MatchBool
	Jump
	Never
	Stub

	// Defer function or method call.
	Defer

	maxKind
)

var text = [...]string{
	empty: "<nil>",

	Block:     "block",
	Assign:    "assign",
	Ret:       "ret",
	Let:       "let",
	Var:       "var",
	If:        "if",
	Invoke:    "invoke",
	Loop:      "loop",
	While:     "while",
	ForRange:  "for.range",
	Match:     "match",
	MatchBool: "match.bool",
	Defer:     "defer",
	Jump:      "jump",
	Never:     "never",
	Stub:      "stub",
}

func (k Kind) String() string {
	return text[k]
}
