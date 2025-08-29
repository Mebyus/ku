package vm

import (
	"encoding/binary"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/mebyus/ku/goku/vm/kvx"
	"github.com/mebyus/ku/goku/vm/opc"
)

type Frame struct {
	// Return address.
	Ret uint32

	// Frame pointer.
	Base uint32
}

// Memory segment names.
//
// Segment is used to store a pointer in register or memory.
// Stored pointer is always 8 bytes in size. Its highest byte
// stores segment, other 7 bytes store offset into that segment:
//
//	SS XX XX XX XX XX XX XX
const (
	SegText   = 0x00
	SegData   = 0x01
	SegGlobal = 0x02
	SegStack  = 0x03
	SegHeap   = 0x04
)

type Machine struct {
	// Instruction pointer. Index in text memory.
	ip uint64

	// Stack pointer. Index in stack memory.
	sp uint64

	// Frame pointer. Index in stack memory.
	fp uint64

	// Syscall register.
	// Select syscall number or receive result code.
	sc uint64

	// Comparison flags.
	cf uint64

	// Number of executed instructions.
	clock uint64

	// General-purpose registers.
	r [64]uint64

	// Code of the program being executed, size cannot change during execution.
	text []byte

	// Static, read-only program data. Loaded at program start.
	data []byte

	// Memory for global variables. Loaded and initialized at program start.
	global []byte

	// Stack memory, size cannot change during execution.
	stack []byte

	// Heap memory, size can change during execution.
	heap []byte

	// Stack for storing frames in procedure calls.
	frames []Frame

	// Runtime error occured while executing current instruction.
	err *RuntimeError

	stats Stats

	// Indicates if jump occured while executing current instruction.
	jump bool

	// Indicates if vm was halted by instruction or runtime error.
	halt bool
}

type Stats struct {
	MaxStackSize uint64
	MaxFrames    uint64
}

func (m *Machine) Exec(prog *kvx.Program) *Exit {
	if len(prog.Text) == 0 {
		// TODO: fill this errors
		return &Exit{Error: &RuntimeError{}}
	}
	if int(prog.EntryPoint) >= len(prog.Text) {
		return &Exit{Error: &RuntimeError{}}
	}

	m.ip = uint64(prog.EntryPoint)
	m.text = prog.Text
	m.data = prog.Data

	if int(prog.GlobalSize) > cap(m.global) {
		m.global = slices.Grow(m.global, int(prog.GlobalSize)-cap(m.global))
	}
	m.global = m.global[:prog.GlobalSize]
	clear(m.global)

	// reset vm state
	m.err = nil
	m.halt = false
	m.jump = false
	m.sp = 0
	m.fp = 0
	m.sc = 0
	m.cf = 0
	m.clock = 0
	m.stack = m.stack[:0]
	m.frames = m.frames[:0]
	m.heap = m.heap[:0]
	clear(m.r[:])

	start := time.Now()

	for !m.halt {
		m.step()
		m.clock += 1
	}

	return m.exit(time.Since(start))
}

func (m *Machine) step() {
	m.jump = false

	ip := m.ip
	if ip+2 > uint64(len(m.text)) {
		// every instruction needs at least 2 bytes for Opcode + Layout
		m.stop(&RuntimeError{Code: ErrorTextEnd})
		return
	}

	op := opc.Opcode(m.text[ip])
	lt := m.text[ip+1]

	// Each m.exec*** method returns instruction data size
	var size uint64 // instruction data size
	var err *RuntimeError
	switch op {
	case opc.Sys:
		err = m.execSys(lt)
	case opc.Jump:
		size, err = m.execJump(lt)
	case opc.Call:
		size, err = m.execCall(lt)
	case opc.Clear:
		// m.clear()
	case opc.Set:
		size, err = m.execSet(lt)
	case opc.Load:
	case opc.Store:
		// err = m.loadValReg()
		// case LoadRegReg:
		// 	err = m.loadRegReg()
		// case LoadValSysReg:
		// 	m.loadValSysReg()
	case opc.Inc:
	case opc.Add:
	case opc.Sub:
		// err = m.addRegReg()
		// err = m.incReg()
	case opc.Test:
		// err = m.testRegVal()
	// case JumpFlagAddr:
	// 	err = m.jumpFlagAddr()
	default:
		m.stop(&RuntimeError{
			Code: ErrorBadOpcode,
			Aux:  uint64(op),
		})
		return
	}
	if err != nil {
		m.stop(err)
		return
	}
	if m.halt {
		return
	}

	if !m.jump {
		m.ip += 2 + size
	}
}

// switch to halt state with runtime error
func (m *Machine) stop(err *RuntimeError) {
	m.err = err
	m.halt = true
}

// get n bytes of current instruction data (opcode and layout not included)
func (m *Machine) idata(n uint64) ([]byte, *RuntimeError) {
	ip := m.ip
	if ip+2+n > uint64(len(m.text)) {
		return nil, &RuntimeError{
			Code: ErrorBadInstructionDataLength,
			Aux:  n,
		}
	}
	return m.text[ip+2 : ip+2+n], nil
}

