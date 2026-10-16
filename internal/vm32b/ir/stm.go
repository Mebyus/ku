package ir

type Stm interface {
	_stm()
}

type stm struct{}

func (stm) _stm() {}

type Return struct {
	stm

	// Can be nil, if function returns nothing.
	Exp Exp
}

var _ Stm = &Return{}
