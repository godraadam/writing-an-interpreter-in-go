package ast

import (
	"bytes"
	"monkey/token"
	"strings"
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

	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Assignment statement
type AssingmentStmt struct {
	Token token.Token
	Name  *Identifier
	Value Expr
}

func (as *AssingmentStmt) stmtNode() {}
func (as *AssingmentStmt) TokenLiteral() string {
	return as.Token.Literal
}
func (as *AssingmentStmt) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(" = ")

	if as.Value != nil {
		out.WriteString(as.Value.String())
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

type PrintStmt struct {
	Token token.Token
	Value Expr
}

func (ps *PrintStmt) stmtNode() {}
func (ps *PrintStmt) TokenLiteral() string {
	return ps.Token.Literal
}

func (ps *PrintStmt) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	if ps.Value != nil {
		out.WriteString(ps.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Expression statement
type ExprStmt struct {
	Token token.Token
	Expr  Expr
}

func (l *ExprStmt) stmtNode() {}
func (l *ExprStmt) TokenLiteral() string {
	return l.Token.Literal
}

func (l *ExprStmt) String() string {
	if l.Expr != nil {
		return l.Expr.String()
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

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (i *NumberLiteral) exprNode() {}
func (i *NumberLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *NumberLiteral) String() string {
	return i.Token.Literal
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (i *BooleanLiteral) exprNode() {}
func (i *BooleanLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *BooleanLiteral) String() string {
	return i.Token.Literal
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) exprNode() {}
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("\"")
	out.WriteString(sl.Value)
	out.WriteString("\"")
	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expr
}

func (al *ArrayLiteral) exprNode() {}
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type PrefixExpr struct {
	Token    token.Token
	Operator string
	Right    Expr
}

func (i *PrefixExpr) exprNode() {}
func (i *PrefixExpr) TokenLiteral() string {
	return i.Token.Literal
}

func (i *PrefixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Operator)
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpr struct {
	Token    token.Token
	Operator string
	Left     Expr
	Right    Expr
}

func (i *InfixExpr) exprNode() {}
func (i *InfixExpr) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" ")
	out.WriteString(i.Operator)
	out.WriteString(" ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

type IfExpr struct {
	Token       token.Token
	Condition   Expr
	Consequence *BlockStmt
	Alternative *BlockStmt
}

func (i *IfExpr) exprNode() {}
func (i *IfExpr) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IfExpr) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString(" ")
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}
	return out.String()
}

type BlockStmt struct {
	Token token.Token
	Stmts []Stmt
}

func (bs *BlockStmt) exprNode() {}
func (bs *BlockStmt) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStmt) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Stmts {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token  token.Token
	Params []*Identifier
	Body   *BlockStmt
}

func (fl *FunctionLiteral) exprNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Params {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpr struct {
	Token    token.Token
	Args     []Expr
	Function Expr
}

func (ce *CallExpr) exprNode() {}
func (ce *CallExpr) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpr) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range ce.Args {
		params = append(params, p.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	return out.String()
}

type IndexExpr struct {
	Token token.Token
	Left  Expr
	Index Expr
}

func (ie *IndexExpr) exprNode() {}
func (ie *IndexExpr) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpr) String() string {
	var out bytes.Buffer

	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")

	return out.String()
}
