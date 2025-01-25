package ast

// Build represents build block at the top of source text.
//
// Formal definition:
//
//	Build => "#build" Block
type Build struct {
	Body Block
}
