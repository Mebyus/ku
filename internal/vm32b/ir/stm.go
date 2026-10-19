package ir

type Stm interface {
	_stm()
}

type stm struct{}

func (stm) _stm() {}

type Discard struct {
	stm

	// Always not nil.
	Exp Exp
}

var _ Stm = &Discard{}

type Return struct {
	stm

	// Can be nil, if function returns nothing.
	Exp Exp
}

var _ Stm = &Return{}

type Assign struct {
	stm

	Exp Exp

	// each local variable uses zero-based index as name
	Name uint32
}

var _ Stm = &Assign{}

// if without else branch
type IfBlock struct {
	stm

	Body []Stm
	Exp  Exp
}

var _ Stm = &IfBlock{}

type While struct {
	stm

	Body []Stm

	Exp Exp
}

var _ Stm = &While{}
