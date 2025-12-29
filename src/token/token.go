package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    uint
	Column  uint
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifier + literal
	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"

	// Operators
	ASSIGNMENT = "="
	PLUS       = "+"
	MINUS      = "-"
	BANG       = "!"
	STAR       = "*"
	SLASH      = "/"

	LESS    = "<"
	GREATER = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"
	LEFT_BRACE        = "{"
	RIGHT_BRACE       = "}"

	// Two chars
	EQUAL     = "=="
	NOT_EQUAL = "!="

	// Reserved keywords
	FUNCTION = "FUNCTION"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	NIL      = "NIL"
)

var reservedKeywords = map[string]TokenType{
	"fn":     FUNCTION,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"nil":    NIL,
}

func LookupIdentifier(ident string) TokenType {
	if tok, ok := reservedKeywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}
