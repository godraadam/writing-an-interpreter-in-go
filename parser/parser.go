package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

const (
	_ int = iota
	LOWEST
	OR
	AND
	EQ
	LTGT
	SUM
	PROD
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQ,
	token.NOT_EQ:   EQ,
	token.LT:       LTGT,
	token.GT:       LTGT,
	token.LTE:      LTGT,
	token.GTE:      LTGT,
	token.ASTERISK: PROD,
	token.MINUS:    SUM,
	token.PLUS:     SUM,
	token.SLASH:    PROD,
	token.LPAREN:   CALL,
	token.AND:      AND,
	token.OR:       OR,
	token.LBRACKET: INDEX,
}

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
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefixFn(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefixFn(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefixFn(token.STRING, p.parseStringLiteral)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpr)
	p.registerPrefixFn(token.PLUS, p.parsePrefixExpr)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpr)
	p.registerPrefixFn(token.LPAREN, p.parseGroupedExpr)
	p.registerPrefixFn(token.IF, p.parseIfExpr)
	p.registerPrefixFn(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefixFn(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefixFn(token.LBRACE, p.parseMapLiteral)

	p.registerInfixFn(token.MINUS, p.parseInfixExpr)
	p.registerInfixFn(token.OR, p.parseInfixExpr)
	p.registerInfixFn(token.AND, p.parseInfixExpr)
	p.registerInfixFn(token.PLUS, p.parseInfixExpr)
	p.registerInfixFn(token.SLASH, p.parseInfixExpr)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpr)
	p.registerInfixFn(token.EQ, p.parseInfixExpr)
	p.registerInfixFn(token.NOT_EQ, p.parseInfixExpr)
	p.registerInfixFn(token.GT, p.parseInfixExpr)
	p.registerInfixFn(token.GTE, p.parseInfixExpr)
	p.registerInfixFn(token.LT, p.parseInfixExpr)
	p.registerInfixFn(token.LTE, p.parseInfixExpr)
	p.registerInfixFn(token.LPAREN, p.parseCallExpr)
	p.registerInfixFn(token.LBRACKET, p.parseIndexExpr)

	p.advance()
	p.advance()
	return p
}

func (p *Parser) advance() {
	p.currToken = p.peekToken
	p.peekToken = p.l.ScanToken()
}

func (p *Parser) match(t token.TokenType) bool {
	isMatch := p.check(t)
	if isMatch {
		p.advance()
	} else {
		p.addError(p.peekToken.Line, fmt.Sprintf("Expected %s got %s", t, p.peekToken.Literal))
	}
	return isMatch
}

func (p *Parser) check(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) matchOptional(t token.TokenType) bool {
	isMatch := p.check(t)
	if isMatch {
		p.advance()
	}
	return isMatch
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) addError(line int, msg string) {
	p.errors = append(p.errors, fmt.Sprintf("Error on line %d: %s", line, msg))
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

func (p *Parser) parseStmt() ast.Stmt {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	case token.PRINT:
		return p.parsePrintStmt()
	case token.IDENT:
		as := p.parseAssignmentStmt()
		if as != nil {
			return as
		}
		return p.parseExprStmt()
	default:
		return p.parseExprStmt()
	}
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	stmt := &ast.LetStmt{Token: p.currToken}

	if !p.match(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Value: p.currToken.Literal, Token: p.currToken}

	if !p.match(token.ASSIGN) {
		return nil
	}

	p.advance()
	stmt.Value = p.parseExpr(LOWEST)
	p.matchOptional(token.SEMICOLON)

	return stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.currToken}
	p.advance()
	stmt.Value = p.parseExpr(LOWEST)
	p.matchOptional(token.SEMICOLON)

	return stmt
}

func (p *Parser) parsePrintStmt() *ast.PrintStmt {
	stmt := &ast.PrintStmt{Token: p.currToken}
	p.advance()
	stmt.Value = p.parseExpr(LOWEST)
	p.matchOptional(token.SEMICOLON)

	return stmt
}

func (p *Parser) parseAssignmentStmt() *ast.AssingmentStmt {
	name := &ast.Identifier{Value: p.currToken.Literal}
	stmt := &ast.AssingmentStmt{Token: p.currToken, Name: name}

	if !p.check(token.ASSIGN) {
		return nil
	}

	p.advance()
	stmt.Value = p.parseExpr(LOWEST)

	p.matchOptional(token.SEMICOLON)

	return stmt
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{Token: p.currToken}
	stmt.Expr = p.parseExpr(LOWEST)

	p.matchOptional(token.SEMICOLON)

	return stmt
}

