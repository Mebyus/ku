package opc

// JumpFlag specifies jump condition.
type JumpFlag uint8

const (
	FlagZ  JumpFlag = 0x1
	FlagNZ JumpFlag = 0x2
	FlagL  JumpFlag = 0x3
	FlagLE JumpFlag = 0x4
	FlagG  JumpFlag = 0x5
	FlagGE JumpFlag = 0x6
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

func EncodeJumpLayout(flag JumpFlag, lt Layout) Layout {
	return lt | Layout(flag<<4)
}

func DecodeJumpLayout(x uint8) (JumpFlag, Layout) {
	flag := JumpFlag(x >> 4)
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

func (f JumpFlag) String() string {
	return flagText[f]
}
