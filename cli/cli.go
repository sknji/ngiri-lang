// reads input, sends it to the interpreter for evaluation, prints the result and starts again
package main

import (
	"io"
	"bufio"
	"fmt"
	"os"

	"github.com/nerdysquirrel/monkey-lang/lexer"
	"github.com/nerdysquirrel/monkey-lang/token"
)

const PROMPT  = ">> "

func main()  {
	Start(os.Stdin, os.Stdout)
}

func Start(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	for {
		fmt.Printf(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			continue
		}

		line := scanner.Text()
		lex := lexer.NewLexer(line)

		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}