package asm

import (
	"github.com/mebyus/ku/goku/vm/ir"
	"github.com/mebyus/ku/goku/vm/opc"
)

// encode trap instruction.
func (a *Assembler) trap(t ir.Trap) {
	a.opcode(opc.Trap)
	a.layout(0)
}

// encode halt instruction.
func (a *Assembler) halt(t ir.Halt) {
	a.opcode(opc.Halt)
	a.layout(0)
}

// encode nop instruction.
func (a *Assembler) nop(t ir.Nop) {
	a.opcode(opc.Nop)
	a.layout(0)
}

// encode syscall instruction.
func (a *Assembler) syscall(t ir.SysCall) {
	a.opcode(opc.SysCall)
	a.layout(0)
}

func (a *Assembler) ret(t ir.Ret) {
	a.opcode(opc.Ret)
	a.layout(0)
}
