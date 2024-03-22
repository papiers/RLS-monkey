package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"

	EQ     = "=="
	NOT_EQ = "!="

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
)

// TypeToken 标记类型
type TypeToken string

// Token 标记
type Token struct {
	Type    TypeToken
	Literal string
}

// New 创建标记
func New(typeToken TypeToken, ch byte) Token {
	return Token{Type: typeToken, Literal: string(ch)}
}

// NewString 字符串创建标记
func NewString(typeToken TypeToken, str string) Token {
	return Token{Type: typeToken, Literal: str}
}

var keywords = map[string]TypeToken{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent 返回关键字或标识符的类型
func LookupIdent(ident string) TypeToken {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
