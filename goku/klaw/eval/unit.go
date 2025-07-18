package eval

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/klaw/ast"
)

// Unit represents result of evaluating unit build script.
type Unit struct {
	Imports  []string
	Includes []string
}

func (u *Unit) valid() diag.Error {
	if !checkUnique(u.Imports) {
		return &diag.PinlessError{Text: fmt.Sprintf("non-unique imports %v", u.Imports)}
	}
	if !checkUnique(u.Includes) {
		return &diag.PinlessError{Text: fmt.Sprintf("non-unique includes %v", u.Includes)}
	}
	return nil
}

func EvalUnit(env *Env, unit *ast.Unit) (*Unit, diag.Error) {
	r := Interpreter{env: env}
	err := r.eval(unit.Dirs)
	if err != nil {
		return nil, err
	}
	err = r.unit.valid()
	if err != nil {
		return nil, err
	}
	return &r.unit, nil
}

type Interpreter struct {
	// Keeps track of result object.
	unit Unit

	env *Env
}

func (r *Interpreter) eval(dirs []ast.Dir) diag.Error {
	for _, d := range dirs {
		err := r.dir(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Interpreter) dir(dir ast.Dir) diag.Error {
	switch d := dir.(type) {
	case ast.Import:
		if d.Val == "" {
			panic("empty import path")
		}
		r.unit.Imports = append(r.unit.Imports, d.Val)
	case ast.Include:
		if d.Val == "" {
			panic("empty include path")
		}
		r.unit.Includes = append(r.unit.Includes, d.Val)
	case ast.Test:
		if !r.env.TestExe {
			return nil
		}
		return r.eval(d.Dirs)
	case ast.Exe:
		if !r.env.Exe {
			return nil
		}
		return r.eval(d.Dirs)
	default:
		panic(fmt.Sprintf("unexpected node %T", d))
	}
	return nil
}

// returns true if all strings in slice are unique
func checkUnique(ss []string) bool {
	if len(ss) < 2 {
		return true
	}
	if len(ss) == 2 {
		return ss[0] != ss[1]
	}

	set := make(map[string]struct{}, len(ss))
	for _, s := range ss {
		_, ok := set[s]
		if ok {
			return false
		}
		set[s] = struct{}{}
	}
	return true
}
