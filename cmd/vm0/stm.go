package main

type Function struct {
	body []stm
	narg uint32 // number of required arguments
	nret uint32 // number of returned values
}

// represents a language statement
type stm any

type Assign struct {
	exp Exp

	// each local variable uses zero-based index as name
	name uint32
}

var _ stm = &Assign{}

type Return struct {
	exp Exp
}

var _ stm = &Return{}

type Print struct {
	exp Exp
}

var _ stm = &Print{}

// if without else branch
type Ifthen struct {
	body []stm
	exp  Exp
}

var _ stm = &Ifthen{}

type While struct {
	body []stm
	exp  Exp
}

var _ stm = &Ifthen{}

func fun(narg, nret uint32, body []stm) *Function {
	return &Function{narg: narg, nret: nret, body: body}
}

func asn(name uint32, exp Exp) *Assign {
	return &Assign{name: name, exp: iv(exp)}
}

func ret(exp Exp) *Return {
	return &Return{exp: iv(exp)}
}

func prt(exp Exp) *Print {
	return &Print{exp: iv(exp)}
}

func ift(exp Exp, body []stm) *Ifthen {
	return &Ifthen{exp: iv(exp), body: body}
}

func while(exp Exp, body []stm) *While {
	return &While{exp: iv(exp), body: body}
}
