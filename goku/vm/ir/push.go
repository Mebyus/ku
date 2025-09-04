package ir

import "github.com/mebyus/ku/goku/vm/opc"

type PushReg struct {
	nodeAtom

	Reg opc.Register
}

type PopReg struct {
	nodeAtom

	Reg opc.Register
}
