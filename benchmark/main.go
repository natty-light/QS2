package main

import (
	"flag"
	"fmt"
	"quonk/compiler"
	"quonk/evaluator"
	"quonk/lexer"
	"quonk/object"
	"quonk/parser"
	"quonk/vm"
	"time"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")

var source = `
const fib = func(x) {
	if (x == 0) {
		return 0;
	} else {
		if (x == 1) {
			return 1;
		} else {
			fib(x - 1) + fib(x - 2);
		}
	}
}
fib(35);
`

func main() {
	flag.Parse()
	var duration time.Duration
	var result object.Object

	l := lexer.New(source)
	p := parser.New(l)
	prog := p.ParseProgram()

	if *engine == "vm" {
		comp := compiler.New()
		err := comp.Compile(prog)
		if err != nil {
			fmt.Printf("compiler error: %s", err)
			return
		}

		machine := vm.New(comp.Bytecode())

		start := time.Now()

		err = machine.Run()
		if err != nil {
			fmt.Printf("vm error: %s", err)
			return
		}

		duration = time.Since(start)
		result = machine.LastPoppedStackElem()
	} else {
		env := object.NewScope()
		start := time.Now()
		result = evaluator.Eval(prog, env)
		duration = time.Since(start)
	}

	fmt.Printf("engine=%s, result=%s. duration=%s", *engine, result.Inspect(), duration)
}
