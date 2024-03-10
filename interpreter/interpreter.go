package interpreter

import (
	"fmt"
	"os"
	"quonk/evaluator"
	"quonk/lexer"
	"quonk/object"
	"quonk/parser"
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
