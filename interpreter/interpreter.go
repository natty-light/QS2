package interpreter

import (
	"QuonkScript/evaluator"
	"QuonkScript/lexer"
	"QuonkScript/object"
	"QuonkScript/parser"
	"fmt"
	"os"
)

func Run(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Honk! Cannot read file %s\n", filename)
		return
	}

	src := string(bytes)

	l := lexer.New(src)
	p := parser.New(l)
	scope := object.NewScope()

	program := p.ParseProgram()

	result := evaluator.Eval(program, scope)
	fmt.Println(result)
}
