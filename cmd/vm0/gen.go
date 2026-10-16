package main

import "fmt"

// inst represents code instruction
type inst struct {
	typ int

	// not all instructions use this
	//
	//	push  - value pushed on stack
	//	load  - variable name
	//	store - variable name
	//	call  - function name
	//	jump  - instruction address
	//	jumpz - instruction address
	arg int32
}

// vm opcodes
const (
	// should be zero value
	opTRAP = iota
	opHALT
	opPRINT
	opCALL
	opRET
	opJUMP  // unconditional jump to target instruction address
	opJZ    // pop EVS and jump to target instruction address if value is 0
	opJNZ   // pop EVS and jump to target instruction address if value is not 0
	opTAKE  // pop value from EVS and place it into stack memory as call argument
	opDRAIN // pop value from EVS and discard it
	opPUSH  // push value stored in local variable onto EVS
	opPOP   // pop value from EVS and store it into local variable
	opPUT   // push constant value onto EVS
	opADD
	opSUB
	opMUL
	opEQUAL
	opLESS

	// opLOAD     // push global variable value onto EVS
	// opSTORE    // pop value from EVS and save it to global variable
)

type Program struct {
	text []inst

	rom []byte

	// must containt instruction adderreses (=index into text)
	// in ascending order
	funs []finfo
}

