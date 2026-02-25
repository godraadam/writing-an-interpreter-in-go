package lexer

import (
	"fmt"
	"monkey/token"
)

type Lexer struct {
	source  string
	current int // current
	start   int // start
	line    int
}

func New(source string) *Lexer {
	l := &Lexer{source: source}
	return l
}

func isAlpha(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '_' // allow for syntax like 100_000
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

// jumps start to current
func (l *Lexer) jump() {
	l.start = l.current
}

// read current char then advance
func (l *Lexer) advance() byte {
	ch := l.peek()
	l.current++
	return ch
}

// read current char and advance if it matches expected
func (l *Lexer) match(expected byte) bool {
	ch := l.peek()
	if ch != expected {
		return false
	}
	l.current++
	return true
}

// read current char and advance if it matches expected as per fn
func (l *Lexer) matchExpr(fn func(byte) bool) bool {
	ch := l.peek()
	if !fn(ch) {
		return false
	}
	l.current++
	return true
}

// read current char without advancing
func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	ch := l.source[l.current]
	return ch
}

// read char after current without advancing
func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	ch := l.source[l.current+1]
	return ch
}

// get next char after skipping all whitespace
func (l *Lexer) skipWhitespace() byte {
	for l.matchExpr(isWhiteSpace) {
		if l.peek() == '\n' {
			l.line++
		}
	}
	l.jump()
	return l.advance()
}

// advance until new line
func (l *Lexer) skipLineComment() {
	for l.peek() != '\n' && !l.isAtEnd() {
		l.advance()
	}
	l.jump()
}

// greedily consume characters to match keword or identifer
func (l *Lexer) readIdentifier() token.Token {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}
	tok := l.createToken(token.LookupIdentifier(l.source[l.start:l.current]))
	return tok
}

// greedily consume characters to match a number (including fractional part)
func (l *Lexer) readNumber() token.Token {
	for isDigit(l.peek()) {
		l.advance()
	}
	// look for fractionial part
	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance() // jump over '.'
		for isDigit(l.peek()) {
			l.advance()
		}
	}
	return l.createToken(token.NUMBER)

}

func (l *Lexer) readString() token.Token {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}
	if l.isAtEnd() {
		reportError(l.line, "Unterminated string literal")
	}
	l.advance()
	return l.createToken(token.STRING)

}

func (l *Lexer) createToken(ttype token.TokenType) token.Token {
	return token.Token{Type: ttype, Literal: (l.source[l.start:l.current]), Line: l.line}
}

func (l *Lexer) ScanToken() token.Token {
	var tok token.Token

	ch := l.skipWhitespace()

	switch ch {
	case '=':
		if l.match('=') {
			tok = l.createToken(token.EQ)
		} else {
			tok = l.createToken(token.ASSIGN)
		}
	case '.':
		if l.match('.') {
			if l.match('.') {
				tok = l.createToken(token.ELLIPSIS)
			}
		}
	case ';':
		tok = l.createToken(token.SEMICOLON)
	case ':':
		tok = l.createToken(token.COLON)
	case '|':
		if l.match('|') {
			tok = l.createToken(token.OR)
		} else {
			reportError(l.line, fmt.Sprintf("Illegal token %q", ch))

		}
	case '&':
		if l.match('&') {
			tok = l.createToken(token.AND)
		} else {
			reportError(l.line, fmt.Sprintf("Illegal token %q", ch))

		}
	case '(':
		tok = l.createToken(token.LPAREN)
	case ')':
		tok = l.createToken(token.RPAREN)
	case ',':
		tok = l.createToken(token.COMMA)
	case '+':
		tok = l.createToken(token.PLUS)
	case '-':
		tok = l.createToken(token.MINUS)
	case '{':
		tok = l.createToken(token.LBRACE)
	case '[':
		tok = l.createToken(token.LBRACKET)
	case ']':
		tok = l.createToken(token.RBRACKET)
	case '}':
		tok = l.createToken(token.RBRACE)
	case '"':
		tok = l.readString()
	case '/':
		if l.match('/') {
			l.skipLineComment()
		} else {
			tok = l.createToken(token.SLASH)
		}
		tok = l.createToken(token.SLASH)
	case '*':
		tok = l.createToken(token.ASTERISK)
	case '<':
		if l.match('=') {
			tok = l.createToken(token.LTE)
		} else {
			tok = l.createToken(token.LT)
		}
	case '>':
		if l.match('=') {
			tok = l.createToken(token.GTE)
		} else {
			tok = l.createToken(token.GT)
		}
	case '!':
		if l.match('=') {
			tok = l.createToken(token.NOT_EQ)
		} else {
			tok = l.createToken(token.BANG)
		}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: l.line}
	default:
		if isAlpha(ch) {
			tok = l.readIdentifier()
		} else if isDigit(ch) {
			tok = l.readNumber()
		} else {
			reportError(l.line, fmt.Sprintf("Illegal token %q", ch))
			tok = l.createToken(token.ILLEGAL)
		}
	}
	l.jump()

	return tok
}

func (l *Lexer) ScanTokens() []token.Token {
	var tokens []token.Token
	for !l.isAtEnd() {
		l.ScanToken()
	}
	return tokens
}

func reportError(line int, message string) {
	fmt.Printf("Error [line %d]: %s\n", line, message)
}
