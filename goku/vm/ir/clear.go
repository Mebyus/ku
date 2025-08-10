package ir

import "github.com/mebyus/ku/goku/vm/opc"

type ClearReg struct {
	nodeAtom

	Reg opc.Register
}
