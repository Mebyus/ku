package compiler

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/ir"
)

func (c *Compiler) translatePush(s ast.Instruction) (ir.Atom, diag.Error) {
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
		return ir.PushReg{Reg: reg.Name}, nil
	default:
		return nil, wrongOperandsNumber(s.Pin, len(s.Operands), "1")
	}
}

func (c *Compiler) translatePop(s ast.Instruction) (ir.Atom, diag.Error) {
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
		return ir.PopReg{Reg: reg.Name}, nil
	default:
		return nil, wrongOperandsNumber(s.Pin, len(s.Operands), "1")
	}
}
