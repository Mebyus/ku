package genc

import (
	"cmp"
	"slices"

	"github.com/mebyus/ku/goku/compiler/srcmap"
	"github.com/mebyus/ku/goku/kub/eval"
)

type State struct {
	Map srcmap.PinMap

	// Maps type name to its id.
	types NameBook

	// Maps error name to its id.
	errors NameBook

	// Maps enum name to its values.
	enums map[string]*NameBook

	Env *eval.Env

	// If true then generate code for statements marked as debug.
	Debug bool

	// If true then generate code for test runner functions.
	Test bool
}

func (s *State) Init() {
	s.types.Init()
	s.errors.Init()
	s.enums = make(map[string]*NameBook)
}

func (s *State) GetTypeId(name string) uint64 {
	return s.types.Get(name)
}

func (s *State) GetErrorId(name string) uint64 {
	return s.errors.Get(name)
}

func (s *State) GetEnumValue(name, entry string) uint64 {
	enum, ok := s.enums[name]
	if !ok {
		enum = &NameBook{}
		enum.Init()
		s.enums[name] = enum
	}

	return enum.Get(entry)
}

func (s *State) ErrorRecords() []BookRecord {
	return s.errors.Records()
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

type BookRecord struct {
	Name string
	Id   uint64
}

// List returns all stored book records sorted by their ids.
func (b *NameBook) Records() []BookRecord {
	if len(b.m) == 0 {
		return nil
	}

	records := make([]BookRecord, 0, len(b.m))
	for name, id := range b.m {
		records = append(records, BookRecord{
			Name: name,
			Id:   id,
		})
	}
	slices.SortFunc(records, func(a, b BookRecord) int {
		return cmp.Compare(a.Id, b.Id)
	})
	return records
}
