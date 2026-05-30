package stg

import (
	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Pool is a collection of Typer objects. Client code can request available
// object from Pool or return an object that is no longer needed. Pool automatically
// manages objects allocation, reuse and resetting state.
type Pool struct {
	Common

	// Objects available for reuse. Operates as stack.
	free []*Typer

	// Maps unit to a Typer attached to it.
	// Contains only occupied objects. Typer is moved to free stack when
	// removed from this map.
	// m map[*Unit]*Typer
}

func NewPool(pool *sx.Pool) *Pool {
	p := &Pool{}
	p.init(pool)
	return p
}

func (p *Pool) init(pool *sx.Pool) {
	p.Common.Init(pool)
}

// Get returns a Typer object attached to a given unit.
// Allocates a new object if specified unit has no Typer attached to it.
func (p *Pool) Get() *Typer {
	// t := p.m[u]
	// if t != nil {
	// 	return t
	// }

	var t *Typer
	n := len(p.free)
	if n == 0 {
		t = NewTyper(&p.Common)
	} else {
		t = p.free[n-1]
		p.free = p.free[:n-1]
		t.reset()
	}

	u := &Unit{}
	u.init(&p.Global)
	t.unit = u
	// p.m[u] = t
	return t
}

func (p *Pool) Put(t *Typer) {
	// delete(p.m, t.unit)
	t.unit = nil
	p.free = append(p.free, t)
}

type Typer struct {
	box NodeBox

	common *Common

	unit *Unit

	// Signature of current function or method being converted to STG form.
	sig *Signature
}

func NewTyper(c *Common) *Typer {
	return &Typer{common: c}
}

func (t *Typer) reset() {
	t.unit = nil
	t.box.reset()
}

// Translate combines parsed source texts into a single unit and
// translates code into STG form.
//
// Assigns types and does typechecking for all symbols and expressions inside
// resulting unit.
func (t *Typer) Translate(texts []*ast.Text) *Unit {
	t.alloc(texts)
	t.index()
	t.toptype()
	t.convert()
	return t.unit
}
