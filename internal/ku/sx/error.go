package sx

import (
	"fmt"
	"io"
)

// ErrorType discriminator that uniquely identifies what compilation
// phase produced en error, by what subsystem/module and what syntax,
// semantic or typecheck was violated and how significant the error is.
//
// Zero value means that it is unstable error, with no assigned type
// yet. Treat it like todo for compiler development.
//
// Type is stored as unsigned integer which encodes several things
// about the error. Encoding scheme:
//
// Two most significant bytes encode violation level and compilation
// phase/compiler sybsystem/module that reported the error
//
// Violation level is stored in 4 most significant bits. It determines
// how errors should be treated: light warnings, warnings that become
// errors in production builds, errors which block codegeneration, etc.
type ErrorType uint64

func (t ErrorType) Level() Viol {
	return Viol(t >> 60)
}

// Viol specifies error violation level. Lower values mean lighter errors.
type Viol uint32

const (
	// Style means warning that does not affect typechecks and
	// codegeneration in any bad way. Programs with these warnings
	// are fine to compile and run, but can be improved/modernized
	// in semantic way.
	Style Viol = iota

	// Light means light warning that could indicate something suspicious
	// in program source code but it does not warrant build fail in
	// production mode.
	Light

	// Warn means hard warning. Program will compile and run in debug/test
	// mode, but build will fail in production mode until warning is
	// resolved or explicitly suppressed via compiler flags.
	Warn

	// Break means that program is broken in some way, either in syntax,
	// semantic, typecheck rules, etc. It cannot be compiled even with
	// compiler flags for suppressing errors. But compiler can still
	// go further with compilation up until backend phases.
	Break

	// Fatal means that there is no sense trying to going further
	// during compilation. Error is so significant that it breaks
	// current and all subsequent phases of compilation.
	//
	// Compiler should report this error and exit immediately.
	Fatal Viol = 0xF
)

// Phase specifies compilation phase for reporting errors.
type Phase uint32

const (
	Lex Phase = iota
	Parse

	// Unit walk/discovery/import graph verification
	Walk

	// STG + typecheck
	Typc
)

func (t ErrorType) String() string {
	p := Phase((t >> 54) & 0b111111)

	var phase string
	switch p {
	case Lex:
		phase = "L"
	case Parse:
		phase = "P"
	case Walk:
		phase = "W"
	case Typc:
		phase = "T"
	default:
		panic(fmt.Sprintf("invalid phase (=%d)", p))
	}

	x := uint64((t << 16) >> 16) // set two highest bytes to zero
	return fmt.Sprintf("%s.%04d", phase, x)
}

// Error describes compiler error with source position.
type Error struct {
	// Possibly long error description and additional info
	// with meaning varying based on type of error.
	//
	// May provide user help, suggestions, advice on how to fix.
	Notes []string

	// Short error message that describes what is wrong.
	Short string

	Type ErrorType

	// What place (in program source code) error should be attributed.
	//
	// Equals 0 if such place cannot be specified.
	Pin Pin
}

// Report build/compilation error with optional notes helpful to user.
//
// When pin=0 it is treated as error that is not attached
// to any particular source location.
func (p *Pool) Report(pin Pin, typ ErrorType, short string, notes ...string) {
	v := typ.Level()
	if v > p.vmax {
		p.vmax = v
	}

	p.ers = append(p.ers, Error{
		Pin:   pin,
		Type:  typ,
		Short: short,
		Notes: notes,
	})
}

func (p *Pool) HasErrors() bool {
	return len(p.ers) != 0
}

// NumErrors returns number of collected errors.
func (p *Pool) NumErrors() int {
	return len(p.ers)
}

// DumpErrors format all collected errors into supplied output writer.
func (p *Pool) DumpErrors(out io.Writer) {
	for _, e := range p.ers {
		p.FormatError(out, &e)
	}
}

func (p *Pool) FormatError(out io.Writer, e *Error) {
	loc := p.FormatPin(e.Pin)
	if loc == "" {
		if e.Type == 0 {
			fmt.Fprintf(out, "%s\n", e.Short)
		} else {
			fmt.Fprintf(out, "[%s] %s\n", e.Type, e.Short)
		}
	} else {
		if e.Type == 0 {
			fmt.Fprintf(out, "%s: %s\n", loc, e.Short)
		} else {
			fmt.Fprintf(out, "%s: [%s] %s\n", loc, e.Type, e.Short)
		}
	}

	if len(e.Notes) == 0 {
		return
	}

	io.WriteString(out, "======\n\n")
	for _, note := range e.Notes {
		io.WriteString(out, note)
		io.WriteString(out, "\n\n")
	}
	io.WriteString(out, "======\n")
}

const (
	// break errors in stg/typechecking
	prefixBreakTypc = ErrorType(Break)<<60 | ErrorType(Typc)<<54

	// Definition/declaration of multiple symbols with the same name
	// in the same scope.
	DupSymbol ErrorType = prefixBreakTypc | iota

	UnknownSymbol
)
