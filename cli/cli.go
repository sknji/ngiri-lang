// reads input, sends it to the interpreter for evaluation, prints the result and starts again
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nerdysquirrel/monkey-lang/lexer"
	"github.com/nerdysquirrel/monkey-lang/parser"
	"github.com/nerdysquirrel/monkey-lang/runtime"
	"github.com/nerdysquirrel/monkey-lang/object"
)

const PROMPT = ">> "

var (
	interactive bool
	fileName    string
)

func init() {
	flag.BoolVar(&interactive, "i", false, "interactive mode")
	flag.StringVar(&fileName, "f", "", "filename")
}

func main() {
	flag.Parse()

	if fileName != "" {
		p := parser.NewParser(lexer.NewLexerFromFile(fileName))
		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Printf("Parser error: %s\n", err)
			}
		}

		env := object.NewEnvironment()
		evaluated := runtime.Eval(p.ParseProgram(), env)
		if evaluated != nil {
			io.WriteString(os.Stdout, evaluated.Inspect())
			io.WriteString(os.Stdout, "\n")
		}
	}

	if interactive && fileName == "" {
		StartInteractiveMode(os.Stdin, os.Stdout)
	}
}

func StartInteractiveMode(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		p := parser.NewParser(lexer.NewLexer(line))
		if len(p.Errors()) > 0 {
			printParseErrors(w, p.Errors())
			continue
		}

		evaluated := runtime.Eval(p.ParseProgram(), env)
		if evaluated != nil {
			io.WriteString(w, evaluated.Inspect())
			io.WriteString(w, "\n")
		}
	}
}

func printParseErrors(w io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(w, err)
	}
}
