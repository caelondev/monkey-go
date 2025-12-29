package parser

import "github.com/caelondev/monkey/src/token"

const (
	_ int = iota
	LOWEST
	TERNARY
	EQUALITY
	COMPARISON
	ADDITIVE
	MULTIPLICATIVE
	UNARY
	CALL
)

var precedence = map[token.TokenType]int{
	token.EQUAL:            EQUALITY,
	token.NOT_EQUAL:        EQUALITY,
	token.LESS:             COMPARISON,
	token.GREATER:          COMPARISON,
	token.PLUS:             ADDITIVE,
	token.MINUS:            ADDITIVE,
	token.STAR:             MULTIPLICATIVE,
	token.SLASH:            MULTIPLICATIVE,
	token.LEFT_PARENTHESIS: CALL,
	token.IF:               TERNARY,
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
	p.registerPrefix(token.BANG, p.parseUnaryExpression)
	p.registerPrefix(token.MINUS, p.parseUnaryExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)

	p.registerInfix(token.PLUS, p.parseBinaryExpression)
	p.registerInfix(token.MINUS, p.parseBinaryExpression)
	p.registerInfix(token.SLASH, p.parseBinaryExpression)
	p.registerInfix(token.STAR, p.parseBinaryExpression)

	p.registerInfix(token.EQUAL, p.parseBinaryExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseBinaryExpression)
	p.registerInfix(token.LESS, p.parseBinaryExpression)
	p.registerInfix(token.GREATER, p.parseBinaryExpression)

	p.registerPrefix(token.LEFT_PARENTHESIS, p.parseGroupExpression)
	p.registerInfix(token.IF, p.parseTernaryExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerInfix(token.LEFT_PARENTHESIS, p.parseCallExpression)
}
