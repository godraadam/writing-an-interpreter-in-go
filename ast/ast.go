package ast

import (
	"bytes"
	"monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type Program struct {
	Stmts []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Stmts) >= 1 {
		return p.Stmts[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Stmts {
		out.WriteString(s.String())
	}

	return out.String()
}

// Let statement
type LetStmt struct {
	Token token.Token
	Name  *Identifier
	Value Expr
}

func (l *LetStmt) stmtNode() {}
func (l *LetStmt) TokenLiteral() string {
	return l.Token.Literal
}
func (l *LetStmt) String() string {
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")

	if l != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Return statement
type ReturnStmt struct {
	Token token.Token
	Value Expr
}

func (l *ReturnStmt) stmtNode() {}
func (l *ReturnStmt) TokenLiteral() string {
	return l.Token.Literal
}

func (l *ReturnStmt) String() string {
	var out bytes.Buffer
	out.WriteString(l.TokenLiteral() + " ")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Expression statement
type ExprStmt struct {
	Token      token.Token
	Expression Expr
}

func (l *ExprStmt) stmtNode() {}
func (l *ExprStmt) TokenLiteral() string {
	return l.Token.Literal
}

func (l *ExprStmt) String() string {
	if l != nil {
		return l.Expression.String()
	}
	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) exprNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}