func (p *Parser) parseExpr(precedence int) ast.Expr {
	prefixFn := p.prefixParseFns[p.currToken.Type]
	if prefixFn == nil {
		return nil
	}
	leftExpr := prefixFn()
	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil {
			return leftExpr
		}
		p.advance()
		leftExpr = infixFn(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expr {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseNumberLiteral() ast.Expr {
	number, err := strconv.ParseFloat(p.currToken.Literal, 64)
	if err != nil {
		p.addError(p.currToken.Line, fmt.Sprintf("Could not parse number literal %s", p.currToken.Literal))
	}
	return &ast.NumberLiteral{Token: p.currToken, Value: number}
}

func (p *Parser) parseBooleanLiteral() ast.Expr {
	var value bool
	switch p.currToken.Type {
	case token.TRUE:
		value = true
	case token.FALSE:
		value = false
	default:
		p.addError(p.currToken.Line, fmt.Sprintf("Could not parse boolean literal %s", p.currToken.Literal))

	}
	return &ast.BooleanLiteral{Token: p.currToken, Value: value}
}

func (p *Parser) parseStringLiteral() ast.Expr {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expr {
	array := &ast.ArrayLiteral{Token: p.currToken}

	array.Elements = p.parseExprList(token.RBRACKET)

	return array
}

func (p *Parser) parseMapLiteral() ast.Expr {
	theMap := &ast.MapLiteral{Token: p.currToken, Pairs: make(map[ast.Expr]ast.Expr)}

	for !p.check(token.RBRACE) {
		p.advance()
		key := p.parseExpr(LOWEST)

		if !p.match(token.COLON) {
			return nil
		}
		p.advance()
		value := p.parseExpr(LOWEST)
		theMap.Pairs[key] = value

		if !p.check(token.RBRACE) && !p.match(token.COMMA) {
			return nil
		}
	}
	if !p.match(token.RBRACE) {
		return nil
	}

	return theMap
}

func (p *Parser) parsePrefixExpr() ast.Expr {
	prefixExpr := &ast.PrefixExpr{Token: p.currToken, Operator: p.currToken.Literal}
	p.advance()
	prefixExpr.Right = p.parseExpr(PREFIX)
	return prefixExpr
}

func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	infixExpr := &ast.InfixExpr{Token: p.currToken, Left: left, Operator: p.currToken.Literal}
	precedence := p.currPrecedence()
	p.advance()
	infixExpr.Right = p.parseExpr(precedence)
	return infixExpr
}

func (p *Parser) parseGroupedExpr() ast.Expr {
	p.advance()
	expr := p.parseExpr(LOWEST)
	if !p.match(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseIfExpr() ast.Expr {
	ifExpr := &ast.IfExpr{Token: p.currToken}

	// TODO: make parantheses optional in if stmts go-style
	if !p.match(token.LPAREN) {
		return nil
	}
	p.advance()
	ifExpr.Condition = p.parseExpr(LOWEST)

	if !p.match(token.RPAREN) {
		return nil
	}

	ifExpr.Consequence = p.parseBlockStmt()
	if p.matchOptional(token.ELSE) {
		ifExpr.Alternative = p.parseBlockStmt()
	}
	return ifExpr
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	blockStmt := &ast.BlockStmt{Token: p.currToken}
	blockStmt.Stmts = []ast.Stmt{}

	if !p.match(token.LBRACE) {
		return nil
	}

	p.advance()
	for p.currToken.Type != token.RBRACE && p.currToken.Type != token.EOF {
		stmt := p.parseStmt()
		if stmt != nil {
			blockStmt.Stmts = append(blockStmt.Stmts, stmt)
		}
		p.advance()
	}

	return blockStmt
}

func (p *Parser) parseFunctionLiteral() ast.Expr {
	fnLit := &ast.FunctionLiteral{Token: p.currToken}

	fnLit.Params = p.parseFunctionParams()
	fnLit.Body = p.parseBlockStmt()

	return fnLit
}

func (p *Parser) parseExprList(end token.TokenType) []ast.Expr {
	list := []ast.Expr{}

	if p.check(end) {
		p.advance()
		return list
	}

	p.advance()
	list = append(list, p.parseExpr(LOWEST))

	for p.check(token.COMMA) {
		p.advance()
		p.advance()
		list = append(list, p.parseExpr(LOWEST))
	}

	if !p.match(end) {
		return nil
	}
	return list
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	params := []*ast.Identifier{}

	if !p.match(token.LPAREN) {
		return nil
	}

	if p.check(token.RPAREN) {
		p.advance()
		return params
	}
	p.advance()

	param := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	params = append(params, param)

	for p.check(token.COMMA) {
		p.advance()
		p.advance()
		param := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		params = append(params, param)
	}

	if !p.match(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseCallExpr(function ast.Expr) ast.Expr {
	exp := &ast.CallExpr{Token: p.currToken, Function: function}

	exp.Args = p.parseExprList(token.RPAREN)

	return exp
}

func (p *Parser) parseIndexExpr(left ast.Expr) ast.Expr {
	exp := &ast.IndexExpr{Token: p.currToken, Left: left}

	p.advance()
	exp.Index = p.parseExpr(LOWEST)

	if !p.match(token.RBRACKET) {
		return nil
	}
	return exp
}
