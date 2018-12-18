// reads input, sends it to the interpreter for evaluation, prints the result and starts again
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/marmotini/monkey-lang/compiler"
	"github.com/marmotini/monkey-lang/vm"
	"io"
	"os"

	"github.com/marmotini/monkey-lang/interpreter"
	"github.com/marmotini/monkey-lang/lexer"
	"github.com/marmotini/monkey-lang/object"
	"github.com/marmotini/monkey-lang/parser"
)

const PROMPT = ">> "

var (
	interactive bool
	fileName    string
	runVm       bool
)

func init() {
	flag.BoolVar(&interactive, "i", false, "interactive mode")
	flag.StringVar(&fileName, "f", "", "filename")
	flag.BoolVar(&runVm, "vm", true, "run virtual machine")
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
		evaluated := interpreter.Eval(p.ParseProgram(), env)
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

		if runVm {
			comp := compiler.NewCompiler()
			err := comp.Compile(p)
			if err != nil {
				fmt.Printf(w, "Woops! Compilation failed:\n %s\n", err)
				continue
			}

			machine := vm.NewVM(comp.Bytecode())
			err = machine.Run()
			if err != nil {
				fmt.Fprintf(w, "Woops! Executing bytecode failed:\n %s\n", err)
				continue
			}

			lastPopped := machine.LastPoppedStackElem()
			io.WriteString(w, lastPopped.Inspect())
			io.WriteString(w, "\n")
		}

		evaluated := interpreter.Eval(p.ParseProgram(), env)
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
