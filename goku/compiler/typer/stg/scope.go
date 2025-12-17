package stg

import (
	"github.com/mebyus/ku/goku/compiler/enums/sck"
	"github.com/mebyus/ku/goku/compiler/enums/smk"
	"github.com/mebyus/ku/goku/compiler/srcmap"
)

type Scope struct {
	// List of all symbols defined inside this scope. Symbols are
	// listed in order they appear in source code (except for global and unit scopes).
	Symbols []*Symbol

	Parent *Scope

	Types *TypeIndex

	// Symbol map. Maps name to its local symbol.
	m map[string]*Symbol

	// Unit where this scope is located.
	// Always nil for global scope.
	// Unit *Unit

	// Scope's nesting level. Starts from 0 for global scope. Language structure
	// implies that first levels are dependant on Kind:
	//
	//	- global     -> 0
	//	- unit       -> 1
	//	- test       -> 2
	//	- unsafe     -> 2
	//	- node       -> 3
	//
	// Subsequent levels are created inside function and method bodies by means of
	// various language constructs.
	//
	// First levels of scope hierarchy are fixed:
	//
	//	- global
	//	- unit
	//	- node
	//
	// Next levels may vary based on source code that defines the scope.
	Level uint32

	// How many loops deep this scope resides (within a function or method).
	// Starts from 0 for top scope (where there are no surrounding loops).
	LoopLevel uint32

	// How many logic branches deep (excluding loops) this scope resides (within a function or method).
	// Starts from 0 for top scope (where there are no surrounding branches).
	BranchLevel uint32

	Kind sck.Kind
}

func (s *Scope) Init(kind sck.Kind, parent *Scope) {
	if kind == sck.Global {
		if parent != nil {
			panic("global scope with parent")
		}
	} else {
		if parent == nil {
			panic("no parent scope")
		}
		s.Types = parent.Types
	}

	s.Kind = kind
	s.Parent = parent

	s.m = make(map[string]*Symbol)
}

// Has checks if symbol with a given name already exists inside the scope.
func (s *Scope) Has(name string) bool {
	return s.m[name] != nil
}

// Alloc allocates new symbol inside the scope.
func (s *Scope) Alloc(kind smk.Kind, name string, pin srcmap.Pin) *Symbol {
	symbol := &Symbol{
		Name: name,
		Pin:  pin,
		Kind: kind,
	}
	s.Bind(symbol)
	return symbol
}

func (s *Scope) Bind(symbol *Symbol) {
	symbol.Scope = s
	s.Symbols = append(s.Symbols, symbol)
	s.m[symbol.Name] = symbol
}

// Lookup finds a symbol by its name inside the scope or by doing lookup in
// parent scope.
func (s *Scope) Lookup(name string) *Symbol {
	symbol := s.Get(name)
	if symbol != nil {
		return symbol
	}
	if s.Parent != nil {
		return s.Parent.Lookup(name)
	}
	return nil
}

// Get finds a symbol by its name inside the scope. Does not check parent scope.
func (s *Scope) Get(name string) *Symbol {
	return s.m[name]
}
