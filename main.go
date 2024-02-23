package main

import (
	"QuonkScript/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
