package ast

import "github.com/mebyus/ku/internal/fbuf"

// Status describes result of parsing source text.
//
// Zero value means that whole text was parsed successfully.
type Status uint32

const (
	// Stop used as parser status to exit parsing after
	// reaching EOF. Parsed texts cannot have this status.
	Stop Status = 1 + iota

	// Flawed means that text was parsed to the end, but
	// contains parsing errors.
	Flawed

	// Parsing stopped because too many errors occured.
	LimitReached

	// Parsing stopped because sync after error failed.
	SyncFail

	// Parsing stopped because text contains too many illegal tokens.
	InvalidTokens
)

// Text represents result of parsing a single source text.
//
// Contains all top-level nodes organized into lists of each type.
type Text struct {
	Imports []Import

	Funs []Fun

	Stubs []FunStub

	Types []Type

	Consts []Const

	Status Status
}

func (t *Text) HasErrors() bool {
	return t.Status != 0
}

func (t *Text) AddImport(i Import) {
	t.Imports = append(t.Imports, i)
}

func (t *Text) AddFun(f Fun) {
	t.Funs = append(t.Funs, f)
}

func (t *Text) AddStub(s FunStub) {
	t.Stubs = append(t.Stubs, s)
}

func (t *Text) AddType(typ Type) {
	t.Types = append(t.Types, typ)
}

func (t *Text) AddConst(c Const) {
	t.Consts = append(t.Consts, c)
}

func (t *Text) Format(buf *fbuf.Buffer) {

}
