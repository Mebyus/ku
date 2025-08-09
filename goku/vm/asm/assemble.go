package asm

import (
	"encoding/binary"
	"fmt"

	"github.com/mebyus/ku/goku/vm/ir"
	"github.com/mebyus/ku/goku/vm/kvx"
)

func Assemble(prog *ir.Program) *kvx.Program {
	a := Assembler{
		tab: OffsetsTable{
			Data:      make([]uint32, len(prog.Data)),
			Labels:    make([]uint32, prog.LabelsCount),
			Functions: make([]uint32, len(prog.Functions)),
		},
	}

	a.encodeDataSegment(prog.Data)
	a.encodeTextSegment(prog.Functions)

	a.patchCalls()
	a.patchJumps()

	return &a.prog
}

type OffsetsTable struct {
	// Translates data entry integer name to its offset
	// in data segment.
	Data []uint32

	// Translates label integer name to label offset in text segment.
	Labels []uint32

	// Translates function integer name to function offset in text segment.
	Functions []uint32
}

type PatchTable struct {
	Calls []CallPatchEntry

	Jumps []JumpPatchEntry
}

type Assembler struct {
	prog kvx.Program

	tab OffsetsTable

	patch PatchTable
}

func (a *Assembler) encodeDataSegment(data []ir.DataEntry) {
	for i, d := range data {
		offset := uint32(len(a.prog.Data))
		a.prog.Data = append(a.prog.Data, d.Val...)
		a.tab.Data[i] = offset
	}
}

func (a *Assembler) encodeTextSegment(functions []ir.Fun) {
	for _, f := range functions {
		a.encodeFun(f)
	}
}

// Returns encoder offset into text segment.
func (a *Assembler) textOffset() uint32 {
	return uint32(len(a.prog.Text))
}

func (a *Assembler) encodeFun(f ir.Fun) {
	a.alignFun()
	a.tab.Functions[f.Name] = a.textOffset()

	for _, atom := range f.Atoms {
		a.encodeAtom(atom)
	}
}

func (a *Assembler) encodeAtom(atom ir.Atom) {
	switch t := atom.(type) {
	case ir.Halt:
	case ir.Place:
		a.tab.Labels[t.Label] = a.textOffset()
	case ir.CallFun:
		a.callFun(t)
	case ir.JumpLabel:
	default:
		panic(fmt.Sprintf("unexpected atom (%T)", t))
	}
}

// Encode 32-bit integer into text segment and advance encoder offset.
func (a *Assembler) val32(v uint32) {
	a.prog.Text = binary.LittleEndian.AppendUint32(a.prog.Text, v)
}

// Encode 8-bit integer into text segment and advance encoder offset.
func (a *Assembler) val8(v uint8) {
	a.prog.Text = append(a.prog.Text, v)
}

// aligns encoder offset to start function encoding.
func (a *Assembler) alignFun() {

}
