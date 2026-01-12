package stg

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/sm"
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
	// temporary intermidiate values.
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

	// Link name. Not empty only if it differs from standard link name
	// mangling algorithm.
	// Link string

	// Source position of symbol origin (where this symbol was declared).
	Pin sm.Pin

	// Always not nil in a completed graph. Can be nil during tree graph.
	Type *Type

	// Scope where this symbol was defined.
	Scope *Scope

	// Auxiliary symbol information. Used for different purposes during graph
	// construction.
	//
	// After indexing phase this field contains corresponding AST node index
	// for nodes inside unit scope.
	//
	// For builtin generic function contains its kind.
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

	// Symbol is declared as public.
	// Applicable only to unit-level symbols.
	SymbolPublic

	// Symbol should be skipped during compilation.
	// Flag is set as a result of symbol usage analysis and program tree pruning.
	SymbolSkip

	// Only applicable for functions. Symbol should be exported in produced binary object file.
	SymbolExport

	// Local symbols are the ones created in function or method bodies (as well as their parameters).
	// Non-local symbols are global or unit-level.
	SymbolLocal
)

func (s *Symbol) IsPublic() bool {
	return s.Flags&SymbolPublic != 0
}

func (s *Symbol) IsExport() bool {
	return s.Flags*SymbolExport != 0
}

func (s *Symbol) IsLocal() bool {
	return s.Flags&SymbolLocal != 0
}

func (s *Symbol) MarkSkip() {
	s.Flags |= SymbolSkip
}

func (s *Symbol) ShouldSkip() bool {
	return s.Flags&SymbolSkip != 0
}

// GetMethodName get original method name without receiver prefix.
// Panics on symbols that are not methods.
func (s *Symbol) GetMethodName() string {
	if s.Kind != smk.Method {
		panic(fmt.Sprintf("%s (=%d) symbol \"%s\"", s.Kind, s.Kind, s.Name))
	}

	for i := range len(s.Name) {
		c := s.Name[i]
		if c == '.' {
			name := s.Name[i+1:]
			if name == "" {
				break
			}
			return name
		}
	}

	panic(fmt.Sprintf("method \"%s\" has bad symbol name format", s.Name))
}

type SymDef interface {
	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_symdef()
}

type nodeSymDef struct{}

func (nodeSymDef) _symdef() {}

// StaticValue is a SymDef for constant symbols. It holds a value known
// at compile time.
type StaticValue struct {
	nodeSymDef

	Exp Exp
}

// SymDefType symbol definition for symbols which refer to a type.
type SymDefType struct {
	nodeSymDef

	Type *Type
}

// SymDefUnit symbol definition for import symbols.
type SymDefUnit struct {
	nodeSymDef

	Unit *Unit
}
