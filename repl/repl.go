package repl

import (
	"bufio"
	"fmt"
	"io"
	"quonk/compiler"
	"quonk/evaluator"
	"quonk/lexer"
	"quonk/object"
	"quonk/parser"
	"quonk/vm"
)

const PROMPT = ">>"

func StartEval(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	scope := object.NewScope()
	macroScope := object.NewScope()
	fmt.Fprint(out, "QuonkScript REPL v0.1\n")
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		evaluator.DefineMacros(program, macroScope)
		expanded := evaluator.ExpandMacros(program, macroScope)

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(expanded, scope)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func StartVM(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Honk! compiler error:\n %s\n", err)
			continue
		}

		machine := vm.NewWithGlobalStore(comp.Bytecode(), globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Honk! runtime error:\n %s\n", err)
			continue
		}

		stackTop := machine.LastPoppedStackElem()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
