// takes input data(tokens), builds a data structure(AST) giving a structural representation
// of the input, checking for correct syntax in the process.
package parser

import (
	"github.com/nerdysquirrel/monkey-lang/lexer"
	"github.com/nerdysquirrel/monkey-lang/token"
	"github.com/nerdysquirrel/monkey-lang/ast"
)

type Parser struct {
	l *lexer.Lexer
	
	currToken token.Token
	peekToken token.Token
}

func NewParser() *Parser {
	p := &Parser{}
	
	p.nextToken()
	p.nextToken()
}

func (p *Parser) nextToken()  {
	
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}