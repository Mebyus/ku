package stg

import (
	"github.com/mebyus/ku/goku/compiler/ast"
	"github.com/mebyus/ku/goku/compiler/diag"
)

// Pool is a collection of Typer objects. Client code can request available
// object from Pool or return an object that is no longer needed. Pool automatically
// manages objects allocation, reuse and resetting state.
type Pool struct {
	// Objects available for reuse. Operates as stack.
	free []*Typer

	// Maps unit to a Typer attached to it.
	// Contains only occupied objects. Typer is moved to free stack when
	// removed from this map.
	m map[*Unit]*Typer

	com *Common
}

func NewPool(c *Common) *Pool {
	return &Pool{
		m:   make(map[*Unit]*Typer),
		com: c,
	}
}

// Get returns a Typer object attached to a given unit.
// Allocates a new object if specified unit has no Typer attached to it.
func (p *Pool) Get(u *Unit) *Typer {
	t := p.m[u]
	if t != nil {
		return t
	}

	n := len(p.free)
	if n == 0 {
		t = NewTyper(p.com)
	} else {
		t = p.free[n-1]
		p.free = p.free[:n-1]
		t.reset()
	}

	t.unit = u
	p.m[u] = t
	return t
}

func (p *Pool) Put(t *Typer) {
	delete(p.m, t.unit)
	t.unit = nil
	p.free = append(p.free, t)
}

// State indicates translation stage of unit inside a Typer.
//
// Name of each state corresponds to current phase of translation.
type State uint8

const (
	// Initial translation phase.
	// Typer gathers source texts during this phase.
	// Phase ends when Index is called.
	//
	// Must be zero (for proper reset).
	StateAlloc State = iota

	StateIndex

	StateScan
)

// Typer is a high-level object that drives unit translation from source AST
// to STG form. It keeps track of various intermediate states, temporary buffers, etc.
//
// Typers can be reused for translating different units as long as unit translation
// was finished. Unit translation can be paused at certain "checkpoints" and resumed
// later, this is used for preloading symbol declarations from prelude units.
type Typer struct {
	box NodeBox

	gb GraphBuilder

	warns  []diag.Error
	errors []diag.Error

	// After index phase is complete, contains all functions defined inside unit.
	funs []*Symbol

	vars []*Symbol

	// After index phase is complete, contains all custom type symbols defined inside unit.
	types []*Symbol

	// After index phase is complete, contains all unit-level constants defined inside unit.
	consts []*Symbol

	// After index phase is complete, contains all methods defined inside unit.
	methods []*Symbol

	com *Common

	unit *Unit

	// Maps custom type symbol (receiver) to a list of its method symbols.
	// Filled during index phase.
	mr map[ /* receiver type symbol */ *Symbol][]*Symbol

	deps DepSet

	state State
}

func NewTyper(c *Common) *Typer {
	t := &Typer{
		com: c,
		mr:  make(map[*Symbol][]*Symbol),
	}
	t.gb.init()
	t.deps.init()
	return t
}

func (t *Typer) reset() {
	t.box.reset()

	t.warns = t.warns[:0]
	t.errors = t.errors[:0]

	t.funs = t.funs[:0]
	t.vars = t.vars[:0]
	t.types = t.types[:0]
	t.consts = t.consts[:0]
	t.methods = t.methods[:0]

	clear(t.mr)

	t.state = StateAlloc
}

// Translate perform full unit translation.
func (t *Typer) Translate(texts []*ast.Text) diag.Error {
	if t.state != StateAlloc {
		return t.resume()
	}

	t.Alloc(texts)
	t.init()

	err := t.Index()
	if err != nil {
		return err
	}

	err = t.Scan()
	if err != nil {
		return err
	}

	return nil
}

func (t *Typer) init() {
	if len(t.box.texts) == 0 {
		panic("no texts")
	}

	t.unit.Init(&t.com.Global)
}

func (t *Typer) resume() diag.Error {
	return nil
}

// Alloc is a more efficient way (compared to Add) to add multiple Texts.
func (t *Typer) Alloc(texts []*ast.Text) {
	t.box.alloc(texts)
}

// Add is used by client to populate unit with source AST before translation.
func (t *Typer) Add(text *ast.Text) {
	t.box.addText(text)
}

// Index and create unit-level symbols.
func (t *Typer) Index() diag.Error {
	t.state = StateIndex

	t.indexImports()
	t.indexConsts()
	t.indexTypes()
	t.indexVars()
	t.indexFuns()
	// t.index
	t.indexMethods()

	if len(t.errors) != 0 {
		// need to refactor this to enable reporting multiple errors at once
		return t.errors[0]
	}
	return nil
}

func (t *Typer) Scan() diag.Error {
	t.state = StateScan

	num := len(t.consts) + len(t.types)
	if num != 0 {
		t.gb.reset(num)
		t.deps.discard = false

		t.scanConstSymbols()
		t.scanTypeSymbols()

		t.gb.Build()
	}

	if len(t.vars) != 0 {
		t.deps.discard = true
		t.scanVarSymbols()
	}

	if len(t.errors) != 0 {
		// need to refactor this to enable reporting multiple errors at once
		return t.errors[0]
	}
	return nil
}
