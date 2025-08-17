package asm

import (
	"github.com/mebyus/ku/goku/vm/ir"
	"github.com/mebyus/ku/goku/vm/opc"
)

func (a *Assembler) incReg(t ir.IncReg) {
	a.opcode(opc.Inc)
	a.layout(opc.IncReg)
	a.register(t.Dest)
	a.register(t.Source)
}

func (a *Assembler) incVal(t ir.IncVal) {
	v := t.Val
	if v < 16 {
		a.incVal4(t.Dest, uint8(v))
		return
	}
	if v <= 0xFFFFFFFF {
		a.incVal32(t.Dest, uint32(v))
		return
	}

	a.incVal64(t.Dest, v)
}

func (a *Assembler) incVal4(dest opc.Register, v uint8) {
	a.opcode(opc.Inc)
	a.layout(opc.EncodeIncTinyLayout(v))
	a.register(dest)
}

func (a *Assembler) incVal32(dest opc.Register, v uint32) {
	a.opcode(opc.Inc)
	a.layout(opc.IncVal32)
	a.register(dest)
	a.val32(v)
}

func (a *Assembler) incVal64(dest opc.Register, v uint64) {
	a.opcode(opc.Inc)
	a.layout(opc.IncVal64)
	a.register(dest)
	a.val64(v)
}
