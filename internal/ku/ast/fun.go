package ast

import (
	"github.com/mebyus/ku/internal/fbuf"
	"github.com/mebyus/ku/internal/ku/sx"
)

// Fun represents unit-level named function definition (signature + body).
//
// Formal definition:
//
//	Fun  -> "fun" Name Signature Body
//	Body -> Block
//	Name -> word
type Fun struct {
	Sig  Signature
	Body Block

	Charms []Charm

	Name string

	// Function name pin.
	Pin sx.Pin

	// True when function is marked as public.
	Pub bool
}

func (f *Fun) Format(buf *fbuf.Buffer) {
	for _, c := range f.Charms {
		c.Format(buf)
		buf.Line()
	}
	if f.Pub {
		buf.Puts("pub\n")
	}
	buf.Puts("fun ")
	buf.Puts(f.Name)
	f.Sig.Format(buf)
	buf.Space()
	f.Body.Format(buf)
}

func (f *Fun) Stub() FunStub {
	return FunStub{
		Sig:    f.Sig,
		Charms: f.Charms,
		Name:   f.Name,
		Pin:    f.Pin,
		Pub:    f.Pub,
	}
}

// FunStub represents function declaration with no implementation.
//
// Formal definition:
//
//	Fun  -> "fun" Name Signature ";"
type FunStub struct {
	Sig Signature

	Charms []Charm

	Name string

	// Function name pin.
	Pin sx.Pin

	// True when function is marked as public.
	Pub bool
}

// Signature represents function signature (params + return result).
//
// Formal definition:
//
//	Signature -> "(" ParamList ")" [ "->" ( Result | "never" ) ]
//	ParamList -> { Param "," } // last comma is optional
//	Result    -> TypeSpec
type Signature struct {
	// Equals nil if there are no parameters in signature.
	Params []Param

	// Equals nil if function returns nothing or never returns.
	Result TypeSpec

	// Equals true if function never returns.
	Never bool
}

func (s *Signature) Format(buf *fbuf.Buffer) {
	if len(s.Params) == 0 {
		buf.Puts("()")
	} else {
		buf.Putb('(')
		s.Params[0].Format(buf)
		for _, p := range s.Params[1:] {
			buf.Puts(", ")
			p.Format(buf)
		}
		buf.Putb(')')
	}

	if s.Never {
		buf.Puts(" -> never")
		return
	}
	if s.Result == nil {
		return
	}

	buf.Puts(" -> ")
	s.Result.Format(buf)
}

// Param represents a single parameter in function signature.
//
// Formal definition:
//
//	Param -> Name ":" TypeSpec
//	Name  -> word
type Param struct {
	Type TypeSpec
	Name string

	// Parameter name pin.
	Pin sx.Pin
}

func (p *Param) Format(buf *fbuf.Buffer) {
	buf.Puts(p.Name)
	buf.Puts(": ")
	p.Type.Format(buf)
}
