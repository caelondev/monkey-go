package lexer

import (
	"github.com/caelondev/monkey/src/token"
)

type Lexer struct {
	source          string
	lastPosition    int
	currentPosition int
	currentChar     byte
}

func New(source string) *Lexer {
	lexer := Lexer{source: source}

	// Initialize first char
	lexer.readChar()
	return &lexer
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	l.skipComments()

	switch l.currentChar {
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar)
		l.readChar()
	case '+':
		tok = newToken(token.PLUS, l.currentChar)
		l.readChar()
	case '-':
		tok = newToken(token.MINUS, l.currentChar)
		l.readChar()
	case '*':
		tok = newToken(token.STAR, l.currentChar)
		l.readChar()
	case '/':
		tok = newToken(token.SLASH, l.currentChar)
		l.readChar()
	case '!':
		tok = l.newCompound(token.BANG, token.NOT_EQUAL)
		l.readChar()
	case '=':
		tok = l.newCompound(token.ASSIGNMENT, token.EQUAL)
		l.readChar()
	case '<':
		tok = newToken(token.LESS, l.currentChar)
		l.readChar()
	case '>':
		tok = newToken(token.GREATER, l.currentChar)
		l.readChar()
	case ',':
		tok = newToken(token.COMMA, l.currentChar)
		l.readChar()
	case '(':
		tok = newToken(token.LEFT_PARENTHESIS, l.currentChar)
		l.readChar()
	case ')':
		tok = newToken(token.RIGHT_PARENTHESIS, l.currentChar)
		l.readChar()
	case '{':
		tok = newToken(token.LEFT_BRACE, l.currentChar)
		l.readChar()
	case '}':
		tok = newToken(token.RIGHT_BRACE, l.currentChar)
		l.readChar()
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
		l.readChar()
	default:
		if isLetter(l.currentChar) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
		} else if isNumber(l.currentChar) {
			tok.Literal = l.readNumber()
			tok.Type = token.NUMBER
		} else {
			tok = newToken(token.ILLEGAL, l.currentChar)
			l.readChar()
		}
	}

	return tok
}

func (l *Lexer) readNumber() string {
	start := l.lastPosition

	for isNumber(l.currentChar) {
		l.readChar()
	}

	return l.source[start:l.lastPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' ||
		l.currentChar == '\t' ||
		l.currentChar == '\r' ||
		l.currentChar == '\n' {
		l.readChar()
	}
}

func (l *Lexer) skipComments() {
	if l.currentChar == '/' {
		if l.peekChar() == '/' {
			l.readChar() // Advance first '/'

			for !l.isAtEnd() && l.currentChar != '\n' {
				l.readChar()
			}

			l.readChar() // Eat '\n'
		} else if l.peekChar() == '*' { // Selective comment
			l.readChar() // Eat '/'

			for !l.isAtEnd() && !(l.currentChar == '*' && l.peekChar() == '/') {
				l.readChar()
			}

			// Consume '*/'
			l.readChar()
			l.readChar()
		}
	}

	// Skip whitespace after comment
	l.skipWhitespace()
}

func (l *Lexer) readIdentifier() string {
	start := l.lastPosition

	for isLetter(l.currentChar) {
		l.readChar()
	}

	return l.source[start:l.lastPosition]
}

func (l *Lexer) peekChar() byte {
	if l.isAtEnd() {
		return 0
	}

	return l.source[l.currentPosition]
}

// Advance
func (l *Lexer) readChar() {
	if l.isAtEnd() {
		l.currentChar = 0
	} else {
		l.currentChar = l.source[l.currentPosition]
	}

	l.lastPosition = l.currentPosition
	l.currentPosition++
}

func (l *Lexer) newCompound(regular, compound token.TokenType) token.Token {
	if l.peekChar() == '=' {
		start := l.currentChar
		l.readChar()
		return token.Token{
			Type:    compound,
			Literal: string(start) + string(l.currentChar),
		}
	} else {
		return newToken(regular, l.currentChar)
	}
}

func newToken(token_type token.TokenType, c byte) token.Token {
	return token.Token{Type: token_type, Literal: string(c)}
}

func (l *Lexer) isAtEnd() bool {
	return l.currentPosition >= len(l.source)
}

func isLetter(ch byte) bool {
	return ('a' <= ch && 'z' >= ch) || ('A' <= ch && 'Z' >= ch) || ch == '_'
}

func isNumber(ch byte) bool {
	return '0' <= ch && '9' >= ch
}
