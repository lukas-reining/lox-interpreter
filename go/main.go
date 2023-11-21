package main

import (
	"fmt"
	"github.com/lukas-reining/lox/interpeter"
	"github.com/lukas-reining/lox/lox"
	"github.com/lukas-reining/lox/parser"
	"os"
)

func main() {
	args := os.Args[1:]
	loxEngine := lox.NewLox()

	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		err := loxEngine.RunFile(args[0])

		switch err.(type) {
		case parser.ParseError:
			os.Exit(65)
		case interpeter.RuntimeError:
			os.Exit(70)
		}
	} else {
		loxEngine.RunPrompt()
	}
}
