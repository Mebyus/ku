package vm32b

import (
	"fmt"

	"github.com/mebyus/ku/internal/ku/enums/bop"
	"github.com/mebyus/ku/internal/vm32b/ir"
)

func Compile(funs []ir.Fun) *Prog {
	var c compiler
	c.init()
	c.compile(funs)
	return &c.prog
}

type compiler struct {
	prog Prog

	// maintains current set of defined local variables
	local map[ /* local variable name */ uint32]struct{}

	// map for string interning
	//
	// maps interned string into its pointer in rom
	intern map[string]uint32

	// total number of functions in program
	nfun uint32

	// number of allocated local variables
	// including function arguments and
	// stack memory gaps
	lcap uint32
}

func (c *compiler) init() {
	c.local = make(map[uint32]struct{})
	c.intern = make(map[string]uint32)
}

func (c *compiler) compile(funs []ir.Fun) {
	c.nfun = uint32(len(funs))
	c.prog.funs = make([]finfo, len(funs))
	for i := range len(funs) {
		c.prog.funs[i] = finfo{
			narg: funs[i].Narg,
			nret: funs[i].Nret,
		}
	}

	for i, f := range funs {
		c.prog.funs[i].addr = ia(len(c.prog.text))
		c.startfun(&f)
		c.block(f.Body)
		c.prog.funs[i].nloc = c.lcap - f.Narg
		// TODO: ? insert return if function
		// has no return values and has no return
		// at the end of body, otherwise should
		// panic
	}
}

func (c *compiler) block(body []ir.Stm) {
	for _, s := range body {
		c.stm(s)
	}
}

func (c *compiler) stm(s ir.Stm) {
	switch s := s.(type) {
	// case *Print:
	// 	c.exp(s.exp)
	// 	c.add(inst{typ: opPRINT})
	case *ir.Assign:
		c.exp(s.Exp)
		c.emit(opPOP, s.Name)
		c.def(s.Name)
	case *ir.Return:
		if s.Exp != nil {
			panic("stub")
			// c.exp(s.exp)
		}
		c.emit(opRET, 0)
	// case *Call:
	// 	c.call(s)
	// 	nret := c.prog.funs[s.name].nret
	// 	for range nret {
	// 		// since call was used as statement
	// 		// that means we need to discard
	// 		// returned values from EVS
	// 		c.add(inst{typ: opDRAIN})
	// 	}
	// case *Ifthen:
	// 	c.exp(s.exp)
	// 	i := c.add(inst{typ: opJZ})
	// 	c.block(s.body)

	// 	// backpatch jump address
	// 	c.prog.text[i].arg = int32(len(c.prog.text))
	case *ir.While:
		i := c.emit(opJUMP, 0) // jump to loop epilogue
		c.block(s.Body)
		j := len(c.prog.text)
		c.exp(s.Exp)
		c.emit(opJNZ, uint32(i+1))     // jump to loop body start
		c.prog.text[i].arg = uint32(j) // backpatch init jump address to epilogue
	default:
		panic(fmt.Sprintf("unexpected %T statement", s))
	}
}

func (c *compiler) exp(exp ir.Exp) {
	switch e := exp.(type) {
	case *ir.Integer:
		c.emit(opPUT, e.Val)
	// case *String:
	// 	p, n := c.str(e.Val)
	// 	c.add(inst{typ: opPUT, arg: int32(n)})
	// 	c.add(inst{typ: opPUT, arg: int32(p)})
	case *ir.Local:
		if !c.has(e.Name) {
			panic(fmt.Sprintf("undefined local variable [%d] used in expression", e.Name))
		}
		c.emit(opPUSH, e.Name)
	// case *Call:
	// 	c.call(e)
	case *ir.BinExp:
		c.exp(e.A)
		c.exp(e.B)

		var op op
		switch e.Op {
		case bop.Add:
			op = opADD
		case bop.Sub:
			op = opSUB
		case bop.Mul:
			op = opMUL
		case bop.Equal:
			op = opEQUAL
		case bop.Less:
			op = opLESS
		default:
			panic(e.Op)
		}
		c.emit(op, 0)
	default:
		panic(fmt.Sprintf("unexpected %T expression", e))
	}
}

// start new function
func (c *compiler) startfun(f *ir.Fun) {
	// clear(c.local)
	// for i := range f.narg {
	// 	c.local[i] = struct{}{}
	// }
	// c.lcap = f.narg
}

// add new instruction to compiled program output
//
// instruction argument is optional for some instructions
//
// returns address of emitted instruction
func (c *compiler) emit(op op, arg uint32) ia {
	i := ia(len(c.prog.text))
	c.prog.text = append(c.prog.text, inst{
		op:  op,
		arg: arg,
	})
	return i
}

// checks that local variable with the given name is defined
func (c *compiler) has(name uint32) bool {
	_, ok := c.local[name]
	return ok
}

// mark local variable as defined
func (c *compiler) def(name uint32) {
	// number of required allocated local variables
	// for given name to be valid
	nreq := name + 1

	if nreq > c.lcap {
		c.lcap = nreq
	}
	c.local[name] = struct{}{}
}
