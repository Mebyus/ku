package vm

import (
	"github.com/mebyus/ku/goku/vm/opc"
)

func (m *Machine) doJump(addr uint32) {
	m.ip = uint64(addr)
	m.jump = true
}

func (m *Machine) execJump(lt uint8) (uint64, *RuntimeError) {
	flag, variant := opc.DecodeJumpLayout(lt)

	var size uint64
	var addr uint32 // jump address
	switch variant {
	case opc.JumpReg:
		size = 1
		data, err := m.idata(size)
		if err != nil {
			return 0, err
		}
		r := opc.Register(data[0])
		val, err := m.get(r)
		if err != nil {
			return 0, err
		}
		segment, offset := getPointerSegmentAndOffset(val)
		if segment != SegText {
			return 0, &RuntimeError{
				Code: ErrorNonTextJump,
				Aux:  val,
			}
		}
		addr = offset
	case opc.JumpVal32:
		size = 4
		data, err := m.idata(size)
		if err != nil {
			return 0, err
		}
		addr = val32(data)
	default:
		return 0, &RuntimeError{
			Code: ErrorBadVariant,
			Aux:  uint64(variant),
		}
	}

	if addr >= uint32(len(m.text)) {
		return 0, &RuntimeError{
			Code: ErrorBadJumpAddress,
			Aux:  uint64(addr),
		}
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

// Bit flags in CF register.
const (
	FlagZero = 1 << 0
)

func (m *Machine) checkFlag(flag opc.JumpFlag) (bool, *RuntimeError) {
	switch flag {
	case 0:
		return true, nil
	case opc.FlagZ:
		return m.cf&FlagZero != 0, nil
	case opc.FlagNZ:
		return m.cf&FlagZero == 0, nil
	default:
		return false, &RuntimeError{
			Code: ErrorBadJumpFlag,
			Aux:  uint64(flag),
		}
	}
}
