package parser

import (
	"fmt"
	"testing"

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

	if len(prog.Statements) != expectedStatements {
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
		integerValue interface{}
	}{
		{"!15", "!", 15},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
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

		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
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
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},

		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a + b / c", "(a + (b / c))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},

		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},

		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},

		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, tt := range tests {
		actual := testParserSetup(t, tt.input, -1).String()
		if actual != tt.expected {
			t.Errorf("expected=%q, found=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)

	return false
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	prog := testParserSetup(t, input, 1)

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prog.Statement[0] is not *ast.ExpressionStatement. got=%s",
			prog.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("consequence is not 1 statements. got=%T",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements is not nil. got=%+v",
			exp.Alternative)
	}
}

func TestFunctionExpression(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	prog := testParserSetup(t, input, 1)

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("prog.Statements[0] is not ast.ExpressionStatement. got=%T",
			prog.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionExpression. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function paramters wrong. expected=2 got=%d",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statement expected 1 statements. got=%d",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body statement is not *ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	prog := testParserSetup(t, input, 1)

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("prog.Statements[0] is not ast.ExpressionStatement. got=%T",
			prog.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length or arguments. got=%d",
			len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"Hello world";`

	prog := testParserSetup(t, input, 1)
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	expected := "Hello world"
	if literal.Value != expected {
		t.Errorf("literal.Value not %q. got=%q", expected, literal.Value)
	}
}

func TestListLiterals(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	prog := testParserSetup(t, input, -1)
	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	list, ok := stmt.Expression.(*ast.ListLiteral)
	if !ok {
		t.Fatalf("exp not *ast.ListLiteral. got=%T", stmt.Expression)
	}

	if len(list.Elements) != 3 {
		t.Fatalf("len(list.Elements) not 3, got=%d", len(list.Elements))
	}

	testIntegerLiteral(t, list.Elements[0], 1)
	testInfixExpression(t, list.Elements[1], 2, "*", 2)
	testInfixExpression(t, list.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpression(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	prog := testParserSetup(t, input, -1)
	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	list, ok := stmt.Expression.(*ast.ListLiteral)
	if !ok {
		t.Fatalf("exp not *ast.ListLiteral. got=%T", stmt.Expression)
	}

	if len(list.Elements) != 3 {
		t.Fatalf("len(list.Elements) not 3, got=%d", len(list.Elements))
	}

	testIntegerLiteral(t, list.Elements[0], 1)
	testInfixExpression(t, list.Elements[1], 2, "*", 2)
	testInfixExpression(t, list.Elements[2], 3, "+", 3)
}
