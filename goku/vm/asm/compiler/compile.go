package compiler

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/vm/asm/ast"
	"github.com/mebyus/ku/goku/vm/ir"
)

func Compile(text *ast.Text) (*ir.Program, diag.Error) {
	if text.Entry.Name == "" {
		return nil, &diag.PinlessError{Text: "program has no entrypoint"}
	}

	c := Compiler{
		prog: ir.Program{
			Functions: make([]ir.Fun, 0, len(text.Functions)),
		},
		funs:   make(map[string]ir.FunName, len(text.Functions)),
		labels: make(map[string]ir.Label), // TODO: we can optimize this allocation after function indexing
	}

	err := c.indexFunctions(text.Functions)
	if err != nil {
		return nil, err
	}

	entry := text.Entry.Name
	i, ok := c.funs[entry]
	if !ok {
		return nil, &diag.SimpleMessageError{
			Pin:  text.Entry.Pin,
			Text: fmt.Sprintf("program has no \"%s\" function", entry),
		}
	}
	c.prog.EntryFun = i

	err = c.translateFunctions(text.Functions)
	if err != nil {
		return nil, err
	}

	return &c.prog, nil
}

type Compiler struct {
	prog ir.Program

	// Maps function string name to its integer name.
	funs map[string]ir.FunName

	// Maps label string name to its integer name.
	//
	// Only contains labels for currently translated function.
	labels map[string]ir.Label
}

func (c *Compiler) indexFunctions(funs []ast.Fun) diag.Error {
	for i, f := range funs {
		name := f.Name
		_, ok := c.funs[name]
		if ok {
			return &diag.SimpleMessageError{
				Pin:  f.Pin,
				Text: fmt.Sprintf("function \"%s\" was already declared in program", name),
			}
		}
		c.funs[name] = ir.FunName(i)
	}
	return nil
}

func (c *Compiler) translateFunctions(funs []ast.Fun) diag.Error {
	for i, f := range funs {
		err := c.translateFunction(ir.FunName(i), &f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) translateFunction(name ir.FunName, f *ast.Fun) diag.Error {
	if len(f.Atoms) == 0 {
		return &diag.SimpleMessageError{
			Pin:  f.Pin,
			Text: "empty function body",
		}
	}

	clear(c.labels)
	count := c.prog.LabelsCount
	for _, label := range f.Labels {
		_, ok := c.labels[label]
		if ok {
			return &diag.SimpleMessageError{
				Pin:  f.Pin,
				Text: fmt.Sprintf("label \"%s\" was already placed in function \"%s\"", label, f.Name),
			}
		}
		c.labels[label] = ir.Label(count)
		count += 1
	}
	c.prog.LabelsCount = count

	fun := ir.Fun{Name: name}
	for _, atom := range f.Atoms {
		a, err := c.translateAtom(atom)
		if err != nil {
			return err
		}
		fun.Atoms = append(fun.Atoms, a)
	}
	c.prog.Functions = append(c.prog.Functions, fun)
	return nil
}

func (c *Compiler) translateAtom(atom ast.Atom) (ir.Atom, diag.Error) {
	switch a := atom.(type) {
	case ast.Place:
		return ir.Place{Label: c.labels[a.Name]}, nil
	case ast.Instruction:
		switch a.Mnemonic {
		case "halt":
			// TODO: check operands
			return ir.Halt{}, nil
		case "ret":
			return ir.Ret{}, nil
		case "nop":
			return ir.Nop{}, nil
		case "inc":
			return c.translateInc(a)
		case "set":
			return c.translateSet(a)
		default:
			return nil, &diag.SimpleMessageError{
				Pin:  a.Pin,
				Text: fmt.Sprintf("unknown instruction mnemonic \"%s\"", a.Mnemonic),
			}
		}
	default:
		panic(fmt.Sprintf("unexpected atom type (%T)", a))
	}
}
