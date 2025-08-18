package vm

import (
	"fmt"

	"github.com/mebyus/ku/goku/vm/opc"
)

func (m *Machine) execSet(lt uint8) (uint64, error) {
	layout, v := opc.DecodeSetLayout(lt)

	var size uint64
	var val uint64 // set value
	var err error
	switch layout {
	case opc.SetReg:
		size = 2
	case opc.SetVal4:
		size = 1
	case opc.SetVal8:
		size = 2
	case opc.SetVal16:
		size = 3
	case opc.SetVal32:
		size = 5
	case opc.SetVal64:
		size = 9
	default:
		return 0, fmt.Errorf("unknown layout 0x%02X", uint8(layout))
	}
	if err != nil {
		return 0, err
	}

	err = m.set(0, val)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (m *Machine) getSetValReg() (uint64, error) {

}
