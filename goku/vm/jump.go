package vm

import (
	"fmt"

	"github.com/mebyus/ku/goku/vm/opc"
)

func (m *Machine) doJump(addr uint32) {
	m.ip = uint64(addr)
	m.jump = true
}

func (m *Machine) execJump(lt uint8) (uint64, error) {
	flag, layout := opc.DecodeJumpLayout(lt)
	size, err := opc.GetJumpDataSize(layout)
	if err != nil {
		return 0, err
	}
	addr, err := m.getJumpAddress(layout, size)
	if err != nil {
		return 0, err
	}
	if addr >= uint32(len(m.text)) {
		return 0, fmt.Errorf("jump to 0x%08X is out of program text range", addr)
	}
	c, err := m.checkFlag(flag)
	if err != nil {
		return 0, err
	}

	if c {
		m.doJump(addr)
	}
	return size, nil
}

func (m *Machine) getJumpAddress(layout opc.Layout, size uint64) (uint32, error) {
	data, err := m.idata(size)
	if err != nil {
		return 0, err
	}

	switch layout {
	case opc.JumpReg:
		r := opc.Register(data[0])
		val, err := m.get(r)
		if err != nil {
			return 0, err
		}
		segment, offset := getPointerSegmentAndOffset(val)
		if segment != SegText {
			return 0, fmt.Errorf("jump to pointer from (=0x%02X) segment", segment)
		}
		return offset, nil
	case opc.JumpVal32:
		val := val32(data)
		return val, nil
	default:
		return 0, fmt.Errorf("unexpected layout (=0x%02X)", layout)
	}
}

// Bit flags in CF register.
const (
	FlagZero = 1 << 0
)

func (m *Machine) checkFlag(flag opc.Flag) (bool, error) {
	switch flag {
	case 0:
		return true, nil
	case opc.FlagZ:
		return m.cf&FlagZero != 0, nil
	case opc.FlagNZ:
		return m.cf&FlagZero == 0, nil
	default:
		return false, fmt.Errorf("unexpected flag ?(%s) (=0x%02X)", flag, flag)
	}
}
