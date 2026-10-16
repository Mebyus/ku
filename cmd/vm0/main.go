package main

import (
	"fmt"
	"time"
)

func main() {
	funs := []*Function{
		fun(0, 0, []stm{
			prt("Hello, world!\n"),
			asn(6, sub(10, 20)),
			asn(1, mul(add(2, mul(3, 0)), 11)),
			asn(0, add(-5, 2)),
			asn(3, add(loc(0), loc(1))),
			asn(2, mul(add(1, loc(3)), -7)),
			asn(5, eql(1, 1)),
			call(1, 2),
			asn(4, call(1, 127)),
			ift(val(1), []stm{
				asn(8, 777),
			}),
			while(less(3, 2), []stm{
				asn(0, add(loc(0), 1)),
			}),
			asn(7, call(2, 99, 6)),
			ret(eql(loc(2), loc(3))),
		}),
		fun(1, 1, []stm{
			prt("called 1\n"),
			ret(add(loc(0), 1)),
		}),
		fun(2, 1, []stm{
			asn(2, 5),
			asn(4, add(13, 2)),
			ret(sub(loc(0), loc(1))),
		}),
	}

	var c compiler
	c.init()
	c.compile(funs)
	prog := &c.prog
	printProgram(prog)

	// TODO: compilation can also output needed local
	// memory size in advance
	var vm vm

	vm.exec(prog)
	fmt.Println()
	fmt.Printf("evs:  %v\n", vm.evs)
	fmt.Printf("smem: %v\n", vm.smem)

	const bench = false

	if bench {
		const n = 10000000
		start := time.Now()
		for range n {
			vm.exec(prog)
		}
		elapsed := time.Since(start)
		fmt.Printf("elapsed: %v\n", elapsed)
		fmt.Printf("pass:    %v\n", elapsed/n)
		fmt.Printf("per 10:  %v\n", elapsed*10/n/time.Duration(vm.total))
		fmt.Printf("itotal:  %d\n", vm.total)
	}
}
