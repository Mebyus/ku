package main

import (
	"fmt"
	"os"
)

// Implements stack-based bytecode virtual machine.
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
type vm struct {
	// evaluation stack
	evs []int32

	// function call stack
	frames []frame

	// contains information about all program functions
	funs []finfo

	// stack memory which is used for:
	//	- call arguments passing
	//	- storing local variables
	smem []int32

	// contains program read-only memory, typically
	// it carries string literals, may be nil if program
	// does not have rom
	rom []byte

	// total number of executed instructions
	total uint64

	// instruction pointer (=index into text)
	ip uint32

	// frame base pointer (=index into stack memory)
	bp uint32

	// number of pushed call arguments
	na uint32
}

const (
	ptrROM = 0b00
	ptrVAR = 0b01
)

type frame struct {
	// return address (=index into program text)
	addr uint32

	// base pointer (=index into vm stack memory)
	base uint32
}

// static information about compiled function
type finfo struct {
	// start instruction index
	addr uint32

	// number of arguments required to make a call
	narg uint32

	// number of allocated local variables
	// does not include arguments but does
	// include stack memory gaps
	nloc uint32

	// number of returned (via evs) values
	nret uint32
}

func (m *vm) reset(rom []byte, funs []finfo) {
	m.ip = 0
	m.bp = 0
	m.na = 0
	m.total = 0

	m.evs = m.evs[:0]
	m.smem = m.smem[:0]
	m.frames = m.frames[:0]
	m.rom = rom
	m.funs = funs

	m.allocn(funs[0].nloc)
}

func (m *vm) exec(prog *Program) {
	if len(prog.text) == 0 || len(prog.funs) == 0 {
		return
	}
	m.reset(prog.rom, prog.funs)

	for m.ip < uint32(len(prog.text)) {
		x := prog.text[m.ip]
		m.ip += 1
		m.total += 1

		switch x.typ {
		case opTRAP:
			panic(fmt.Sprintf("trap at %d", m.ip-1))
		case opHALT:
			return
		case opPRINT:
			p := uint32(m.pop())
			n := uint32(m.pop())

			var data []byte
			typ := p >> 30
			switch typ {
			case ptrROM:
				data = m.rom[p : p+n]
			case ptrVAR:
				panic("stub")
			default:
				panic(typ)
			}
			if len(data) != 0 {
				os.Stdout.Write(data)
			}
		case opCALL:
			m.call(uint32(x.arg))
		case opRET:
			if len(m.frames) == 0 {
				// stop execution since we are exiting
				// start function frame
				return
			}

			i := len(m.frames) - 1
			frame := m.frames[i]
			m.frames = m.frames[:i]

			m.smem = m.smem[:m.bp]
			m.ip = frame.addr
			m.bp = frame.base
		case opJUMP:
			m.ip = uint32(x.arg)
		case opJZ:
			v := m.pop()
			if v == 0 {
				m.ip = uint32(x.arg)
			}
		case opJNZ:
			v := m.pop()
			if v != 0 {
				m.ip = uint32(x.arg)
			}
		case opTAKE:
			v := m.pop()
			m.pusha(v)
		case opDRAIN:
			m.pop()
		case opPUSH:
			name := uint32(x.arg)
			v := m.get(name)
			m.push(v)
		case opPOP:
			name := uint32(x.arg)
			v := m.pop()
			m.set(name, v)
		case opPUT:
			v := x.arg
			m.push(v)
		case opADD:
			b := m.pop()
			a := m.pop()
			m.push(a + b)
		case opSUB:
			b := m.pop()
			a := m.pop()
			m.push(a - b)
		case opMUL:
			b := m.pop()
			a := m.pop()
			m.push(a * b)
		case opEQUAL:
			b := m.pop()
			a := m.pop()
			var v int32
			if a == b {
				v = 1
			}
			m.push(v)
		case opLESS:
			b := m.pop()
			a := m.pop()
			var v int32
			if a < b {
				v = 1
			}
			m.push(v)
		default:
			panic(fmt.Sprintf("ins %d at %d", x.typ, m.ip-1))
		}
	}
}

func (m *vm) push(v int32) {
	m.evs = append(m.evs, v)
}

// push call argument to stack memory
func (m *vm) pusha(v int32) {
	m.smem = append(m.smem, v)
	m.na += 1
}

// call function under specified name (=index into funs)
func (m *vm) call(name uint32) {
	f := m.funs[name]
	if m.na != f.narg {
		panic("number of arguments mismatch")
	}

	m.frames = append(m.frames, frame{
		addr: m.ip,
		base: m.bp,
	})
	m.bp = uint32(len(m.smem)) - m.na
	m.allocn(f.nloc)
	m.na = 0
	m.ip = f.addr
}

// allocate n local variables on stack memory
func (m *vm) allocn(n uint32) {
	if n == 0 {
		return
	}
	m.smem = append(m.smem, make([]int32, n)...)
}

func (m *vm) pop() int32 {
	i := len(m.evs) - 1
	v := m.evs[i]
	m.evs = m.evs[:i]
	return v
}

func (m *vm) set(name uint32, v int32) {
	m.smem[m.bp+name] = v
}

func (m *vm) get(name uint32) int32 {
	return m.smem[m.bp+name]
}

func (m *vm) pstate() {
	fmt.Printf("evs:  %v\n", m.evs)
	fmt.Printf("smem: %v\n", m.smem)
}
