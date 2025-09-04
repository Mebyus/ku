package asm

import (
	"github.com/mebyus/ku/goku/vm/ir"
	"github.com/mebyus/ku/goku/vm/opc"
)

func (a *Assembler) testVal(t ir.TestVal) {
	v := t.Val
	if v <= 0xF {
		a.testVal4(t.Dest, uint8(v))
		return
	}
	if v <= 0xFF {
		a.testVal8(t.Dest, uint8(v))
		return
	}
	if v <= 0xFFFF {
		a.testVal16(t.Dest, uint16(v))
	}
	if v <= 0xFFFFFFFF {
		a.testVal32(t.Dest, uint32(v))
		return
	}

	a.testVal64(t.Dest, v)
}

func (a *Assembler) testVal4(dest opc.Register, v uint8) {
	a.opcode(opc.Test)
	a.layout(opc.EncodeTestValLayout(opc.TestVal4, v))
	a.register(dest)
}

func (a *Assembler) testVal8(dest opc.Register, v uint8) {
	a.opcode(opc.Test)
	a.layout(opc.EncodeTestLayout(opc.TestVal8))
	a.register(dest)
	a.val8(v)
}

func (a *Assembler) testVal16(dest opc.Register, v uint16) {
	a.opcode(opc.Test)
	a.layout(opc.EncodeTestLayout(opc.TestVal16))
	a.register(dest)
	a.val16(v)
}

func (a *Assembler) testVal32(dest opc.Register, v uint32) {
	a.opcode(opc.Test)
	a.layout(opc.EncodeTestLayout(opc.TestVal32))
	a.register(dest)
	a.val32(v)
}

func (a *Assembler) testVal64(dest opc.Register, v uint64) {
	a.opcode(opc.Test)
	a.layout(opc.EncodeTestLayout(opc.TestVal64))
	a.register(dest)
	a.val64(v)
}
