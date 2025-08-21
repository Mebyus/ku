package compiler

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/ir"
	"github.com/mebyus/ku/goku/vm/opc"
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

func (c *Compiler) translateSet(s ast.Instruction) (ir.Atom, diag.Error) {
	switch len(s.Operands) {
	case 2:
		dest := s.Operands[0]
		d, ok := dest.(ast.Register)
		if !ok {
			return nil, &diag.SimpleMessageError{
				Pin:  s.Pin,
				Text: fmt.Sprintf("bad operand: destination must be register, got (%T)", dest),
			}
		}
		if d.Name.Special() && d.Name != opc.RegSC {
			return nil, &diag.SimpleMessageError{
				Pin:  s.Pin,
				Text: fmt.Sprintf("bad operand: set instruction cannot have \"%s\" register as destination", d.Name),
			}
		}

		source := s.Operands[1]
		switch c := source.(type) {
		case ast.Register:
			return ir.SetReg{
				Dest:   d.Name,
				Source: c.Name,
			}, nil
		case ast.Integer:
			return ir.SetVal{
				Dest: d.Name,
				Val:  c.Val,
			}, nil
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  s.Pin,
				Text: fmt.Sprintf("bad operand: source must be register or immediate, got (%T)", c),
			}
		}
	default:
		return nil, wrongOperandsNumber(s.Pin, len(s.Operands), "2")
	}
}

func (c *Compiler) translateJump(s ast.Instruction) (ir.Atom, diag.Error) {
	if s.Variant != "" {
		panic("stub")
	}

	switch len(s.Operands) {
	case 1:
		d := s.Operands[0]
		label, ok := d.(ast.Label)
		if !ok {
			return nil, &diag.SimpleMessageError{
				Pin:  s.Pin,
				Text: fmt.Sprintf("bad operand: destination must be label, got (%T)", d),
			}
		}
		l, ok := c.labels[label.Name]
		if !ok {
			return nil, &diag.SimpleMessageError{
				Pin:  s.Pin,
				Text: fmt.Sprintf("label \"%s\" not placed in this function", label.Name),
			}
		}
		return ir.JumpLabel{Label: l}, nil
	default:
		return nil, wrongOperandsNumber(s.Pin, len(s.Operands), "1")
	}
}

func (c *Compiler) translateCall(s ast.Instruction) (ir.Atom, diag.Error) {
	if s.Variant != "" {
		// TODO: return error
	}

	switch len(s.Operands) {
	case 1:
		d := s.Operands[0]
		symbol, ok := d.(ast.Symbol)
		if !ok {
			return nil, &diag.SimpleMessageError{
				Pin:  s.Pin,
				Text: fmt.Sprintf("bad operand: destination must be symbol, got (%T)", d),
			}
		}
		fun, ok := c.funs[symbol.Name]
		if !ok {
			return nil, &diag.SimpleMessageError{
				Pin:  s.Pin,
				Text: fmt.Sprintf("function \"%s\" not declared", symbol.Name),
			}
		}
		return ir.CallFun{Fun: fun}, nil
	default:
		return nil, wrongOperandsNumber(s.Pin, len(s.Operands), "1")
	}
}
