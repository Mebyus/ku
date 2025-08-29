package opc

// Sys opcode variants.
const (
	// Trap causes immediate abnormal program exit.
	// This instruction is much like halt, but intented
	// to be placed as a trap for code which must not be executed
	// thus causing runtime error.
	//
	// Most notable usecase for this is to detect illegal
	// control flow execution of unreachable code.
	Trap uint8 = iota

	// Halts program execution.
	Halt

	// Reserved empty instruction. Does nop (no operation).
	Nop

	// Call builtin procedure.
	SysCall

	// Return from procedure.
	Ret
)
