package ir

import "github.com/mebyus/ku/goku/vm/opc"

// JumpLabel represents instruction jump to label.
type JumpLabel struct {
	nodeAtom

	// Label name of jump destination.
	Label Label

	// Equals 0 for unconditional jump.
	Flag opc.JumpFlag
}

// JumpReg represents instruction jump to address stored in register.
type JumpReg struct {
	nodeAtom

	Reg opc.Register

	// Equals 0 for unconditional jump.
	Flag opc.JumpFlag
}
