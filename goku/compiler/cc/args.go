package cc

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/enums/bk"
)

type ArgsBuilder struct {
	args []string
}

func (g *ArgsBuilder) Init() {
	const initialArgsCapacity = 32
	g.args = make([]string, 0, initialArgsCapacity)
}

func (g *ArgsBuilder) Args() []string {
	return g.args
}

func (g *ArgsBuilder) add(arg ...string) {
	g.args = append(g.args, arg...)
}

func (g *ArgsBuilder) common() {
	g.add(codegenFlags...)
	g.add(maxErrorsFlag)
	g.add(warningFlags...)
	g.add(otherFlags...)
}

func (g *ArgsBuilder) out(path string) {
	if path == "" {
		panic("empty output path")
	}
	g.add("-o", path)
}

func (g *ArgsBuilder) src(path string) {
	if path == "" {
		panic("empty source path")
	}
	g.add(path)
}

func (g *ArgsBuilder) srcObj(path string) {
	g.add("-c")
	g.src(path)
}

func (g *ArgsBuilder) srcExe(path string) {
	g.src(path)
}

func (g *ArgsBuilder) optimizations(o string) {
	g.add("-O" + o)
}

func (g *ArgsBuilder) optimizationsAndDebugInfo(k bk.Kind) {
	switch k {
	case 0:
		panic("unspecified build kind")
	case bk.Debug:
		g.optimizations(debugCompilerOptimizations)
		g.add(debugInfoFlag)
	case bk.Test:
		g.optimizations(testCompilerOptimizations)
		g.add(debugInfoFlag)
	case bk.Safe:
		g.optimizations(safeCompilerOptimizations)
	case bk.Fast:
		g.optimizations(fastCompilerOptimizations)
	default:
		panic(fmt.Sprintf("unexpected build kind \"%s\" (=%d)", k, k))
	}
}
