package ir

type Exp interface {
	_exp()
}

type exp struct{}

func (exp) _exp() {}
