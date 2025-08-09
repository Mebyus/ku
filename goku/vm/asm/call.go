package asm

import "github.com/mebyus/ku/goku/vm/ir"

// encode call instruction with static function address.
func (a *Assembler) callFun(t ir.CallFun) {
	address := a.tab.Functions[t.Fun]

	// TODO: encode call opcode and layout

	if address == 0 {
		// this check reduces number of patches we need to apply later,
		// other address values are obviously already filled with correct
		// values
		//
		// since only one function in all program can have address 0
		// number of "false-positive" patches will be relatively low
		a.patch.Calls = append(a.patch.Calls, CallPatchEntry{
			Fun:    t.Fun,
			Offset: a.textOffset(),
		})
	}

	a.val32(address)
}
