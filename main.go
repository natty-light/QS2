package main

import (
	"QuonkScript/interpreter"
	"QuonkScript/repl"
	"os"
)

func main() {

	args := os.Args

	if len(args) == 1 {
		// If no filename was passed as a command line argument, run the repl
		repl.Start(os.Stdin, os.Stdout)
	} else {
		// Script name should be second arg
		interpreter.Run(args[1])
	}

}
