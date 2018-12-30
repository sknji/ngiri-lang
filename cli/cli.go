// reads input, sends it to the interpreter for evaluation, prints the result and starts again
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/marmotini/monkey-lang/compiler"
	"github.com/marmotini/monkey-lang/interpreter"
	"github.com/marmotini/monkey-lang/lexer"
	"github.com/marmotini/monkey-lang/object"
	"github.com/marmotini/monkey-lang/parser"
	"github.com/marmotini/monkey-lang/vm"
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

		var evaluated object.Object
		if runVm {
			var err error
			evaluated, err = executeVM(p, os.Stdout)
			if err != nil {
				fmt.Fprintf(os.Stdout, err.Error())
			}
		} else {
			env := object.NewEnvironment()
			evaluated = interpreter.Eval(p.ParseProgram(), env)
		}

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

		var evaluated object.Object
		if runVm {
			var err error
			evaluated, err = executeVM(p, os.Stdout)
			if err != nil {
				fmt.Fprintf(os.Stdout, err.Error())
			}
		} else {
			evaluated = interpreter.Eval(p.ParseProgram(), env)
		}

		if evaluated != nil {
			io.WriteString(w, evaluated.Inspect())
			io.WriteString(w, "\n")
		}
	}
}

func executeVM(p *parser.Parser, w io.Writer) (object.Object, error) {
	comp := compiler.NewCompiler()
	err := comp.Compile(p.ParseProgram())
	if err != nil {
		return nil, fmt.Errorf("Woops! Compilation failed:\n %s\n", err)
	}

	machine := vm.NewVM(comp.Bytecode())
	err = machine.Run()
	if err != nil {
		return nil, fmt.Errorf("Woops! Executing bytecode failed:\n %s\n", err)
	}

	return machine.LastPoppedStackElem(), nil
}

func printParseErrors(w io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(w, err)
	}
}
