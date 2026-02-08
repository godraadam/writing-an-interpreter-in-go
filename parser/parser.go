package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

const (
	_ int = iota
	LOWEST
	EQ
	LTGT
	SUM
	PROD
	PREFIX
	CALL
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	return p
}

func (p *Parser) advance() {
	p.currToken = p.peekToken
	p.peekToken = p.l.ScanToken()
}

func (p *Parser) matchPeek(t token.TokenType) bool {
	isMatch := p.peekToken.Type == t
	if isMatch {
		p.advance()
	} else {
		p.addError(p.peekToken.Line, fmt.Sprintf("Expected %s got %s", t, p.peekToken.Literal))
	}
	return isMatch
}

func (p *Parser) addError(line int, msg string) {
	p.errors = append(p.errors, fmt.Sprintf("Error on line %d: %s", line, msg))
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Stmts = []ast.Stmt{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStmt()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
		p.advance()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	default:
		return p.parseExprStmt()
	}
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	stmt := &ast.LetStmt{Token: p.currToken}

	if !p.matchPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Value: p.currToken.Literal, Token: p.currToken}

	if !p.matchPeek(token.ASSIGN) {
		return nil
	}

	// TODO parse expression
	for p.currToken.Type != token.SEMICOLON {
		p.advance()
	}
	return stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.currToken}
	p.advance()
	// TODO parse expression
	for p.currToken.Type != token.SEMICOLON {
		p.advance()
	}
	return stmt
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{Token: p.currToken}
	stmt.Expression = p.parseExpr(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.advance()
	}
	return stmt
}

func (p *Parser) parseExpr(precedence int) ast.Expr {
	prefixFn := p.prefixParseFns[p.currToken.Type]
	if prefixFn == nil {
		return nil
	}
	leftExpr := prefixFn()

	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expr {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}
