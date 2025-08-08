package ir

// Halt represents halt instruction.
type Halt struct {
	nodeAtom
}

// Trap represents explicit trap instruction.
type Trap struct {
	nodeAtom
}

// Nop represents nop instruction.
type Nop struct {
	nodeAtom
}

// SysCall represents syscall instruction.
type SysCall struct {
	nodeAtom
}

// Ret represents ret instruction.
type Ret struct {
	nodeAtom
}
