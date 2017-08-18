// takes input data(tokens), builds a data structure(AST) giving a structural representation
// of the input, checking for correct syntax in the process.
package parser

import (
	"fmt"

	"github.com/nerdysquirrel/monkey-lang/lexer"
	"github.com/nerdysquirrel/monkey-lang/token"
	"github.com/nerdysquirrel/monkey-lang/ast"
)

type Parser struct {
	l *lexer.Lexer
	
	currToken token.Token
	peekToken token.Token

	errors []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
	t, p.peekToken.Type)

	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken()  {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(type_ token.TokenType) bool {
	return p.currToken.Type == type_
}

func (p *Parser) peekTokenIs(type_ token.TokenType) bool {
	return p.peekToken.Type == type_
}

func (p *Parser) expectPeek(type_ token.TokenType) bool {
	if p.peekTokenIs(type_) {
		p.nextToken()
		return true
	}

	p.peekError(type_)

	return false
}

func (p *Parser) ParseProgram() (prog *ast.Program) {
	prog = &ast.Program{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}

		p.nextToken()
	}

	return
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
	return nil
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO
	for !p.expectPeek(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	ret := &ast.ReturnStatement{Token: p.currToken}

	for !p.expectPeek(token.SEMICOLON) {
		p.nextToken()
	}

	return ret
}