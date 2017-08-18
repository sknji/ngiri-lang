package parser

import (
	"testing"

	"github.com/nerdysquirrel/monkey-lang/lexer"
	"github.com/nerdysquirrel/monkey-lang/ast"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5; let y = 10; let foobar = 838383;`

	p := NewParser(lexer.NewLexer(input))

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if program == nil {
		t.Errorf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Errorf("program.Statements does not contain 3 statements, got=%d",
			len(program.Statements))
	}

	tests := []struct{
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, v := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, v.expectedIdentifier) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors{
		t.Errorf("parser error: %q", msg)
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%q", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("stmt.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	input := `return 5; return 10; return 993322;`

	p := NewParser(lexer.NewLexer(input))
	prog := p.ParseProgram()

	if len(prog.Statements) != 3 {
		t.Errorf("program.Statements does not contain 3 statements. got=%q",
			len(prog.Statements))
		return
	}

	for _, stmt := range prog.Statements {
		rtn, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%q",
			stmt)
			return
		}

		if rtn.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral not `return`. got=%q",
			rtn.TokenLiteral())
		}
	}
}