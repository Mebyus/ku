package compiler

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/ir"
)

func wrongOperandsNumber(pin srcmap.Pin, got int, want string) diag.Error {
	return &diag.SimpleMessageError{
		Pin:  pin,
		Text: fmt.Sprintf("wrong number of operands: got %d, want %s", got, want),
	}
}

func (c *Compiler) translateInc(s ast.Instruction) (ir.Atom, diag.Error) {
	switch len(s.Operands) {
	case 1:
		op := s.Operands[0]
		reg, ok := op.(ast.Register)
		if !ok {
			return nil, &diag.SimpleMessageError{
				Pin:  s.Pin,
				Text: fmt.Sprintf("bad operand: want register, got (%T)", op),
			}
		}
		if reg.Name.Special() {
			return nil, &diag.SimpleMessageError{
				Pin:  s.Pin,
				Text: "inc instruction cannot have special register as destination operand",
			}
		}
		return ir.IncVal{
			Dest: reg.Name,
			Val:  1,
		}, nil
	case 2:
	default:
		return nil, wrongOperandsNumber(s.Pin, len(s.Operands), "1-2")
	}

	panic("stub")
}
