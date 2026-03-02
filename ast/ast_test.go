package ast

import (
	"monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Stmts: []Stmt{
			&LetStmt{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Target: &Identifier{
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
	expected := "let myVar = anotherVar;"
	result := program.String()
	if result != expected {
		t.Errorf("Incorrect program.String() Expected %q, got %q", expected, result)
	}
}
