package vm32b

import "github.com/mebyus/ku/internal/vm32b/ir"

func Compile(funs []ir.Fun) *Prog {
	var c compiler
	c.compile(funs)
	return &c.prog
}

type compiler struct {
	prog Prog
}

func (c *compiler) compile(funs []ir.Fun) {

}
