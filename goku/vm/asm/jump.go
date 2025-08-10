package asm

import (
	"github.com/mebyus/ku/goku/vm/ir"
	"github.com/mebyus/ku/goku/vm/opc"
)

// encode jump instruction with static label address.
func (a *Assembler) jumpLabel(t ir.JumpLabel) {
	address := a.tab.Labels[t.Label]

	a.opcode(opc.Jump)
	a.layout(opc.EncodeJumpLayout(t.Flag, opc.JumpVal32))

	if address == 0 {
		// this check reduces number of patches we need to apply later,
		// other address values are obviously already filled with correct
		// values
		//
		// since only one label in all program can have address 0
		// number of "false-positive" patches will be relatively low
		a.patch.Jumps = append(a.patch.Jumps, JumpPatchEntry{
			Label:  t.Label,
			Offset: a.textOffset(),
		})
	}

	a.val32(address)
}
