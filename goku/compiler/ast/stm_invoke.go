package ast

// Invoke represents a call expression statement.
//
// Formal definition:
//
//	Invoke => Call ";"
type Invoke struct {
	Call Call
}

var _ Statement = Invoke{}
