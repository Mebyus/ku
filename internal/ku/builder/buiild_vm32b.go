package builder

import (
	"fmt"
	"io"
	"os"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/builder/kub"
	"github.com/mebyus/ku/internal/ku/genb"
	"github.com/mebyus/ku/internal/ku/parser"
	"github.com/mebyus/ku/internal/ku/stg"
	"github.com/mebyus/ku/internal/ku/sx"
	"github.com/mebyus/ku/internal/vm32b"
)

func CompileForVM32B(path string) (*vm32b.Prog, error) {
	pool := sx.New()
	tf, err := pool.Load(path)
	if err != nil {
		return nil, err
	}

	text := parser.ParseText(tf)

	n := 0 // total number of errors
	n += len(text.Errors)

	for _, e := range text.Errors {
		pos := pool.DecodePin(e.Pin)
		fmt.Fprintf(os.Stderr, "%s: %s\n", pos, e.Short)
	}

	var prog stg.Program
	typers := stg.NewPool(kub.Build{
		Target: kub.Target{
			Arch: kub.ArchVM32B,
			OS:   kub.Bare,
		},
	}, pool)
	typer := typers.Get()
	unit := stg.Unit{Name: "ku"}
	typer.Do(&unit, []*ast.Text{text})
	typers.Put(typer)
	prog.Units = append(prog.Units, &unit)

	if pool.HasErrors() {
		pool.DumpErrors(os.Stderr)
		os.Exit(1)
	}

	// TODO: should we pass Common to pool by pointer from Program instead?
	prog.Common = typers.Common

	funs := genb.Gen(&prog)
	return vm32b.Compile(funs), nil
}

// CompileEncodeForVM32B compile Ku source file specified by path and
// encode resulting bytecode program into output writer.
func CompileEncodeForVM32B(path string, out io.Writer) error {
	prog, err := CompileForVM32B(path)
	if err != nil {
		return err
	}

	return vm32b.Encode(out, prog)
}

func CompileFileForVM32B(path string, out string) error {
	prog, err := CompileForVM32B(path)
	if err != nil {
		return err
	}

	// vm32b.PrintProg(p)

	// var vm vm32b.VM
	// vm.Exec(p, nil)

	// fmt.Println()
	// vm.PrintState()

	if out == "" {
		return nil
	}

	return vm32b.Save(prog, out)
}
