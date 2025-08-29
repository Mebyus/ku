package opc

// Opcode denotes which type of operation instruction performs
// inside VM.
type Opcode uint8

const (
	// This instruction family serves as collection for various
	// VM control operations. All instructions in this family have no
	// data part.
	Sys Opcode = iota

	// Change instruction pointer via jump in text.
	Jump

	// Call procedure stored in text.
	Call

	// Push value to stack memory.
	Push

	// Pop value from stack memory.
	Pop

	// Set register or memory location to zero.
	Clear

	// Copy register value or immediate value to another register.
	Set

	// Load value to register from memory location.
	Load

	// Store value from register or immediate value to memory location.
	Store

	// Compare two values.
	Test

	// Increase value in register.
	Inc

	// Decrease value in register.
	Dec

	// Add two values.
	Add

	// Subtract two values.
	Sub
)

var opcodeText = [...]string{
	Sys:   "sys",
	Jump:  "jump",
	Call:  "call",
	Push:  "push",
	Pop:   "pop",
	Clear: "clear",
	Set:   "set",
	Load:  "load",
	Store: "store",
	Test:  "test",
	Inc:   "inc",
	Add:   "add",
	Sub:   "sub",
}

func (c Opcode) String() string {
	return opcodeText[c]
}
