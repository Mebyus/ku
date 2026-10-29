package kub

import "fmt"

// Arch represents host or target machine architecture.
type Arch uint32

const (
	ArchVM32B = iota
	ArchX86
	ArchAMD64
	ArchARM64
)

// Size returns size (in bytes) of pointer for this architecture.
func (a Arch) Size() uint32 {
	switch a {
	case ArchVM32B, ArchX86:
		return 4
	case ArchAMD64, ArchARM64:
		return 8
	default:
		panic(fmt.Sprintf("unexpected arch (=%d)", a))
	}
}

// OS represents host or target machine operating system.
type OS uint32

const (
	// No operating system.
	Bare OS = iota

	Linux
	Windows
	Darwin
)
