package vm

import (
	"fmt"

	"github.com/mebyus/ku/goku/vm/opc"
)

func (m *Machine) execCall(lt uint8) (uint64, error) {
	var data []byte // instruction data
	var size uint64
	var addr uint32 // call address
	var err error

	switch opc.Layout(lt) {
	case opc.CallReg:
		panic("stub")
	case opc.CallVal32:
		size = 4
		data, err = m.idata(size)
		if err != nil {
			return 0, err
		}
		addr = val32(data)
	default:
		return 0, fmt.Errorf("unknown layout 0x%02X", lt)
	}
	if err != nil {
		return 0, err
	}

	if addr >= uint32(len(m.text)) {
		return 0, fmt.Errorf("address 0x%08X outside of text segment", addr)
	}
	m.doCall(size, addr)
	return size, nil
}

func (m *Machine) doCall(size uint64, addr uint32) {
	m.frames = append(m.frames, Frame{
		Base: uint32(m.fp),

		// TODO: we probably need to check return address as uint64 first for correctness
		Ret: uint32(m.ip) + 2 + uint32(size), // return to next instruction after the call
	})
	m.fp = m.sp
	m.ip = uint64(addr)
	m.jump = true

	if uint64(len(m.frames)) > m.stats.MaxFrames {
		m.stats.MaxFrames = uint64(len(m.frames))
	}
}

func (m *Machine) execRet(lt uint8) (uint64, error) {
	if len(m.frames) == 0 {
		return 0, fmt.Errorf("empty frame stack")
	}
	frame := m.frames[len(m.frames)-1]
	m.frames = m.frames[:len(m.frames)-1]

	m.sp = m.fp
	m.fp = uint64(frame.Base)
	m.ip = uint64(frame.Ret)
	m.jump = true
	return 0, nil
}
