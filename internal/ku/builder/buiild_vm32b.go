package builder

import (
	"fmt"
	"os"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/genb"
	"github.com/mebyus/ku/internal/ku/parser"
	"github.com/mebyus/ku/internal/ku/stg"
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/vm32b"
)

func CompileFileForVM32B(path string) error {
	pool := sx.New()
	tf, err := pool.Load(path)
	if err != nil {
		return err
	}

	text := parser.ParseText(tf)

	n := 0 // total number of errors
	n += len(text.Errors)

	for _, e := range text.Errors {
		pos := pool.DecodePin(e.Pin)
		fmt.Fprintf(os.Stderr, "%s: %s\n", pos, e.Short)
	}

	var prog stg.Program
	typers := stg.NewPool(pool)
	typer := typers.Get()
	unit := stg.Unit{Name: "ku"}
	typer.Do(&unit, []*ast.Text{text})
	typers.Put(typer)
	prog.Units = append(prog.Units, &unit)

	n += len(unit.Errors)
	for _, e := range unit.Errors {
		pos := pool.DecodePin(e.Pin)
		fmt.Fprintf(os.Stderr, "%s: %s\n", pos, e.Short)
	}

	if n != 0 {
		os.Exit(1)
	}

	// TODO: should we pass Common to pool by pointer from Program instead?
	prog.Common = typers.Common

	funs := genb.Gen(&prog)
	p := vm32b.Compile(funs)
	vm32b.PrintProg(p)

	var vm vm32b.VM
	vm.Exec(p)

	fmt.Println()
	vm.PrintState()

	return nil
}
