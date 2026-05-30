package stg

import (
	"github.com/mebyus/ku/internal/ku/enums/scok"
	"github.com/mebyus/ku/internal/ku/enums/symk"
	"github.com/mebyus/ku/internal/ku/sx"
)

type Scope struct {
	// List of all symbols defined inside this scope. Symbols are
	// listed in order they appear in source code (except for global and unit scopes).
	Symbols []*Symbol

	// Types *TypeIndex

	Parent *Scope

	// Types *TypeIndex

	// Gens *GenIndex

	// Symbol map. Maps name to its local symbol.
	m map[ /* symbol name */ string]*Symbol

	Kind scok.Kind
}

func (s *Scope) InitGlobal() {
	s.init(scok.Global)
	// s.Types = types
	// s.Gens = gens
}

func (s *Scope) Init(kind scok.Kind, parent *Scope) {
	s.Parent = parent
	// s.Types = parent.Types
	// s.Gens = parent.Gens
	// s.Level = parent.Level + 1
	// s.LoopLevel = parent.LoopLevel

	// switch kind {
	// case sck.Loop:
	// 	s.LoopLevel += 1
	// }

	s.init(kind)
}

func (s *Scope) init(kind scok.Kind) {
	s.Kind = kind
	s.m = make(map[string]*Symbol)
}

// Get finds a symbol by its name inside the scope. Does not check parent scope.
func (s *Scope) Get(name string) *Symbol {
	return s.m[name]
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

func (s *Scope) Bind(symbol *Symbol) {
	// if s.IsLocal() {
	// 	symbol.Flags |= SymbolLocal
	// }
	symbol.Scope = s
	s.Symbols = append(s.Symbols, symbol)
	s.m[symbol.Name] = symbol
}

// New allocates new symbol inside the scope.
func (s *Scope) New(kind symk.Kind, name string, pin sx.Pin) *Symbol {
	symbol := &Symbol{
		Name: name,
		Pin:  pin,
		Kind: kind,
	}
	s.Bind(symbol)
	return symbol
}
