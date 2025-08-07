package opc

// Layout specifies instruction subkind as well as how instruction data
// is encoded in text. For example in jump instruction layout determines
// jump condition (if any) and where destination is stored.
type Layout uint8

// Generic layouts.
const (
	// Data:
	//	RR - 1 byte - destination
	//	RR - 1 byte - source
	//	RR - 1 byte - source
	RegRegReg Layout = 0x0

	// Data:
	//	RR - 1 byte - destination
	//	RR - 1 byte - source
	//	II II II II - 4 bytes - source
	RegRegVal32 Layout = 0x1

	// Data:
	//	RR - 1 byte - destination
	//	RR - 1 byte - source
	//	II II II II II II II II - 8 bytes - source
	RegRegVal64 Layout = 0x2
)
