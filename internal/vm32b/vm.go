package vm32b

import "fmt"

// Implements 32-bit stack-based bytecode virtual machine.
//
// Evaluation stack (EVS for short) is used as "infinite"
// set of general purpose registers. It can hold
// temporary and intermediate values for evaluating
// expressions, passing return values, preparing
// call arguments, etc. EVS is not addressable and
// offers only push/pop like operations.
//
// Call frames are stored in a separate stack
// inaccessible to executed code, so there is
// no way to corrupt it, besided call stack overflow.
//
// Each call frame gets its own local variables via
// access to shared stack memory which is automatically
// managed by VM. Call arguments and local variables
// share "namespace" and are addressed by their local "names".
// Local variable "name" is basically index into call frame
// stack memory. For example if function was called
// with 2 arguments, then:
//
//	0 - first argument
//	1 - second argument
//	2 - first local variable
//	3 - second local variable
//	...
//
// Essentially arguments and local variables are all just
// low-level stack memory address slots for a function.
//
// Calling convention is quite simple. To make a call
// caller must place required (and it should exactly match
// number of function parameters) number of arguments
// onto stack memory in direct order, meaning that
// first argument gets lower address. VM offers special
// instruction TAKE for this operation of arranging call
// argument, it pops value from EVS and places it into
// correct stack memory slot. Later (after executing CALL
// instruction) that slot just becomes addressable local
// variable in the new call frame.
//
// When function call ends it must push returned values
// (if any) onto EVS. Returned values are pushed in
// direct order, meaning that first returned value
// comes onto EVS first. Thus when caller wants to
// pop and use results from EVS they will come in
// reverse order.
//
// Beyond arranging call arguments via TAKE instruction
// both caller and callee are not required to do any magic
// function prologue or epilogue, just CALL + RET pair.
// However there is a special case when function is
// called, but its returned values are not needed by the
// caller. In that case caller must explicitly pop
// returned values and discard them via DRAIN.
//
// Pointers can work with 4 distinct memory regions:
//
//	0 - read-only memory (comes with the program) (=0b00)
//	1 - slots for global variables                (=0b01)
//	2 - stack memory                              (=0b10)
//	3 - heap memory                               (=0b11)
//
// In each pointer 2 highest bits are reserved for
// encoding memory region. Zero pointer value is used
// as nil.
//
// Strings are recognized as a pair of pointer + length
// (in bytes). Each function that accepts string as an
// argument (at high-level) must expect two argument
// slots used for a single string argument. First slot
// of the pair must contain string pointer and second
// its length.
//
// PRINT instruction accepts string pair of pointer +
// length and pops them from EVS in direct order,
// meaning pointer will be poped first and length second.
type VM struct {
	// evaluation stack
	evs []uint32

	// stack memory which is used for:
	//	- call arguments passing
	//	- storing local variables
	smem []uint32

	// function call stack
	frames []frame

	// contains information about all program functions
	funs []finfo

	// instruction pointer (=index into text)
	ip ia

	// frame base pointer (=index into stack memory)
	bp uint32

	// number of pushed call arguments
	na uint32

	// number of expected return values
	// nr uint32
}

type frame struct {
	// return address (=index into program text)
	addr ia

	// base pointer (=index into vm stack memory)
	base uint32
}

func (g *VM) Exec(prog *Prog, args []uint32) {
	if len(prog.text) == 0 || len(prog.funs) == 0 {
		return
	}
	g.reset(prog, args)
	g.exec(prog.text)
}

func (g *VM) reset(prog *Prog, args []uint32) {
	g.ip = 0
	g.bp = 0
	g.na = 0
	// g.total = 0

	g.evs = g.evs[:0]
	g.smem = g.smem[:0]
	g.frames = g.frames[:0]
	// m.rom = rom
	g.funs = prog.funs

	// setup start function initial frame
	start := prog.funs[0]
	g.allocn(start.nloc)
	// g.nr = start.nret
	if uint32(len(args)) != start.narg {
		panic(fmt.Sprintf("number of start arguments mismatch (got %d, want %d)", len(args), start.narg))
	}
}

