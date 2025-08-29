package vm

import (
	"github.com/mebyus/ku/goku/vm/opc"
)

func (m *Machine) execSet(lt uint8) (uint64, *RuntimeError) {
	variant, v := opc.DecodeSetLayout(lt)

	var data []byte // instruction data
	var size uint64
	var val uint64 // set value
	var err *RuntimeError
	switch variant {
	case opc.SetReg:
		size = 2
		data, err = m.idata(size)
		if err != nil {
			return 0, err
		}
		r := opc.Register(data[1])
		val, err = m.get(r)
	case opc.SetVal4:
		val = uint64(v)
		size = 1
		data, err = m.idata(size)
	case opc.SetVal8:
		size = 2
		data, err = m.idata(size)
		if err != nil {
			return 0, err
		}
		val = uint64(data[1])
	case opc.SetVal16:
		size = 3
		data, err = m.idata(size)
		if err != nil {
			return 0, err
		}
		val = uint64(val16(data[1:]))
	case opc.SetVal32:
		size = 5
		data, err = m.idata(size)
		if err != nil {
			return 0, err
		}
		val = uint64(val32(data[1:]))
	case opc.SetVal64:
		size = 9
		data, err = m.idata(size)
		if err != nil {
			return 0, err
		}
		val = val64(data[1:])
	default:
		return 0, &RuntimeError{
			Code: ErrorBadVariant,
			Aux:  uint64(variant),
		}
	}
	if err != nil {
		return 0, err
	}

	r := opc.Register(data[0])
	err = m.set(r, val)
	if err != nil {
		return 0, err
	}
	return size, nil
}
