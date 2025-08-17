package asm

import (
	"github.com/mebyus/ku/goku/vm/ir"
	"github.com/mebyus/ku/goku/vm/opc"
)

func (a *Assembler) setReg(t ir.SetReg) {
	a.opcode(opc.Set)
	a.layout(opc.EncodeSetLayout(opc.SetReg))
	a.register(t.Dest)
	a.register(t.Source)
}

func (a *Assembler) setVal(t ir.SetVal) {
	v := t.Val
	if v <= 0xF {
		a.setVal4(t.Dest, uint8(v))
		return
	}
	if v <= 0xFF {
		a.setVal8(t.Dest, uint8(v))
		return
	}
	if v <= 0xFFFF {
		a.setVal16(t.Dest, uint16(v))
	}
	if v <= 0xFFFFFFFF {
		a.setVal32(t.Dest, uint32(v))
		return
	}

	a.setVal64(t.Dest, v)
}

func (a *Assembler) setVal4(dest opc.Register, v uint8) {
	a.opcode(opc.Set)
	a.layout(opc.EncodeSetValLayout(opc.SetVal4, v))
	a.register(dest)
}

func (a *Assembler) setVal8(dest opc.Register, v uint8) {
	a.opcode(opc.Set)
	a.layout(opc.EncodeSetLayout(opc.SetVal8))
	a.register(dest)
	a.val8(v)
}

func (a *Assembler) setVal16(dest opc.Register, v uint16) {
	a.opcode(opc.Set)
	a.layout(opc.EncodeSetLayout(opc.SetVal16))
	a.register(dest)
	a.val16(v)
}

func (a *Assembler) setVal32(dest opc.Register, v uint32) {
	a.opcode(opc.Set)
	a.layout(opc.EncodeSetLayout(opc.SetVal32))
	a.register(dest)
	a.val32(v)
}

func (a *Assembler) setVal64(dest opc.Register, v uint64) {
	a.opcode(opc.Set)
	a.layout(opc.EncodeSetLayout(opc.SetVal64))
	a.register(dest)
	a.val64(v)
}
