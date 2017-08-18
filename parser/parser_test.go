package parser

import (
	"testing"

	"fmt"
	"github.com/nerdysquirrel/monkey-lang/ast"
	"github.com/nerdysquirrel/monkey-lang/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5; let y = 10; let foobar = 838383;`

	program := testParserSetup(t, input, 3)

	tests := []struct {
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
	if p == nil {
		t.Errorf("program is nil")
	}

	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
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
	prog := testParserSetup(t, input, 3)

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

func TestIdentifierExpression(t *testing.T) {
	input := `foobar`
	prog := testParserSetup(t, input, 1)

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] is not as.ExpressionStatement. got=%q",
			prog.Statements[0])
		return
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%q", stmt.Expression)
		return
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not %s. got=%s", "foobar", ident.Value)
		return
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}

}

func testParserSetup(t *testing.T, input string, expectedStatements int) *ast.Program {
	p := NewParser(lexer.NewLexer(input))
	prog := p.ParseProgram()
	checkParserErrors(t, p)

	if expectedStatements == -1 {
		return prog
	}

	if len(prog.Statements) != 1 {
		t.Fatalf("program does not have enough statements. Statements %+v. expected=%d, got=%d",
			prog.Statements, expectedStatements, len(prog.Statements))
	}

	return prog
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`
	prog := testParserSetup(t, input, 1)

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] is not as.ExpressionStatement. got=%q",
			prog.Statements[0])
		return
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp not *ast.IntegerLiteral. got=%q", stmt.Expression)
		return
	}

	if literal.Value != 5 {
		t.Errorf("Literal.Value is not %s. got=%s", "foobar", literal.Value)
		return
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("Literal.TokenLiteral not %s. got=%s", "foobar", literal.TokenLiteral())
	}

}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTesting := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!15", "!", 15},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTesting {
		prog := testParserSetup(t, tt.input, 1)

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				prog.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		prog := testParserSetup(t, tt.input, 1)

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				prog.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Errorf("exp is not *ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Left, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T)  {
	tests := []struct{
		input string
		expected string
	}{
		{"-a * b", "((-a) * b)",},
		{"!-a", "(!(-a))",},
		{"a + b + c", "((a + b) + c)",},
		{"a + b - c", "((a + b) - c)",},
		{"a + b / c", "(a + (b / c))",},
	}

	for _, tt := range tests {
		actual := testParserSetup(t, tt.input, -1).String()
		if actual != tt.expected {
			t.Errorf("expected=%q, found=%q", tt.expected, actual)
		}
	}
}
