package cc

import (
	"os"
	"os/exec"

	"github.com/mebyus/ku/goku/compiler/enums/bk"
)

const compiler = "cc"

func InvokeCompiler(args []string) error {
	cmd := exec.Command(compiler, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CompileObj(out, src string, k bk.Kind) error {
	var g ArgsBuilder
	g.Init()

	g.common()
	g.optimizationsAndDebugInfo(k)
	g.out(out)
	g.srcObj(src)

	return InvokeCompiler(g.Args())
}

func CompileExe(out, src string, k bk.Kind) error {
	var g ArgsBuilder
	g.Init()

	g.common()
	g.optimizationsAndDebugInfo(k)
	g.out(out)
	g.srcExe(src)

	return InvokeCompiler(g.Args())
}
