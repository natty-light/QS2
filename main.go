package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"quonk/compiler"
	"quonk/evaluator"
	"quonk/lexer"
	"quonk/object"
	"quonk/parser"
	"quonk/repl"
	"quonk/vm"
)

func main() {

	args := os.Args

	if len(args) == 1 {
		// If no filename was passed as a command line argument, run the repl
		repl.Start(os.Stdin, os.Stdout)
	} else {
		if args[1] == "interpret" {
			Interpret(args[2])
		} else if args[1] == "run" {
			Run(args[2])
		} else if args[1] == "compile" {
			Compile(args[2])
		} else if args[1] == "exec" {
			// TODO: implement reading intermediate bytecode file
		} else if args[1] == "help" {
			fmt.Println("Usage: quonk [interpret|compile|exec|help] [filename]")
		}
	}

}

func Interpret(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Honk! Cannot read file %s\n", filename)
		return
	}

	src := string(file)

	l := lexer.New(src)
	p := parser.New(l)
	scope := object.NewScope()

	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Honk! Parser errors:")
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
	} else {
		result := evaluator.Eval(program, scope)
		if result != nil {
			fmt.Println(result.Inspect())
		}
	}

}

func Run(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Honk! Cannot read file %s\n", filename)
		return
	}

	src := string(file)

	l := lexer.New(src)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	comp := compiler.New()
	err = comp.Compile(program)
	if err != nil {
		fmt.Printf("Compiler error: %s\n", err)
		return
	}

	machine := vm.New(comp.Bytecode())
	err = machine.Run()
	if err != nil {
		fmt.Printf("Runtime error: %s\n", err)
		return
	}

	stackTop := machine.StackTop()
	fmt.Println(stackTop.Inspect())
}

func Compile(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Honk! Cannot read file %s\n", filename)
		return
	}

	src := string(file)

	l := lexer.New(src)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	comp := compiler.New()
	err = comp.Compile(program)
	if err != nil {
		fmt.Printf("Compiler error: %s\n", err)
		return
	}

	bytecode := comp.Bytecode()
	var out bytes.Buffer
	out.Write(bytecode.Instructions)
	out.Write([]byte("\n"))
	// TODO: come up with way of encoding constants as bytes
}

func Exec(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Honk! Cannot read file %s\n", filename)
		return
	}

	bytecode := parseFileIntoBytecode(file)
	machine := vm.New(bytecode)
	err = machine.Run()
	if err != nil {
		fmt.Printf("Runtime error: %s\n", err)
		return
	}

	stackTop := machine.StackTop()
	fmt.Println(stackTop.Inspect())
}

func parseFileIntoBytecode(file []byte) *compiler.Bytecode {
	return &compiler.Bytecode{}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
