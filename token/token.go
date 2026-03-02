package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

const (
	ILLEGAL  = "ILLEGAL"
	EOF      = "EOF"
	IDENT    = "IDENT"
	NUMBER   = "NUMBER"
	STRING   = "STRING"
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	BANG      = "!"
	QUOTES    = "\""
	ELLIPSIS  = "..."

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="
	LTE    = ">="
	GTE    = "<="

	OR  = "||"
	AND = "&&"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	PRINT    = "PRINT"
	WHILE    = "WHILE"
	NIL      = "NIL"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"print":  PRINT,
	"while":  WHILE,
	"nil":    NIL,
}

func LookupIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
