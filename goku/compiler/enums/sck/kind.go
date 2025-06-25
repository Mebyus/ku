package sck

// Kind indicates scope kind.
type Kind uint8

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Scope with global program-wide access. Holds builtin language symbols.
	Global

	// Scope with unit-wide access (across all unit atoms).
	//
	// Does not include unit tests.
	Unit

	Generic

	// Scope created by unit level function, method or test (inside function body).
	Node

	// Scope that holds collection of all tests inside a unit.
	Test

	// Scope that holds collection of all unsafe nodes inside a unit.
	Unsafe

	Block
	Loop

	// If or else branch block.
	Branch
	Case
)

var text = [...]string{
	empty: "<nil>",

	Global: "global",
	Unit:   "unit",
	Node:   "node",
	Test:   "test",
	Unsafe: "unsafe",
	Block:  "block",
	Loop:   "loop",
	Branch: "branch",
	Case:   "case",
}

func (k Kind) String() string {
	return text[k]
}
