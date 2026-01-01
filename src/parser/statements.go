package parser

import (
	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatements()
	case token.ASSIGN:
		return p.parseBatchAssignStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken() // Advance
	}

	return stmt
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	// SYNTAX ---
	//
	// Null Vals
	// var <Identifier>;
	// var <Identifier>, <Identifier>;
	//
	// Null Vals
	// var <Identifier> = <expr>;
	// var <Identifier>, <Identifier> = <expr>;
	//

	stmt := &ast.VarStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	// First var
	stmt.Names = append(stmt.Names, &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Eat Identifier
		p.nextToken() // Eat comma

		stmt.Names = append(stmt.Names, &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		})
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken() // Eat last var name
		stmt.Value = &ast.NilLiteral{Token: p.currentToken}
		return stmt
	}

	if !p.expectPeek(token.ASSIGNMENT) {
		return nil
	}

	p.nextToken() // Advance ASSIGNMENT

	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// SYNTAX ---
	//
	// return;
	// return <value>;
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken() // Advance past "return" keyword

	// return;
	if p.currentTokenIs(token.SEMICOLON) {
		return stmt
	}

	stmt.ReturnValue = p.parseExpression(LOWEST)

	p.expectPeek(token.SEMICOLON)

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken, Statements: make([]ast.Statement, 0)}

	p.nextToken() // move current token to { ---

	for !p.currentTokenIs(token.RIGHT_BRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken() // Advance next statement
	}

	return block
}

func (p *Parser) parseIfStatements() *ast.IfStatement {
	// Syntax ---
	//
	// if (condition) stmt;
	// if (condition) stmt; else stmt;
	// if (condition) stmt; else if (condition) stmt; else stmt;
	// if (condition) { ... }
	// if (condition) { ... } else { ... }
	// if (condition) { ... } else if (condition) { ... } else { ... } ---
	//
	stmt := &ast.IfStatement{Token: p.currentToken}

	if !p.expectPeek(token.LEFT_PARENTHESIS) { // expectPeek: '(' -> nextToken to '('
		return nil
	}

	p.nextToken() // nextToken: advance into condition (eat '(')
	stmt.Condition = p.parseExpression(LOWEST)
	// current = last token of condition, peek = ')'

	if !p.expectPeek(token.RIGHT_PARENTHESIS) { // expectPeek: ')' -> nextToken to ')'
		return nil
	}
	// current = ')', peek = either '{' (block) or start of one-line stmt

	// Consequence: one-line statement
	if !p.peekTokenIs(token.LEFT_BRACE) {
		p.nextToken() // nextToken: move to start of one-line statement (eat ')')
		stmt.Consequence = p.parseStatement()
		// after parseStatement, peek might be ELSE

		if p.peekTokenIs(token.ELSE) {
			p.nextToken() // nextToken: advance to ELSE (eat whatever is after stmt)
			// check for else-if (hybrid support)
			if p.peekTokenIs(token.IF) {
				p.nextToken() // nextToken: advance to IF
				stmt.Alternative = p.parseIfStatements()
				return stmt
			}

			// else followed by a block or a one-line stmt
			if p.peekTokenIs(token.LEFT_BRACE) {
				p.nextToken() // nextToken: advance to '{'
				stmt.Alternative = p.parseBlockStatement()
			} else {
				p.nextToken() // nextToken: advance to start of else one-line stmt
				stmt.Alternative = p.parseStatement()
			}
		}
		return stmt
	}

	// Consequence: block statement
	p.nextToken() // nextToken: advance to '{'
	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // nextToken: advance to ELSE (eat whatever is after block)
		// else-if after a block
		if p.peekTokenIs(token.IF) {
			p.nextToken() // nextToken: advance to IF
			stmt.Alternative = p.parseIfStatements()
			return stmt
		}

		// else followed by a block or one-line statement
		if p.peekTokenIs(token.LEFT_BRACE) {
			p.nextToken() // nextToken: advance to '{'
			stmt.Alternative = p.parseBlockStatement()
		} else {
			p.nextToken() // nextToken: advance to start of else one-line stmt
			stmt.Alternative = p.parseStatement()
		}
	}

	return stmt
}

func (p *Parser) parseBatchAssignStatement() *ast.BatchAssignmentStatement {
	stmt := &ast.BatchAssignmentStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	firstAssignee := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	stmt.Assignees = append(stmt.Assignees, firstAssignee)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Eat first assignee

		if !p.expectPeek(token.IDENTIFIER) {
			return nil
		}

		assignee := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		stmt.Assignees = append(stmt.Assignees, assignee)
	}

	if !p.expectPeek(token.ASSIGNMENT) {
		return nil
	}

	p.nextToken() // Advance past ASSIGNMENT TOKEN

	stmt.NewValue = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseFunctionStatement() *ast.FunctionDeclarationStatement {
	stmt := &ast.FunctionDeclarationStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.LEFT_PARENTHESIS) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}
