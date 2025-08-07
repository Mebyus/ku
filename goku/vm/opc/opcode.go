package opc

// Opcode denotes which type of operation instruction performs
// inside VM.
type Opcode uint8

const (
	// Trap causes immediate abnormal program exit.
	// This instruction is much like halt, but intented
	// to be placed as a trap for code which must not be executed
	// thus causing runtime error.
	//
	// Most notable usecase for this is to detect illegal
	// control flow execution of unreachable code.
	Trap Opcode = 0x00

	// Halts program execution.
	Halt Opcode = 0x01

	// Reserved empty instruction. Does nop (no operation).
	Nop Opcode = 0x02

	// Call builtin procedure.
	SysCall Opcode = 0x03

	// Change instruction pointer via jump in text.
	Jump Opcode = 0x04

	// Call procedure stored in text.
	Call Opcode = 0x05

	// Return from procedure.
	Ret Opcode = 0x06

	// Push value to stack memory.
	Push Opcode = 0x07

	// Pop value from stack memory.
	Pop Opcode = 0x08

	// Set register or memory location to zero.
	Clear Opcode = 0x09

	// Copy register value or immediate value to another register.
	Copy Opcode = 0x0A

	// Load value to register from memory location.
	Load Opcode = 0x0B

	// Store value from register or immediate value to memory location.
	Store Opcode = 0x0C

	// Compare two values.
	Test Opcode = 0x0D

	// Increase value in register.
	Inc Opcode = 0x0E

	// Decrease value in register.
	Dec Opcode = 0x0F

	// Add two values.
	Add Opcode = 0x10

	// Subtract two values.
	Sub Opcode = 0x11
)

var opcodeText = [...]string{
	Trap:    "trap",
	Halt:    "halt",
	Nop:     "nop",
	SysCall: "syscall",
	Jump:    "jump",
	Call:    "call",
	Ret:     "ret",
	Push:    "push",
	Pop:     "pop",
	Clear:   "clear",
	Copy:    "copy",
	Load:    "load",
	Store:   "store",
	Test:    "test",
	Inc:     "inc",
	Add:     "add",
	Sub:     "sub",
}

func (c Opcode) String() string {
	return opcodeText[c]
}
