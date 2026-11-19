package stg

import "github.com/mebyus/ku/internal/ku/sx"

// SelectField represents expression which selects a field
// from expression of custom type (struct or union).
type SelectField struct {
	operand

	// Target for select application.
	// Always not nil.
	Exp Exp

	// Selected field.
	// Always not nil
	Field *Field

	// pin where select happens
	pin sx.Pin
}

var _ Operand = &SelectField{}

func (s *SelectField) Pin() sx.Pin {
	return s.pin
}

func (s *SelectField) Type() *Type {
	return s.Field.Type
}

func (s *SelectField) Target() bool {
	return true
}

// DelectField represents expression which dereferences expression and
// then selects a field from pointer to custom type (struct or union).
type DelectField struct {
	operand

	// Target for delect application.
	// Always not nil.
	Exp Exp

	// Selected field.
	// Always not nil
	Field *Field

	// pin where delect happens
	pin sx.Pin
}

var _ Operand = &SelectField{}

func (d *DelectField) Pin() sx.Pin {
	return d.pin
}

func (d *DelectField) Type() *Type {
	return d.Field.Type
}

func (d *DelectField) Target() bool {
	return true
}

// BoundMethod represents expression which is formed by binding
// a method with receiver expression.
//
// For example:
//
//	buf.push
//	file.write
type BoundMethod struct {
	operand

	// Bound receiver.
	Recv Exp

	// Which method was bound.
	Symbol *Symbol

	// pin where binding happens
	pin sx.Pin
}

func (m *BoundMethod) Pin() sx.Pin {
	return m.pin
}

func (m *BoundMethod) Type() *Type {
	panic("TODO")
}

func (m *BoundMethod) Target() bool {
	return false
}
