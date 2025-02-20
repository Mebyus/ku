package typer

import "github.com/mebyus/ku/goku/compiler/ast"

// Box is a container for collecting AST nodes from all unit texts.
type Box struct {
	// List of top custom type definition nodes.
	Types []ast.Type

	// List of top constant definition nodes.
	Constants []ast.TopConst

	// List of top variable definition nodes.
	Variables []ast.TopVar

	// List of top alias nodes.
	Aliases []ast.TopAlias

	// List of top function definition nodes.
	Functions []ast.Fun

	// List of unit test functions.
	Tests []ast.Fun

	// List of top function stub nodes.
	FunStubs []ast.FunStub

	// List of method nodes.
	Methods []ast.Method

	// List of generic nodes.
	Generics []ast.Gen

	// List of generic bind nodes.
	GenBinds []ast.GenBind

	// Maps custom type receiver name to a list of its method indices inside
	// Methods slice.
	MethodsByReceiver map[ /* receiver type name */ string][]uint32

	// Maps generic name to a list of its generic bind indices inside GenBinds slice.
	GenBindsByName map[ /* generic name */ string][]uint32
}

func (b *Box) init(texts []*ast.Text) {
	var (
		funs      uint32
		vars      uint32
		tests     uint32
		types     uint32
		aliases   uint32
		methods   uint32
		genbinds  uint32
		funstubs  uint32
		generics  uint32
		constants uint32
	)
	for _, t := range texts {
		funs += uint32(len(t.Functions))
		vars += uint32(len(t.Variables))
		tests += uint32(len(t.Tests))
		types += uint32(len(t.Types))
		aliases += uint32(len(t.Aliases))
		methods += uint32(len(t.Methods))
		genbinds += uint32(len(t.GenBinds))
		funstubs += uint32(len(t.FunStubs))
		generics += uint32(len(t.Generics))
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
	if aliases != 0 {
		b.Aliases = make([]ast.TopAlias, 0, aliases)
	}
	if methods != 0 {
		b.Methods = make([]ast.Method, 0, methods)
		b.MethodsByReceiver = make(map[string][]uint32)
	}
	if funstubs != 0 {
		b.FunStubs = make([]ast.FunStub, 0, funstubs)
	}
	if constants != 0 {
		b.Constants = make([]ast.TopConst, 0, constants)
	}
	if generics != 0 {
		b.Generics = make([]ast.Gen, 0, generics)
		b.GenBindsByName = make(map[string][]uint32, generics)
	}
	if genbinds != 0 {
		b.GenBinds = make([]ast.GenBind, 0, genbinds)
	}
}

// returns internal index of saved node
func (b *Box) addType(node ast.Type) uint32 {
	i := uint32(len(b.Types))
	b.Types = append(b.Types, node)
	return i
}

func (b *Box) addAlias(node ast.TopAlias) uint32 {
	i := uint32(len(b.Aliases))
	b.Aliases = append(b.Aliases, node)
	return i
}

func (b *Box) addFun(node ast.Fun) uint32 {
	i := uint32(len(b.Functions))
	b.Functions = append(b.Functions, node)
	return i
}

func (b *Box) addConst(node ast.TopConst) uint32 {
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

	receiver := node.Receiver.Name.Str
	b.bindMethod(receiver, i)
	return i
}

func (b *Box) addTest(node ast.Fun) uint32 {
	i := uint32(len(b.Tests))
	b.Tests = append(b.Tests, node)
	return i
}

func (b *Box) addGen(node ast.Gen) uint32 {
	i := uint32(len(b.Generics))
	b.Generics = append(b.Generics, node)
	return i
}

func (b *Box) addGenBind(node ast.GenBind) uint32 {
	i := uint32(len(b.GenBinds))
	b.GenBinds = append(b.GenBinds, node)

	generic := node.Name.Str
	b.bindGeneric(generic, i)
	return i
}

func (b *Box) bindMethod(receiver string, i uint32) {
	b.MethodsByReceiver[receiver] = append(b.MethodsByReceiver[receiver], i)
}

func (b *Box) bindGeneric(generic string, i uint32) {
	b.GenBindsByName[generic] = append(b.GenBindsByName[generic], i)
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

func (b *Box) Const(i uint32) ast.TopConst {
	return b.Constants[i]
}

func (b *Box) Var(i uint32) ast.TopVar {
	return b.Variables[i]
}

func (b *Box) Method(i uint32) ast.Method {
	return b.Methods[i]
}
