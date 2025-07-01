package genc

import "github.com/mebyus/ku/goku/compiler/srcmap"

type State struct {
	Map srcmap.PinMap

	// Maps type name to its id.
	types NameBook

	// Maps error name to its id.
	errors NameBook

	// If true then generate code for statements marked as debug.
	Debug bool

	// If true then generate code for test runner functions.
	Test bool
}

func (s *State) Init() {
	s.types.Init()
	s.errors.Init()
}

func (s *State) GetTypeId(name string) uint64 {
	return s.types.Get(name)
}

func (s *State) GetErrorId(name string) uint64 {
	return s.errors.Get(name)
}

// NameBook generates sequential integer ids for each unique given name (string).
// Generated ids start from 1.
type NameBook struct {
	// Maps name to its id.
	m map[string]uint64

	// Previous generated id.
	prev uint64
}

func (b *NameBook) Init() {
	b.m = make(map[string]uint64)
}

func (b *NameBook) Get(name string) uint64 {
	id, ok := b.m[name]
	if ok {
		return id
	}
	b.prev += 1
	id = b.prev
	b.m[name] = id
	return id
}
