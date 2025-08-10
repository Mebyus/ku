package opc

import "fmt"

// Flag specifies jump condition.
type Flag uint8

const (
	FlagZ  Flag = 0x1
	FlagNZ Flag = 0x2
	FlagL  Flag = 0x3
	FlagLE Flag = 0x4
	FlagG  Flag = 0x5
	FlagGE Flag = 0x6
)

// Jump instruction layouts.
//
// Layout of Jump is split into two parts.
// High 4 bits specify Flag (jump condition). Zero value of Flag
// means unconditional jump.
// Low 4 bits specify Data encoding.
const (
	// Jump to address stored in register.
	// Address must be in text segment.
	//
	// Data:
	//	RR - 1 byte - destination
	JumpReg Layout = 0x0

	// Jump to address stored in instruction immediate value.
	// Value is extended to text segment address.
	//
	// Data:
	//	II II II II - 4 bytes - destination
	JumpVal32 Layout = 0x1
)

// GetJumpSize returns size of jump instruction Data part.
func GetJumpDataSize(layout Layout) (uint64, error) {
	switch layout {
	case JumpReg:
		return 1, nil
	case JumpVal32:
		return 4, nil
	default:
		return 0, fmt.Errorf("unexpected layout (=0x%02X)", layout)
	}
}

func EncodeJumpLayout(flag Flag, lt Layout) Layout {
	return lt | Layout(flag<<4)
}

func DecodeJumpLayout(x uint8) (Flag, Layout) {
	flag := Flag(x >> 4)
	layout := Layout(x & 0xF)
	return flag, layout
}

var flagText = [...]string{
	FlagZ:  "z (== 0)",
	FlagNZ: "nz (!= 0)",
	FlagL:  "l (< x)",
	FlagLE: "le (<= x)",
	FlagG:  "g (> x)",
	FlagGE: "ge (>= x)",
}

func (f Flag) String() string {
	return flagText[f]
}
