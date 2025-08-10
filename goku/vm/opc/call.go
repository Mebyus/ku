package opc

// Call instruction layouts.
const (
	// Call to function which address is stored in register.
	// Address must be in text segment.
	//
	// Data:
	//	RR - 1 byte - destination
	CallReg Layout = 0x0

	// Call to function which address is stored in instruction
	// immediate value.
	// Value is extended to text segment address.
	//
	// Data:
	//	II II II II - 4 bytes - destination
	CallVal32 Layout = 0x1
)