func (g *VM) exec(text []inst) {
	for g.ip < ia(len(text)) {
		x := text[g.ip]
		g.ip += 1

		switch x.op {
		case opTRAP:
			panic(fmt.Sprintf("TRAP (0x%04X)", g.ip-1))
		case opHALT:
			return
		case opCALL:
			g.call(x.arg)
		case opRET:
			if len(g.frames) == 0 {
				// stop execution since we are exiting
				// start function frame
				return
			}
			// if g.nr != uint32(len(g.evs)) {
			// 	panic(fmt.Sprintf("number of return values mismatch (got %d, want %d)", len(g.evs), g.nr))
			// }

			i := len(g.frames) - 1
			frame := g.frames[i]
			g.frames = g.frames[:i]

			g.smem = g.smem[:g.bp]
			g.ip = frame.addr
			g.bp = frame.base
		case opJUMP:
			g.ip = ia(x.arg)
		case opJZ:
			v := g.pop()
			if v == 0 {
				g.ip = ia(x.arg)
			}
		case opJNZ:
			v := g.pop()
			if v != 0 {
				g.ip = ia(x.arg)
			}
		case opTAKE:
			v := g.pop()
			g.pusharg(v)
		case opDRAIN:
			g.pop()
		case opPUSH:
			name := x.arg
			v := g.get(name)
			g.push(v)
		case opPOP:
			name := x.arg
			v := g.pop()
			g.set(name, v)
		case opMOVE:
			a := g.pop()
			v := g.pop()
			g.move(a, v)
		case opADDR:
			name := x.arg
			a := g.addr(name)
			g.push(a)
		case opINDEX:
			index := g.pop()
			addr := g.pop()
			stride := x.arg
			g.index(addr, index, stride)
		case opPUT:
			v := x.arg
			g.push(v)
		case opADD:
			b := g.pop()
			a := g.pop()
			g.push(a + b)
		case opSUB:
			b := g.pop()
			a := g.pop()
			g.push(a - b)
		case opMUL:
			b := g.pop()
			a := g.pop()
			g.push(a * b)
		case opEQUAL:
			b := g.pop()
			a := g.pop()
			var v uint32
			if a == b {
				v = 1
			}
			g.push(v)
		case opLESS:
			b := g.pop()
			a := g.pop()
			var v uint32
			if a < b {
				v = 1
			}
			g.push(v)
		default:
			panic(fmt.Sprintf("ins %d (0x%04X)", x.op, g.ip-1))
		}
	}
}

func (g *VM) index(addr uint32, index uint32, stride uint32) {
	typ := addr >> 30
	switch typ {
	case ptrROM:
		panic("stub")
	case ptrVAR:
		panic("stub")
	case ptrLOC:
		slot := stride*index + (addr<<2)>>4
		for i := range stride {
			g.push(g.smem[slot+i])
		}
	default:
		panic("stub")
	}
}

// call function under specified name (=index into funs)
func (g *VM) call(name uint32) {
	f := g.funs[name]
	if g.na != f.narg {
		panic("number of arguments mismatch")
	}

	g.frames = append(g.frames, frame{
		addr: g.ip,
		base: g.bp,
	})
	g.bp = uint32(len(g.smem)) - g.na
	g.allocn(f.nloc)
	g.na = 0
	// g.nr = f.nret
	g.ip = f.addr
}

// push call argument to stack memory
func (g *VM) pusharg(v uint32) {
	g.smem = append(g.smem, v)
	g.na += 1
}

// push value onto EVS
func (g *VM) push(v uint32) {
	g.evs = append(g.evs, v)
}

// pop value from EVS
func (g *VM) pop() uint32 {
	i := len(g.evs) - 1
	v := g.evs[i]
	g.evs = g.evs[:i]
	return v
}

// set local variable/argument to specified value
func (g *VM) set(name uint32, v uint32) {
	g.smem[g.bp+name] = v
}

// get value stored in local variable/argument
func (g *VM) get(name uint32) uint32 {
	return g.smem[g.bp+name]
}

// get encoded address of a local variable
func (g *VM) addr(name uint32) uint32 {
	a := g.bp + name
	return (ptrLOC << 30) | (a << 2)
}

// store value at memory address
func (g *VM) move(addr uint32, v uint32) {
	typ := addr >> 30
	switch typ {
	case ptrROM:
		panic("rom write")
	case ptrVAR:
		panic("stub")
	case ptrLOC:
		slot := (addr << 2) >> 4
		g.smem[slot] = v
	default:
		panic("stub")
	}
}

// allocate n local variables on stack memory
func (g *VM) allocn(n uint32) {
	if n == 0 {
		return
	}
	g.smem = append(g.smem, make([]uint32, n)...)
}

func (g *VM) PrintState() {
	fmt.Printf("evs:  %v\n", g.evs)
	fmt.Printf("smem: %v\n", g.smem)
}
