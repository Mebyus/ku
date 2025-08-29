package asm

import (
	"github.com/mebyus/ku/goku/vm/ir"
	"github.com/mebyus/ku/goku/vm/opc"
)

// encode trap instruction.
func (a *Assembler) trap(t ir.Trap) {
	a.opcode(opc.Sys)
	a.layout(opc.Layout(opc.Trap))
}

// encode halt instruction.
func (a *Assembler) halt(t ir.Halt) {
	a.opcode(opc.Sys)
	a.layout(opc.Layout(opc.Halt))
}

// encode nop instruction.
func (a *Assembler) nop(t ir.Nop) {
	a.opcode(opc.Sys)
	a.layout(opc.Layout(opc.Nop))
}

// encode syscall instruction.
func (a *Assembler) syscall(t ir.SysCall) {
	a.opcode(opc.Sys)
	a.layout(opc.Layout(opc.SysCall))
}

func (a *Assembler) ret(t ir.Ret) {
	a.opcode(opc.Sys)
	a.layout(opc.Layout(opc.Ret))
}
