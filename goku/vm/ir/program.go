package ir

type Program struct {
	// Valid program must have at least one function.
	//
	// One of the functions must also be program entrypoint.
	Functions []Fun

	// List of all data entries defined inside the program.
	Data []DataEntry

	// Integer name (index) of entrypoint function.
	EntryFun FunName

	// Total number of distinct labels inside the program.
	LabelsCount uint32
}

// Fun represents a function inside a program.
// Valid function must contain at least one instruction.
type Fun struct {
	// Atoms constitute function body (code).
	Atoms []Atom

	Name FunName
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
