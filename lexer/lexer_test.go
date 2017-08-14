package lexer

import (
	"testing"

	"github.com/nerdysquirrel/monkey-lang/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	l := NewLexer(input)

	for k, v := range tests {
		tok := l.NextToken()

		if tok.Type != v.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				k, v.expectedType, tok.Type)
		}

		if tok.Literal != v.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q",
				k, v.expectedLiteral, tok.Literal)
		}
	}
}
