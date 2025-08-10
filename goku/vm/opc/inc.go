package opc

// Inc instruction layouts.
//
// Layout of Inc is split into two parts.
// Low 4 bits specify variant and Data encoding.
// High 4 bits carry increment value for IncTiny variant.
const (
	// Increase destination register by immediate value
	// obtained from 4 high bits of instruction layout.
	//
	// Data:
	//	RR - 1 byte - destination
	IncTiny Layout = 0x0

	// Increase destination register by value stored
	// in another register.
	//
	// Data:
	//	RR - 1 byte - destination
	//	RR - 1 byte - source
	IncReg Layout = 0x1

	// Increase destination register by value immediate value
	// from Data.
	//
	// Data:
	//	RR - 1 byte - destination
	//	II II II II - 4 bytes - source
	IncVal32 Layout = 0x2

	// Increase destination register by value immediate value
	// from Data.
	//
	// Data:
	//	RR - 1 byte - destination
	//	II II II II II II II II - 8 bytes - source
	IncVal64 Layout = 0x3
)

func EncodeIncTinyLayout(v uint8) Layout {
	return IncTiny | Layout(v<<4)
}
