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

// new LetStmt: let <identifier> | <destructuring_expr> = <expr>
// destructuring_expr: <array_destructuring_expr> | <map_destructuring_expr>
// array_destructuring_expr: [<identifer> (, <identifier>)*, <ellipsis_expr>?] | [<ellipsis_expr> (, <identifier>*)]
// map_destructuring_expr: {<identifier> (, <identifier>)*, <ellipsis_expr>?}
// ellipsis_expr: ...<identifer>

// Let statement
type LetStmt struct {
	Token  token.Token
	Target Expr
	Value  Expr
}

func (l *LetStmt) stmtNode() {}
func (l *LetStmt) TokenLiteral() string {
	return l.Token.Literal
}
func (l *LetStmt) String() string {
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Target.String())
	out.WriteString(" = ")

	if l.Value != nil {
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

type WhileStmt struct {
	Token token.Token
	Cond  Expr
	Body  *BlockStmt
}

func (ws *WhileStmt) stmtNode() {}
func (ws *WhileStmt) TokenLiteral() string {
	return ws.Token.Literal
}

func (ws *WhileStmt) String() string {
	var out bytes.Buffer
	out.WriteString(ws.TokenLiteral())
	out.WriteString("(" + ws.Cond.String() + ")")
	out.WriteString(" {")
	out.WriteString(ws.Body.String() + "}")
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

type MapLiteral struct {
	Token token.Token
	Pairs map[Expr]Expr
}

func (ml *MapLiteral) exprNode() {}
func (ml *MapLiteral) TokenLiteral() string {
	return ml.Token.Literal
}

func (ml *MapLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	pairs := []string{}
	for k, v := range ml.Pairs {
		pairs = append(pairs, k.String()+": "+v.String())
	}
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

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

func (ie *IfExpr) exprNode() {}
func (ie *IfExpr) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpr) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" ")
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type AssignExpr struct {
	Token token.Token
	Field Expr
	Value Expr
}

func (ae *AssignExpr) exprNode() {}
func (ae *AssignExpr) TokenLiteral() string {
	return ae.Token.Literal
}

func (ae *AssignExpr) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Field.String())
	out.WriteString(" = ")
	out.WriteString(ae.Value.String())
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

type EllipsisExpr struct {
	Token token.Token
	Name  Identifier
}

func (ee *EllipsisExpr) exprNode() {}
func (ee *EllipsisExpr) TokenLiteral() string {
	return ee.Token.Literal
}

func (ee *EllipsisExpr) String() string {
	return ee.TokenLiteral() + ee.Name.String()
}

// ellipsis where?
type ArrayDestructuringExpr struct {
	Token                token.Token
	Names                []Identifier
	EllipsisExprPosition int // -1 for none
	EllipsisExpr         *EllipsisExpr
}

func (ade *ArrayDestructuringExpr) exprNode() {}
func (ade *ArrayDestructuringExpr) TokenLiteral() string {
	return ade.Token.Literal
}

func (ade *ArrayDestructuringExpr) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	elements := []string{}
	for i, el := range ade.Names {
		if ade.EllipsisExpr != nil && i == ade.EllipsisExprPosition {
			elements = append(elements, ade.EllipsisExpr.String())
		}
		elements = append(elements, el.String())
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type MapDestructuringExpr struct {
	Token        token.Token
	Names        []Identifier
	EllipsisExpr *EllipsisExpr
}

func (mde *MapDestructuringExpr) exprNode() {}
func (mde *MapDestructuringExpr) TokenLiteral() string {
	return mde.Token.Literal
}

func (mde *MapDestructuringExpr) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	elements := []string{}
	for _, el := range mde.Names {
		elements = append(elements, el.String())
	}
	out.WriteString(strings.Join(elements, ", "))
	if mde.EllipsisExpr != nil {
		out.WriteString(", ")
		out.WriteString(mde.EllipsisExpr.String())
	}
	out.WriteString("}")

	return out.String()
}
