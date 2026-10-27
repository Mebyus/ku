package vm32b

// Exit describes vm exit state after program execution.
// Includes both normal and abnormal exits.
type Exit struct {
	// Values returned by start function.
	//
	// Valid only for normal exits.
	//
	// Empty if start returns nothing.
	//
	// Callers should not take ownership or modify this slice.
	// Next call to VM.Exec will reuse it.
	Ret []uint32

	// Non-zero code only for abnormal exits.
	Error RuntimeError
}

func (x *Exit) IsOK() bool {
	return x.Error.Code == 0
}

type RuntimeError struct {
	Code ErrorCode
}

type ErrorCode uint32

const (
	// Execution reached end of program text (instructions).
	CodeProgEnd ErrorCode = 1 + iota

	// Executed empty program.
	CodeEmpty

	// Execution reached trap instruction.
	CodeTrap

	// Nil pointer dereference.
	CodeNil

	// Panic inside VM host.
	CodePanic
)
