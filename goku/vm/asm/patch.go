package asm

import (
	"encoding/binary"

	"github.com/mebyus/ku/goku/vm/ir"
)

// CallPatchEntry contains information on how to patch
// call address in program text.
type CallPatchEntry struct {
	// Patch will apply address (text offset) of this function.
	Fun ir.FunName

	// Where patch will be placed.
	Offset uint32
}

// JumpPatchEntry contains information on how to patch
// jump address in program text.
type JumpPatchEntry struct {
	// Patch will apply address (text offset) of this label.
	Label ir.Label

	// Where patch will be placed.
	Offset uint32
}

// Fill call address placeholders in program text.
func (a *Assembler) patchCalls() {
	for _, p := range a.patch.Calls {
		address := a.tab.Functions[p.Fun]
		a.patch32(p.Offset, address)
	}
}

// Fill jump address placeholders in program text.
func (a *Assembler) patchJumps() {
	for _, p := range a.patch.Jumps {
		address := a.tab.Labels[p.Label]
		a.patch32(p.Offset, address)
	}
}

// Patch text segment with value {v} at specified offset.
//
// Value will be encoded in little endian.
func (a *Assembler) patch32(offset, v uint32) {
	b := a.prog.Text[offset : offset+4]
	binary.LittleEndian.PutUint32(b, v)
}
