package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Assign represents assign statement where assign target can
// be any syntactically correct expression.
type Assign struct {
	stm

	// What is being changed by assignment.
	//
	// Always not nil.
	Target Exp

	// Represents value being assigned.
	//
	// Always being nil.
	Exp Exp
}

var _ Statement = &Assign{}

func (a *Assign) Format(buf *fbuf.Buffer) {
	a.Target.Format(buf)
	buf.Puts(" = ")
	a.Exp.Format(buf)
	buf.Semi()
}

// AssignSymbol represents statement that assign a new value to
// existing symbol.
//
// Formal definition:
//
//	AssignSymbol -> Name = Exp ";"
//	Name         -> word
type AssignSymbol struct {
	stm

	// Always not nil.
	Exp Exp

	// Name of symbol being assigned.
	Name string

	// Name pin.
	Pin sx.Pin
}

var _ Statement = &AssignSymbol{}

func (a *AssignSymbol) Format(buf *fbuf.Buffer) {
	buf.Puts(a.Name)
	buf.Puts(" = ")
	a.Exp.Format(buf)
	buf.Semi()
}

// AssignDiscard represents statement that discards result of
// evaluating expression.
//
// Formal definition:
//
//	AssignDiscard -> "_" = Exp ";"
type AssignDiscard struct {
	stm

	// Always not nil.
	Exp Exp

	// Discard pin.
	Pin sx.Pin
}

var _ Statement = &AssignDiscard{}

func (a *AssignDiscard) Format(buf *fbuf.Buffer) {
	buf.Puts("_ = ")
	a.Exp.Format(buf)
	buf.Semi()
}
