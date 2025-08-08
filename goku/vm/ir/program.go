package ir

type Program struct {
	// Valid program must have at least one function.
	//
	// One of the functions must also be program entrypoint.
	Functions []Fun

	// List of all data entries defined inside the program.
	Data []DataEntry

	// Index of entrypoint function.
	EntryFun uint32
}

// Fun represents a function inside a program.
// Valid function must contain at least one instruction.
type Fun struct {
	// Atoms constitute function body (code).
	Atoms []Atom

	// List of all labels used inside function body.
	Labels []LabelEntry
}

// Label contains label name in integer form.
//
// Directly corresponds to label entry index inside list of
// all function labels.
type Label uint32

type LabelEntry struct {
	Name Label

	// Atom index which this label points to.
	Atom uint32
}

// Data contains data name in integer form.
//
// Directly corresponds to data entry index inside list of
// all program data.
type Data uint32

type DataEntry struct {
	// Contains raw array of bytes. Using string here to avoid
	// wasting space on redundant capacity.
	//
	// Valid data entry always has non-zero byte length.
	Val string

	Name Data
}
