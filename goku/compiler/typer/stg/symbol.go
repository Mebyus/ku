package stg

import (
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/source"
)

// Symbol. Symbols represent objects in source code which we can reference and
// use by using their names (identifiers). Most symbols are created with custom names
// in user's code by using language constructs:
//
//   - var - variable
//   - let - runtime constant, i.e. variable that is read-only, value is assigned to it exactly once
//   - const - build-time constant, i.e. its value must be known or computable at build-time
//   - type - defines a new type or prototype
//   - fun - defines a new function, method, blueprint or prototype blueprint
//   - method - defines a method, can only be used at unit level
//   - import - binds other unit to a local name in current unit
//   - function or prototype parameters
//
// Builtin symbols are created by the compiler and include basic types and primitive functions.
//
// Each symbol has a name, type and parent scope. The latter is determined by the place of origin (declaration)
// of the symbol.
//
// For each unique symbol in a program this struct must be instanced exactly once and then
// passed around in a pointer without creating a dereferenced copy.
type Symbol struct {
	// Symbol definition. Actual value stored in this field depends on Kind field.
	// See implementations of SymDef interface for more information.
	//
	// This field can be nil for dried symbol. Typical case for this is separate units compilation.
	//
	// During indexing and type checking this field may contain
	// temporary intermidiate values. For example for most symbols
	// during text gathering it is set to symbol index inside the box
	// of corresponding kind.
	Def SymDef

	// Always not empty.
	// Always an alphanumerical word for all symbol types except methods.
	//
	// For methods this field has special format:
	//	"receiver.name"
	//
	// Since other symbol types cannot have period in their names and
	// each custom type method names must be unique this
	// naming scheme cannot have accidental collisions.
	Name string

	// Source position of symbol origin (where this symbol was declared).
	Pin source.Pin

	// Always not nil in a completed graph. Can be nil during tree graph.
	Type *Type

	// Scope where this symbol was defined.
	Scope *Scope

	// Auxiliary symbol information. Used for different purposes during graph
	// construction.
	//
	// After indexing phase this field contains corresponding AST node index
	// for nodes inside unit scope.
	Aux uint32

	Flags SymbolFlag

	Kind smk.Kind
}

// SymbolFlag bit flags for specifing additional symbol properties.
type SymbolFlag uint8

const (
	// Symbol is language builtin.
	SymbolBuiltin SymbolFlag = 1 << iota

	// Symbol is function stub.
	SymbolStub
)

type SymDef interface {
	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_symdef()
}
