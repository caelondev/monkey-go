package parser

import (
	"fmt"
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
	if leftExpression == nil {
		return nil
	}

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

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()                     // Eat ( token
	expr := p.parseExpression(LOWEST) // Use LOWEST, not CALL

	if !p.expectPeek(token.RIGHT_PARENTHESIS) { // Consume the )
		return nil
	}

	return expr
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	var value bool

	switch p.currentToken.Type {
	case token.TRUE:
		value = true
	case token.FALSE:
		value = false
	}

	return &ast.BooleanExpression{Token: p.currentToken, Value: value}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	expr := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeek(token.LEFT_PARENTHESIS) {
		return nil
	}

	expr.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseNilLiteral() ast.Expression {
	return &ast.NilLiteral{Token: p.currentToken}
}
func (p *Parser) parseInfinityLiteral() ast.Expression {
	return &ast.InfinityLiteral{Token: p.currentToken, Sign: 1}
}

func (p *Parser) parseNaNLiteral() ast.Expression {
	return &ast.NaNLiteral{Token: p.currentToken}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	expr := &ast.ArrayLiteral{Token: p.currentToken}

	if p.peekTokenIs(token.RIGHT_BRACKET) {
		p.nextToken()
		return expr
	}

	p.nextToken() // Eat [

	firstElem := p.parseExpression(LOWEST)
	expr.Elements = append(expr.Elements, firstElem)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Eat expr
		p.nextToken() // Eat comma

		elem := p.parseExpression(LOWEST)
		expr.Elements = append(expr.Elements, elem)
	}

	if !p.expectPeek(token.RIGHT_BRACKET) {
		return nil
	}

	return expr
}

func (p *Parser) parseHashLiteral() ast.Expression {
	expr := &ast.HashLiteral{Token: p.currentToken}
	expr.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RIGHT_BRACE) {
		p.nextToken() // Move to key position
		key := p.parseExpression(LOWEST)
		if key == nil {
			return nil
		}

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken() // Move to value position
		value := p.parseExpression(LOWEST)
		if value == nil {
			return nil
		}

		expr.Pairs[key] = value // ACTUALLY STORE IT!

		// Handle comma or closing brace
		if !p.peekTokenIs(token.RIGHT_BRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RIGHT_BRACE) {
		return nil
	}

	return expr
}

/*
* [ INFIX EXPRESSIONS ]
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

func (p *Parser) parseTernaryExpression(left ast.Expression) ast.Expression {
	// Syntax ---
	//
	// <consequence> if <condition> else <alternate>
	//

	expr := &ast.TernaryExpression{Token: p.currentToken}
	expr.Consequence = left

	pre := p.currentPrecedence()
	p.nextToken() // Eat IF Token
	expr.Condition = p.parseExpression(pre)

	if !p.expectPeek(token.ELSE) {
		return nil
	}

	p.nextToken() // Eat ELSE

	expr.Alternative = p.parseExpression(pre)

	return expr
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Token: p.currentToken}
	expr.Function = left
	expr.Arguments = p.parseCallArguments()

	return expr
}

func (p *Parser) parseExponentExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinaryExpression{Token: p.currentToken, Operator: p.currentToken, Left: left}

	// Right associative parsing
	pre := p.currentPrecedence() - 1
	p.nextToken() // Eat CARET
	expr.Right = p.parseExpression(pre)

	return expr
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	if left == nil {
		return nil
	}

	// Allow both identifiers AND index expressions
	switch leftExpr := left.(type) {
	case *ast.Identifier:
		// Regular variable assignment: a = 10
		expr := &ast.AssignmentExpression{
			Token:    p.currentToken,
			Assignee: leftExpr,
		}
		p.nextToken()
		expr.NewValue = p.parseExpression(ASSIGNMENT + 1)

		if expr.NewValue == nil {
			p.errors = append(p.errors, fmt.Sprintf(
				"[Ln %d:%d] Invalid right-hand side in assignment",
				p.currentToken.Line, p.currentToken.Column))
			return nil
		}
		return expr

	case *ast.IndexExpression:
		expr := &ast.IndexAssignmentExpression{
			Token:  p.currentToken,
			Target: leftExpr.Target, // Take Index Expr's array target
			Index:  leftExpr.Index,
		}
		p.nextToken()
		expr.NewValue = p.parseExpression(ASSIGNMENT + 1)

		if expr.NewValue == nil {
			p.errors = append(p.errors, fmt.Sprintf(
				"[Ln %d:%d] Invalid right-hand side in assignment",
				p.currentToken.Line, p.currentToken.Column))
			return nil
		}
		return expr

	default:
		p.throwError(
			"[Ln %d:%d] Cannot reassign to non-identifier/non-index expression '%s'",
			p.currentToken.Line,
			p.currentToken.Column,
			p.currentToken.Type,
		)
		return nil
	}
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Token: p.currentToken, Target: left}

	p.nextToken() // Eat [ ---

	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_BRACKET) {
		return nil
	}

	return expr
}

/*
* [ HELPERS ]
**/

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	idents := make([]*ast.Identifier, 0)

	// Check if no args passed
	if p.peekTokenIs(token.RIGHT_PARENTHESIS) {
		p.nextToken() // Eat ( ---
		return idents // Return empty
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	firstParam := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	idents = append(idents, firstParam)

	// Will run every comma, and automatically
	// jumps to it
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Eat first param
		if !p.expectPeek(token.IDENTIFIER) {
			return nil
		}

		param := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		idents = append(idents, param)
	}

	p.nextToken() // Eat Ident ---
	return idents
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := make([]ast.Expression, 0)

	// Currently at ( ---

	if p.peekTokenIs(token.RIGHT_PARENTHESIS) {
		p.nextToken() // Advance ) ---
		return args
	}

	// Eat ( ---
	p.nextToken()
	firstArg := p.parseExpression(LOWEST)
	args = append(args, firstArg)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Advance past Ident
		p.nextToken() // Advance comma

		arg := p.parseExpression(LOWEST)
		args = append(args, arg)
	}

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}

	return args
}