type compiler struct {
	prog Program

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

func (c *compiler) compile(funs []*Function) {
	c.nfun = uint32(len(funs))
	c.prog.funs = make([]finfo, len(funs))
	for i := range len(funs) {
		c.prog.funs[i] = finfo{
			narg: funs[i].narg,
			nret: funs[i].nret,
		}
	}

	for i, f := range funs {
		c.prog.funs[i].addr = uint32(len(c.prog.text))
		c.fun(f)
		c.block(f.body)
		c.prog.funs[i].nloc = c.lcap - f.narg
		// TODO: ? insert return if function
		// has no return values and has no return
		// at the end of body, otherwise should
		// panic
	}
}

func (c *compiler) stm(s stm) {
	switch s := s.(type) {
	case *Print:
		c.exp(s.exp)
		c.add(inst{typ: opPRINT})
	case *Assign:
		c.exp(s.exp)
		c.add(inst{typ: opPOP, arg: int32(s.name)})
		c.def(s.name)
	case *Return:
		c.exp(s.exp)
		c.add(inst{typ: opRET})
	case *Call:
		c.call(s)
		nret := c.prog.funs[s.name].nret
		for range nret {
			// since call was used as statement
			// that means we need to discard
			// returned values from EVS
			c.add(inst{typ: opDRAIN})
		}
	case *Ifthen:
		c.exp(s.exp)
		i := c.add(inst{typ: opJZ})
		c.block(s.body)

		// backpatch jump address
		c.prog.text[i].arg = int32(len(c.prog.text))
	case *While:
		i := c.add(inst{typ: opJUMP}) // jump to loop epilogue
		c.block(s.body)
		j := len(c.prog.text)
		c.exp(s.exp)
		c.add(inst{typ: opJNZ, arg: int32(i + 1)}) // jump to loop body start
		c.prog.text[i].arg = int32(j)              // backpatch init jump address to epilogue
	default:
		panic(s)
	}
}

func (c *compiler) block(body []stm) {
	for _, s := range body {
		c.stm(s)
	}
}

func (c *compiler) call(call *Call) {
	if call.name >= c.nfun {
		panic(fmt.Sprintf("call to undefined function (%d)", call.name))
	}

	narg := c.prog.funs[call.name].narg
	if uint32(len(call.args)) != narg {
		panic("number of arguments mismatch")
	}

	for _, arg := range call.args {
		c.exp(arg)
		c.add(inst{typ: opTAKE})
	}
	c.add(inst{typ: opCALL, arg: int32(call.name)})
}

func (c *compiler) exp(exp Exp) {
	switch e := exp.(type) {
	case *Integer:
		c.add(inst{typ: opPUT, arg: e.Val})
	case *String:
		p, n := c.str(e.Val)
		c.add(inst{typ: opPUT, arg: int32(n)})
		c.add(inst{typ: opPUT, arg: int32(p)})
	case *Local:
		if !c.has(e.Name) {
			panic(fmt.Sprintf("undefined local variable [%d] used in expression", e.Name))
		}
		c.add(inst{typ: opPUSH, arg: int32(e.Name)})
	case *Call:
		c.call(e)
	case *BinExp:
		c.exp(e.A)
		c.exp(e.B)

		var typ int
		switch e.Op {
		case bopAdd:
			typ = opADD
		case bopSub:
			typ = opSUB
		case bopMul:
			typ = opMUL
		case bopEql:
			typ = opEQUAL
		case bopLess:
			typ = opLESS
		default:
			panic(e.Op)
		}
		c.add(inst{typ: typ})
	default:
		panic(e)
	}
}

// start new function
func (c *compiler) fun(f *Function) {
	clear(c.local)
	for i := range f.narg {
		c.local[i] = struct{}{}
	}
	c.lcap = f.narg
}

// get instruction address of n instructions ahead
func (c *compiler) ahead(n int) uint32 {
	return uint32(len(c.prog.text) + n)
}

// returns address of added instruction
func (c *compiler) add(x inst) uint32 {
	i := uint32(len(c.prog.text))
	c.prog.text = append(c.prog.text, x)
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

// intern a string and return its (p, n)
//
//	p - pointer into rom
//	n - its length in bytes
func (c *compiler) str(s string) (uint32, uint32) {
	p, ok := c.intern[s]
	if ok {
		return p, uint32(len(s))
	}

	// separate bytes of each distinct string with
	// single 0 byte for debugging purposes and
	// compatibility with C
	//
	// this also reserves nil pointer address at
	// the rom start
	c.prog.rom = append(c.prog.rom, 0)

	p = uint32(len(c.prog.rom))
	c.prog.rom = append(c.prog.rom, s...)

	c.intern[s] = p
	return p, uint32(len(s))
}

func printProgram(prog *Program) {
	if len(prog.text) == 0 {
		return
	}

	var fun int // current function index
	for i, x := range prog.text {
		if fun < len(prog.funs) && i == int(prog.funs[fun].addr) {
			if fun != 0 {
				fmt.Println()
			}
			fmt.Printf("FUN %d (0x%04X):\n", fun, i)
			fun += 1
		} else if i%4 == 0 {
			fmt.Printf("0x%04X:\n", i)
		}

		switch x.typ {
		case opTRAP:
			fmt.Printf("    TRAP\n")
		case opHALT:
			fmt.Printf("    HALT\n")
		case opPRINT:
			fmt.Printf("    PRINT\n")
		case opCALL:
			fmt.Printf("    CALL   #.%d\n", x.arg)
		case opRET:
			fmt.Printf("    RET\n")
		case opJUMP:
			fmt.Printf("    JUMP   0x%X  (%d)\n", x.arg, x.arg-int32(i))
		case opJZ:
			fmt.Printf("    JZ     0x%X  (%d)\n", x.arg, x.arg-int32(i))
		case opJNZ:
			fmt.Printf("    JNZ    0x%X  (%d)\n", x.arg, x.arg-int32(i))
		case opTAKE:
			fmt.Printf("    TAKE\n")
		case opDRAIN:
			fmt.Printf("    DRAIN\n")
		case opPUSH:
			fmt.Printf("    PUSH   [%d]\n", x.arg)
		case opPOP:
			fmt.Printf("    POP    [%d]\n", x.arg)
		case opPUT:
			fmt.Printf("    PUT    0x%X  (%d)\n", uint32(x.arg), x.arg)
		case opADD:
			fmt.Printf("    ADD\n")
		case opSUB:
			fmt.Printf("    SUB\n")
		case opMUL:
			fmt.Printf("    MUL\n")
		case opEQUAL:
			fmt.Printf("    EQUAL\n")
		case opLESS:
			fmt.Printf("    LESS\n")
		default:
			panic(x.typ)
		}
	}
}
