package bgk

// Kind indicates builtin generic kind.
type Kind uint8

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Min function.
	Min

	// Max function.
	Max

	// Copy function.
	Copy

	// Clear function.
	Clear
)