func getPointerSegmentAndOffset(ptr uint64) (uint8, uint32) {
	// most significant byte in pointer encodes the memory segment
	segment := uint8(ptr >> 56)
	offset := uint32(ptr & 0xFFFFFFFF)
	return segment, offset
}

func (m *Machine) memslice(ptr uint64, n uint32) ([]byte, error) {
	if n == 0 {
		return nil, fmt.Errorf("empty slice")
	}

	segment, offset := getPointerSegmentAndOffset(ptr)
	var b []byte
	switch segment {
	case SegText:
		b = m.text
	case SegData:
		b = m.data
	case SegGlobal:
		b = m.global
	case SegStack:
		b = m.stack
	case SegHeap:
		b = m.heap
	default:
		return nil, fmt.Errorf("unknown segment (=0x%02X) in pointer (=0x%016X)", segment, ptr)
	}

	if offset >= uint32(len(b)) {
		return nil, fmt.Errorf("offset (=0x%08X) is out of (=0x%02X) segment range", offset, segment)
	}
	end := offset + n
	if end > uint32(len(b)) {
		return nil, fmt.Errorf("end (=0x%08X) is out of (=0x%02X) segment range", end, segment)
	}
	return b[offset:end], nil
}

// get register value
func (m *Machine) get(r opc.Register) (uint64, *RuntimeError) {
	if !r.Special() {
		if r >= 64 {
			return 0, &RuntimeError{
				Code: ErrorBadRegister,
				Aux:  uint64(r),
			}
		}
		v := m.r[r]
		return v, nil
	}

	switch r {
	case opc.RegIP:
		return m.ip, nil
	case opc.RegSP:
		return m.sp, nil
	case opc.RegFP:
		return m.fp, nil
	case opc.RegSC:
		return m.sc, nil
	case opc.RegCF:
		return m.cf, nil
	case opc.RegClock:
		return m.clock, nil
	default:
		return 0, &RuntimeError{
			Code: ErrorBadSpecialRegister,
			Aux:  uint64(r),
		}
	}
}

// set register value
func (m *Machine) set(r opc.Register, v uint64) *RuntimeError {
	if !r.Special() {
		if r >= 64 {
			return &RuntimeError{
				Code: ErrorBadRegister,
				Aux:  uint64(r),
			}
		}
		m.r[r] = v
		return nil
	}

	switch r {
	case opc.RegIP, opc.RegSP, opc.RegFP, opc.RegCF, opc.RegClock:
		return &RuntimeError{
			Code: ErrorReadOnlyRegister,
			Aux:  uint64(r),
		}
	case opc.RegSC:
		m.sc = v
		return nil
	default:
		return &RuntimeError{
			Code: ErrorBadSpecialRegister,
			Aux:  uint64(r),
		}
	}
}

// Exit describes vm exit state after program execution.
// Includes both normal and abnormal exits.
type Exit struct {
	// Runtime error for abnormal exit. Equals nil for normal exit.
	Error *RuntimeError

	// Real execution time.
	Time time.Duration

	// Value of instruction pointer register.
	IP uint64

	// Exit status of the program. Obtained from syscall register
	// upon program exit. Valid only for normal exit.
	Status uint64

	// Number of executed instructions.
	Clock uint64

	// Maximum stack size during execution.
	MaxStackSize uint64

	// Maximum number of frames (call depth) during execution.
	MaxFrames uint64

	// Last executed instruction opcode.
	Opcode opc.Opcode

	// Last executed instruction layout. Equals 0 if execution
	// stopped because layout would be outside of text segment,
	Layout uint8
}

func (e *Exit) Render(w io.Writer) error {
	_, err := io.WriteString(w, fmt.Sprintf("vm.time:  %s\n", e.Time.String()))
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, fmt.Sprintf("vm.clock: %d\n", e.Clock))
	if err != nil {
		return err
	}

	s := e.String()
	_, err = io.WriteString(w, s)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "\n")
	if err != nil {
		return err
	}

	return nil
}

func (e *Exit) String() string {
	if e.Error == nil {
		return fmt.Sprintf("vm: normal exit (at 0x%08X) with status %d", e.IP, e.Status)
	}

	return fmt.Sprintf("vm: abnormal exit (at 0x%08X) with runtime error: %v", e.IP, e.Error)
}

func (m *Machine) exit(dur time.Duration) *Exit {
	e := &Exit{
		Time:  dur,
		IP:    m.ip,
		Clock: m.clock,

		MaxStackSize: m.stats.MaxStackSize,
		MaxFrames:    m.stats.MaxFrames,
	}

	if m.ip < uint64(len(m.text)) {
		e.Opcode = opc.Opcode(m.text[m.ip])
		if m.ip+1 < uint64(len(m.text)) {
			e.Layout = m.text[m.ip+1]
		}
	}

	if m.err != nil {
		e.Error = m.err
		return e
	}

	e.Status = m.sc
	return e
}

func val64(buf []byte) uint64 {
	return binary.LittleEndian.Uint64(buf)
}

func val32(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}

func val16(buf []byte) uint16 {
	return binary.LittleEndian.Uint16(buf)
}
