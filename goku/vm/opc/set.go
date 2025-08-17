package opc

// Set instruction layouts.
//
// High 4 bits encode actual layout.
// Low 4 bits may carry low bits of immediate value for some layouts.
const (
	// Copy value from source register to destination register.
	//
	// Data:
	//	RR - 1 byte - destination
	//	RR - 1 byte - source
	SetReg Layout = iota

	// Set destination register to immediate value encoded in layout.
	//
	// Data:
	//	RR - 1 byte - destination
	SetVal4

	// Set destination register to immediate value encoded in data.
	//
	// Data:
	//	RR - 1 byte - destination
	//	II - 1 byte - source
	SetVal8

	// Data:
	//	RR - 1 byte - destination
	//	II - 2 byte - source
	SetVal16

	// Data:
	//	RR - 1 byte - destination
	//	II - 4 byte - source
	SetVal32

	// Data:
	//	RR - 1 byte - destination
	//	II - 8 byte - source
	SetVal64
)

func EncodeSetValLayout(lt Layout, v uint8) Layout {
	return (lt << 4) | Layout(v)
}

func EncodeSetLayout(lt Layout) Layout {
	return EncodeSetValLayout(lt, 0)
}

func DecodeSetLayout(x uint8) (Layout, uint8) {
	layout := Layout(x >> 4)
	value := x & 0xF
	return layout, value
}
