package parser

import (
	"strconv"

	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/token"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		p.nextToken()                          // Advance left and inspect operator
		leftExpression = infix(leftExpression) // Bubble up left
	}

	// Returns left ONLY IF THE NEXT TOKEN IS A SEMICOLON
	// and THE PRECEDENCE IS LOWER THAN THE CURRENT PRECEDENCE
	return leftExpression
}

/*
* [ PREFIX EXPRESSIONS ]
**/

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseNumberExpression() ast.Expression {
	// TODO: Errors are being ignored, this might cause a crash
	// when the tokenzer invalidly tokenize a numerical token
	value, _ := strconv.ParseFloat(p.currentToken.Literal, 64)
	return &ast.NumberLiteral{Token: p.currentToken, Value: value}
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	expr := &ast.UnaryExpression{Token: p.currentToken, Operator: p.currentToken}
	p.nextToken() // Advance past unary operator

	expr.Right = p.parseExpression(UNARY)

	return expr
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()                     // Eat ( token
	expr := p.parseExpression(LOWEST) // Use LOWEST, not CALL

	if !p.expectPeek(token.RIGHT_PARENTHESIS) { // Consume the )
		return nil
	}

	return expr
}

/*
* [ PREFIX EXPRESSIONS ]
**/

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinaryExpression{
		Token:    p.currentToken,
		Operator: p.currentToken,
		Left:     left,
	}

	pre := p.currentPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(pre)

	return expr
}
