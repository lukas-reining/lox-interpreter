package lox

import (
	"bufio"
	"fmt"
	"github.com/lukas-reining/lox/interpeter"
	"github.com/lukas-reining/lox/parser"
	"github.com/lukas-reining/lox/scanner"
	"log"
	"os"
)

type Lox struct {
}

func NewLox() *Lox {
	return &Lox{}
}

func (l *Lox) report(line int, where string, messsage string) {
	_, _ = fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, messsage)
}

func (l *Lox) scannerError(err scanner.ScannerError) {
	l.report(err.Line(), "", err.Message())
}

func (l *Lox) parserError(err parser.ParseError) {
	token := err.Lexeme()
	if len(token) == 0 {
		l.report(err.Line(), " at end", err.Message())
	} else {
		l.report(err.Line(), " at '"+token+"'", err.Message())
	}
}

func (l *Lox) runtimeError(err interpeter.RuntimeError) {
	l.report(err.Line(), "", err.Message())
}

func (l *Lox) error(err error) {
	switch e := err.(type) {
	case scanner.ScannerError:
		l.scannerError(e)
	case parser.ParseError:
		l.parserError(e)
	case interpeter.RuntimeError:
		l.runtimeError(e)
	default:
		l.report(-1, "", e.Error())
	}
}

func (l *Lox) RunFile(filePath string) error {
	dat, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatalf("File not found: %s", filePath)
	}

	script := string(dat)
	_, _, err = l.Run(script, nil)

	if err != nil {
		l.error(err)
	}

	return err
}

func (l *Lox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	var currentEnv *interpeter.Environment
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')

		if err == nil {
			value, env, err := l.Run(text, currentEnv)

			if err != nil {
				l.error(err)
			} else {
				currentEnv = env
				println(interpeter.Stringify(value))
			}
		}
	}
}

func (l *Lox) Run(script string, env *interpeter.Environment) (interpeter.LoxValue, *interpeter.Environment, error) {
	sourceScanner := scanner.NewScanner(script)
	tokens, err := sourceScanner.ScanTokens()
	if err != nil {
		return nil, nil, err
	}

	sourceParser := parser.NewParser[any, interpeter.RuntimeError](tokens)
	statements, err := sourceParser.Parse()
	if err != nil {
		return nil, nil, err
	}

	interpreter := interpeter.NewInterpreterWithEnv(env)

	resolver := interpeter.NewResolver(&interpreter)
	if err := resolver.Resolve(statements); err != nil {
		return nil, nil, err
	}

	return interpreter.Interpret(statements)
}
