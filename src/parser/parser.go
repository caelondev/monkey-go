package parser

import (
	"fmt"

	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/lexer"
	"github.com/caelondev/monkey/src/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token
	errors       []string
	hadError     bool

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         make([]string, 0),
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	p.createLookupTable()

	// Initialize currentToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = make([]ast.Statement, 0)

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		// Eat semicolon
		p.nextToken()
	}

	return program
}

func (p *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return p.currentToken.Type == tokenType
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekToken.Type == tokenType {
		p.nextToken()
		return true
	}

	p.peekError(tokenType)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	p.throwError(
		"[Ln %d:%d] Expected token after '%s' to be %s, got '%s' instead",
		p.currentToken.Line,
		p.currentToken.Column,
		p.currentToken.Literal,
		t,
		p.peekToken.Literal,
	)
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.throwError(
		"[Ln %d:%d] -> Unexpected token found: '%s'",
		p.currentToken.Line,
		p.currentToken.Column,
		t,
	)
}

func (p *Parser) throwError(format string, a ...interface{}) {
	if p.hadError {
		return
	}

	p.hadError = true
	msg := fmt.Sprintf(format, a...)
	p.errors = append(p.errors, msg)
}
