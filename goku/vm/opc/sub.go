package opc

// Test instruction layouts.
//
// High 4 bits encode actual layout.
// Low 4 bits may carry low bits of immediate value for some layouts.
const (
	// Compare destination register value to source register.
	//
	// Data:
	//	RR - 1 byte - destination
	//	RR - 1 byte - source
	TestReg Layout = iota

	// Compare destination register to immediate value encoded in layout.
	//
	// Data:
	//	RR - 1 byte - destination
	TestVal4

	// Compare destination register to immediate value encoded in data.
	//
	// Data:
	//	RR - 1 byte - destination
	//	II - 1 byte - source
	TestVal8

	// Data:
	//	RR - 1 byte - destination
	//	II - 2 byte - source
	TestVal16

	// Data:
	//	RR - 1 byte - destination
	//	II - 4 byte - source
	TestVal32

	// Data:
	//	RR - 1 byte - destination
	//	II - 8 byte - source
	TestVal64
)

func EncodeTestValLayout(lt Layout, v uint8) Layout {
	return (lt << 4) | Layout(v)
}

func EncodeTestLayout(lt Layout) Layout {
	return EncodeTestValLayout(lt, 0)
}

func DecodeTestLayout(x uint8) (Layout, uint8) {
	layout := Layout(x >> 4)
	value := x & 0xF
	return layout, value
}
