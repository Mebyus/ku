package ir

import "github.com/mebyus/ku/goku/vm/opc"

type SetReg struct {
	nodeAtom

	Dest   opc.Register
	Source opc.Register
}

type SetVal struct {
	nodeAtom

	Val  uint64
	Dest opc.Register
}
