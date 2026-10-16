package main

type Exp any

// interger constant used as expression
type Integer struct {
	Val int32
}

var _ Exp = &Integer{}

type String struct {
	Val string
}

var _ Exp = &String{}

// local variable used as expression
type Local struct {
	Name uint32
}

var _ Exp = &Local{}

// binary expression
type BinExp struct {
	A, B Exp
	Op   int
}

var _ Exp = &BinExp{}

type Call struct {
	args []Exp

	// function name (=index)
	name uint32
}

var _ Exp = &Call{}

// binary operators
const (
	// a + b
	bopAdd = iota

	// a - b
	bopSub

	// a * b
	bopMul

	// a == b
	//
	// false -> 0
	// true  -> 1
	bopEql

	// a < b
	bopLess
)

// turns plain int into Integer node
func iv(exp Exp) Exp {
	n, ok := exp.(int)
	if ok {
		return &Integer{Val: int32(n)}
	}
	s, ok := exp.(string)
	if ok {
		return &String{Val: s}
	}
	return exp
}

func add(a, b Exp) *BinExp {
	return &BinExp{A: iv(a), B: iv(b), Op: bopAdd}
}

func sub(a, b Exp) *BinExp {
	return &BinExp{A: iv(a), B: iv(b), Op: bopSub}
}

func mul(a, b Exp) *BinExp {
	return &BinExp{A: iv(a), B: iv(b), Op: bopMul}
}

func eql(a, b Exp) *BinExp {
	return &BinExp{A: iv(a), B: iv(b), Op: bopEql}
}

func less(a, b Exp) *BinExp {
	return &BinExp{A: iv(a), B: iv(b), Op: bopLess}
}

func val(v int32) *Integer {
	return &Integer{Val: v}
}

func loc(name uint32) *Local {
	return &Local{Name: name}
}

func call(name uint32, args ...Exp) *Call {
	// TODO: handle strings
	for i := range len(args) {
		args[i] = iv(args[i])
	}
	return &Call{name: name, args: args}
}
