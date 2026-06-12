package stg

import (
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/sx"
)

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
	Pin sx.Pin

	// Scope where this symbol was defined.
	Scope *Scope

	// Type which symbol yeilds on usage as value.
	Type *Type

	// Auxiliary symbol information. Used for different purposes during graph
	// construction.
	//
	// After indexing phase this field contains corresponding AST node index
	// for nodes inside unit scope.
	//
	// For builtin generic function contains its kind.
	Aux uint32

	Flags SymbolFlag

	Kind symk.Kind
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

func (s *Symbol) IsStub() bool {
	return s.Flags&SymbolStub != 0
}

type SymDef interface {
	// Discriminator method for interface implementations.
	// Only serves as a trick to enhance Go typechecking in
	// type assertions.
	//
	// Does nothing when called.
	_symdef()
}

type symdef struct{}

// Explicit interface implementation check.
var _ SymDef = symdef{}

func (symdef) _symdef() {}

// StaticValue is a SymDef for constant symbols. It holds a value known
// at compile time.
type StaticValue struct {
	symdef

	Exp Exp
}
