package opc

// Inc instruction layouts.
const (
	// Increase value in specified register by 1.
	//
	// Data:
	//	RR - 1 byte - destination
	IncReg Layout = 0x0

	// Increase value in specified register by value
	// stored in another register.
	//
	// Data:
	//	RR - 1 byte - destination
	//	RR - 1 byte - source
	IncRegReg Layout = 0x1

	// Increase value in specified register by immediate value.
	//
	// Data:
	//	RR - 1 byte - destination
	//	II II II II - 4 bytes - source
	IncRegVal32 Layout = 0x2

	// Increase value in specified register by immediate value.
	//
	// Data:
	//	RR - 1 byte - destination
	//	II II II II II II II II - 8 bytes - source
	IncRegVal64 Layout = 0x3
)
