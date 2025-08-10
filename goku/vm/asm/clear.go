package asm

import (
	"github.com/mebyus/ku/goku/vm/ir"
	"github.com/mebyus/ku/goku/vm/opc"
)

func (a *Assembler) clearReg(t ir.ClearReg) {
	a.opcode(opc.Clear)
	a.layout(opc.ClearReg)
	a.register(t.Reg)
}
