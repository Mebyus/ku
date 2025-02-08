package typer

import "github.com/mebyus/ku/goku/compiler/ast"

// Box is a container for collecting AST nodes from all unit texts.
type Box struct {
	// List of top custom type definition nodes.
	Types []ast.Type

	// List of top constant definition nodes.
	Constants []ast.TopLet

	// List of top variable definition nodes.
	Variables []ast.TopVar

	// List of top function definition nodes.
	Functions []ast.Fun

	// List of unit test functions.
	Tests []ast.Fun

	// List of top function declaration nodes.
	FunStubs []ast.FunStub

	// List of method nodes.
	Methods []ast.Method

	// Maps custom type receiver name to a list of its method indices inside
	// Methods slice.
	MethodsByReceiver map[ /* receiver type name */ string][]uint32
}

func (b *Box) init(texts []*ast.Text) {
	var (
		funs      uint32
		vars      uint32
		tests     uint32
		types     uint32
		methods   uint32
		funstubs  uint32
		constants uint32
	)
	for _, t := range texts {
		funs += uint32(len(t.Functions))
		vars += uint32(len(t.Variables))
		tests += uint32(len(t.Tests))
		types += uint32(len(t.Types))
		methods += uint32(len(t.Methods))
		funstubs += uint32(len(t.FunStubs))
		constants += uint32(len(t.Constants))
	}

	if funs != 0 {
		b.Functions = make([]ast.Fun, 0, funs)
	}
	if vars != 0 {
		b.Variables = make([]ast.TopVar, 0, vars)
	}
	if tests != 0 {
		b.Tests = make([]ast.Fun, 0, tests)
	}
	if types != 0 {
		b.Types = make([]ast.Type, 0, types)
	}
	if methods != 0 {
		b.Methods = make([]ast.Method, 0, methods)
		b.MethodsByReceiver = make(map[string][]uint32)
	}
	if funstubs != 0 {
		b.FunStubs = make([]ast.FunStub, 0, funstubs)
	}
	if constants != 0 {
		b.Constants = make([]ast.TopLet, 0, constants)
	}
}

// returns internal index of saved node
func (b *Box) addType(node ast.Type) uint32 {
	i := uint32(len(b.Types))
	b.Types = append(b.Types, node)
	return i
}

func (b *Box) addFun(node ast.Fun) uint32 {
	i := uint32(len(b.Functions))
	b.Functions = append(b.Functions, node)
	return i
}

func (b *Box) addConst(node ast.TopLet) uint32 {
	i := uint32(len(b.Constants))
	b.Constants = append(b.Constants, node)
	return i
}

func (b *Box) addVar(node ast.TopVar) uint32 {
	i := uint32(len(b.Variables))
	b.Variables = append(b.Variables, node)
	return i
}

func (b *Box) addFunStub(node ast.FunStub) uint32 {
	i := uint32(len(b.FunStubs))
	b.FunStubs = append(b.FunStubs, node)
	return i
}

func (b *Box) addMethod(node ast.Method) uint32 {
	i := uint32(len(b.Methods))
	b.Methods = append(b.Methods, node)
	return i
}

func (b *Box) addTest(node ast.Fun) uint32 {
	i := uint32(len(b.Tests))
	b.Tests = append(b.Tests, node)
	return i
}

func (b *Box) bindMethod(receiver string, i uint32) {
	b.MethodsByReceiver[receiver] = append(b.MethodsByReceiver[receiver], i)
}

func (b *Box) Type(i uint32) ast.Type {
	return b.Types[i]
}

func (b *Box) Fun(i uint32) ast.Fun {
	return b.Functions[i]
}

func (b *Box) FunStub(i uint32) ast.FunStub {
	return b.FunStubs[i]
}

func (b *Box) Test(i uint32) ast.Fun {
	return b.Tests[i]
}

func (b *Box) Const(i uint32) ast.TopLet {
	return b.Constants[i]
}

func (b *Box) Var(i uint32) ast.TopVar {
	return b.Variables[i]
}

func (b *Box) Method(i uint32) ast.Method {
	return b.Methods[i]
}
