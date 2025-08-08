package asm

import (
	"github.com/mebyus/ku/goku/vm/ir"
	"github.com/mebyus/ku/goku/vm/kvx"
)

func Assemble(prog *ir.Program) *kvx.Program {
	a := Assembler{
		tab: OffsetsTable{
			Data:   make([]uint32, len(prog.Data)),
			Labels: make([][]uint32, len(prog.Functions)),
		},
	}

	a.encodeDataSegment(prog.Data)
	a.encodeTextSegment(prog.Functions)

	return &a.prog
}

type OffsetsTable struct {
	// Translates data entry integer name to its offset
	// in data segment.
	Data []uint32

	// Translates function integer name + label integer name
	// to label offset in text segment.
	Labels [][]uint32
}

type Assembler struct {
	prog kvx.Program

	tab OffsetsTable

	// offset into current segment
	offset uint32
}

func (a *Assembler) encodeDataSegment(data []ir.DataEntry) {
	for i, d := range data {
		offset := uint32(len(a.prog.Data))
		a.prog.Data = append(a.prog.Data, d.Val...)
		a.tab.Data[i] = offset
	}
}

func (a *Assembler) encodeTextSegment(functions []ir.Fun) {
	for i, f := range functions {
		_ = i
		_ = f
	}
}
