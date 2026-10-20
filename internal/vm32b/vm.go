package vm32b

import "fmt"

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
