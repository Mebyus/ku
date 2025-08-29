package vm

import "github.com/mebyus/ku/goku/vm/opc"

func (m *Machine) execSys(lt uint8) *RuntimeError {
	switch lt {
	case opc.Trap:
		return &RuntimeError{Code: ErrorTrap}
	case opc.Halt:
		m.halt = true
	case opc.Nop:
		// no operation
	case opc.Ret:
		return m.execRet()
	case opc.SysCall:
		panic("stub")
	default:
		return &RuntimeError{
			Code: ErrorTrap,
			Aux:  uint64(lt),
		}
	}

	return nil
}
