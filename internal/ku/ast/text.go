package ast

// Status describes result of parsing source text.
type Status uint32

const (
	// Ok is the only error-free result of parsing.
	//
	// Means that whole text was parsed successfully.
	Ok Status = iota

	// Flawed means that text was parsed to the end, but
	// contains parsing errors.
	Flawed

	// Parsing stopped because too many errors occured.
	ErrorLimitReached

	// Parsing stopped because sync after error failed.
	ErrorSyncFailed

	// Parsing stopped because text contains too many illegal tokens.
	ErrorIllegalTokens
)

// Text represents result of parsing a single source text.
//
// Contains all top-level nodes organized into lists of each type.
type Text struct {
	Funs []Fun

	Stubs []FunStub

	// List of all (not only top-level) errors occured during parsing.
	Errors []*Error

	Status Status
}

func (t *Text) IsOk() bool {
	return t.Status == Ok
}

func (t *Text) AddFun(f Fun) {
	t.Funs = append(t.Funs, f)
}

func (t *Text) AddStub(s FunStub) {
	t.Stubs = append(t.Stubs, s)
}
