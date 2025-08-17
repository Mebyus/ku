package opc

import (
	"fmt"
	"strconv"
)

// Register represents register name in used in
// instruction encoding.
//
// Values 0 <= r < 64 represent general purpose registers.
//
// Values with highest bit set to 1 represent special VM
// control registers.
type Register uint8

const (
	// Instruction pointer.
	RegIP Register = 0x80 + iota

	// Stack poiner.
	//
	// Read-only. Managed by VM.
	RegSP

	// Frame pointer.
	//
	// Read-only. Managed by VM.
	RegFP

	// Syscall register.
	//
	// Read-write.
	RegSC

	// Comparison flags register.
	//
	// Read-only. Managed by VM.
	RegCF

	// Register which tracks number of executed instructions.
	//
	// Read-only. Managed by VM.
	RegClock
)

func (r Register) Special() bool {
	return r&0x80 != 0
}

func (r Register) String() string {
	if !r.Special() {
		return "r" + strconv.FormatUint(uint64(r), 10)
	}

	switch r {
	case RegIP:
		return "ip"
	case RegSP:
		return "sp"
	case RegFP:
		return "fp"
	case RegSC:
		return "sc"
	case RegCF:
		return "cf"
	case RegClock:
		return "clock"
	default:
		panic(fmt.Sprintf("unexpected special register (=%d)", r))
	}
}
