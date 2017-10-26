package ast

import (
	"github.com/marmotini/monkey-lang/token"
	"testing"
)

func TestString(t *testing.T) {
	prog := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if str := prog.String(); str != "let myVar = anotherVar;" {
		t.Errorf("Program.String() wrong. got='%s'", str)
	}
}
