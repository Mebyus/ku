package ir

import "github.com/mebyus/ku/goku/vm/opc"

type IncReg struct {
	nodeAtom

	Dest   opc.Register
	Source opc.Register
}

type IncVal struct {
	nodeAtom

	Val  uint64
	Dest opc.Register
}

type DecReg struct {
	nodeAtom

	Dest   opc.Register
	Source opc.Register
}

type DecVal struct {
	nodeAtom

	Val  uint64
	Dest opc.Register
}
