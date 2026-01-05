package parser

import "github.com/caelondev/monkey/src/token"

const (
	_ int = iota
	LOWEST
	ASSIGNMENT
	TERNARY
	OR
	AND
	EQUALITY
	COMPARISON
	ADDITIVE
	MULTIPLICATIVE
	EXPONENTIATION
	UNARY
	CALL
)

var precedence = map[token.TokenType]int{
	token.ASSIGNMENT:       ASSIGNMENT,
	token.IF:               TERNARY,
	token.OR:               OR,
	token.AND:              AND,
	token.EQUAL:            EQUALITY,
	token.NOT_EQUAL:        EQUALITY,
	token.LESS:             COMPARISON,
	token.LESS_EQUAL:       COMPARISON,
	token.GREATER:          COMPARISON,
	token.GREATER_EQUAL:    COMPARISON,
	token.PLUS:             ADDITIVE,
	token.MINUS:            ADDITIVE,
	token.STAR:             MULTIPLICATIVE,
	token.SLASH:            MULTIPLICATIVE,
	token.CARET:            EXPONENTIATION,
	token.NOT:              UNARY,
	token.LEFT_PARENTHESIS: CALL,
	token.LEFT_BRACKET:     CALL,
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Right 10(+)12
func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// Left (10)+12
func (p *Parser) currentPrecedence() int {
	if p, ok := precedence[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) createLookupTable() {
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseNumberExpression)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LEFT_BRACE, p.parseHashLiteral)

	// Array
	p.registerPrefix(token.LEFT_BRACKET, p.parseArrayLiteral)
	p.registerInfix(token.LEFT_BRACKET, p.parseIndexExpression) // Indexing

	p.registerPrefix(token.MINUS, p.parseUnaryExpression)

	p.registerPrefix(token.NIL, p.parseNilLiteral)
	p.registerPrefix(token.NOT_A_NUMBER, p.parseNaNLiteral)
	p.registerPrefix(token.INFINITY, p.parseInfinityLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)

	p.registerInfix(token.PLUS, p.parseBinaryExpression)
	p.registerInfix(token.MINUS, p.parseBinaryExpression)
	p.registerInfix(token.SLASH, p.parseBinaryExpression)
	p.registerInfix(token.STAR, p.parseBinaryExpression)
	p.registerInfix(token.CARET, p.parseExponentExpression)

	p.registerPrefix(token.NOT, p.parseUnaryExpression)
	p.registerInfix(token.AND, p.parseBinaryExpression)
	p.registerInfix(token.OR, p.parseBinaryExpression)
	p.registerInfix(token.EQUAL, p.parseBinaryExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseBinaryExpression)
	p.registerInfix(token.LESS, p.parseBinaryExpression)
	p.registerInfix(token.GREATER, p.parseBinaryExpression)
	p.registerInfix(token.LESS_EQUAL, p.parseBinaryExpression)
	p.registerInfix(token.GREATER_EQUAL, p.parseBinaryExpression)

	p.registerPrefix(token.LEFT_PARENTHESIS, p.parseGroupExpression)
	p.registerInfix(token.IF, p.parseTernaryExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerInfix(token.LEFT_PARENTHESIS, p.parseCallExpression)
	p.registerInfix(token.ASSIGNMENT, p.parseAssignmentExpression)
}
