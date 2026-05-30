package ast

// TypeSpec node that represents type specifier of any kind.
type TypeSpec interface {
	_spec()
}

// Embed this to quickly implement _spec() discriminator from TypeSpec interface.
// Do not use it for anything else.
type spec struct{}

func (spec) _spec() {}

// Explicit interface implementation check.
var _ TypeSpec = spec{}
