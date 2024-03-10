package main

import (
	"fmt"
	"os"
	"quonk/interpreter"
	"quonk/repl"
)

func main() {

	args := os.Args

	if len(args) == 1 {
		// If no filename was passed as a command line argument, run the repl
		repl.Start(os.Stdin, os.Stdout)
	} else {
		if args[1] == "run" {
			interpreter.Run(args[2])
		} else if args[1] == "compile" {
			// TODO: implement compile
		} else if args[1] == "exec" {
			// TODO: implement reading intermediate bytecode file
		} else if args[1] == "help" {
			fmt.Println("Usage: quonk [run|compile|exec|help] [filename]")
		}
	}

}
