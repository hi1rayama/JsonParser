package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"

	EOF = "EOF"

	// 識別子 + リテラル
	IDENT = "IDENT" 
	DIGITS   = "DIGITS"   // 1234567890

	TRUE  = "TRUE"
	FALSE = "FALSE"
	NULL  = "NULL"

	PLUS = "+"
	MINUS = "-"

	// デリミタ
	COMMA      = ","
	PRIOD      = "."
	COLON      = ":"
	DOUBLEQUOT = "\""

	LBRACKET = "["
	RBRACKET = "]"
	LBRACE   = "{"
	RBRACE   = "}"
)

var keywords = map[string]TokenType{
	"true":   TRUE,
	"false":  FALSE,
	"null": NULL,
}

// LookupIdent 
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}