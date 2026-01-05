package lexer

import "github.com/caelondev/monkey/src/token"

type Lexer struct {
	source          string
	lastPosition    int
	currentPosition int
	currentChar     byte
	line            uint
	column          uint
}

func New(source string) *Lexer {
	lexer := Lexer{
		source: source,
		line:   1,
		column: 0,
	}
	lexer.readChar()
	return &lexer
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	for {
		l.skipWhitespace()
		if l.currentChar == '/' && (l.peekChar() == '/' || l.peekChar() == '*') {
			l.skipComments()
		} else {
			break
		}
	}

	// Capture position at START of token
	startLine := l.line
	startColumn := l.column

	switch l.currentChar {
	case ';':
		tok = l.newTokenWithPos(token.SEMICOLON, l.currentChar, startLine, startColumn)
		l.readChar()
	case '^':
		tok = l.newTokenWithPos(token.CARET, l.currentChar, startLine, startColumn)
		l.readChar()
	case '+':
		tok = l.newTokenWithPos(token.PLUS, l.currentChar, startLine, startColumn)
		l.readChar()
	case '-':
		tok = l.newTokenWithPos(token.MINUS, l.currentChar, startLine, startColumn)
		l.readChar()
	case '*':
		tok = l.newTokenWithPos(token.STAR, l.currentChar, startLine, startColumn)
		l.readChar()
	case '/':
		tok = l.newTokenWithPos(token.SLASH, l.currentChar, startLine, startColumn)
		l.readChar()
	case '!':
		tok = l.newCompound(token.BANG, token.NOT_EQUAL, startLine, startColumn)
		l.readChar()
	case '=':
		tok = l.newCompound(token.ASSIGNMENT, token.EQUAL, startLine, startColumn)
		l.readChar()
	case '<':
		tok = l.newCompound(token.LESS, token.LESS_EQUAL, startLine, startColumn)
		l.readChar()
	case '>':
		tok = l.newCompound(token.GREATER, token.GREATER_EQUAL, startLine, startColumn)
		l.readChar()
	case ':':
		tok = l.newTokenWithPos(token.COLON, l.currentChar, startLine, startColumn)
		l.readChar()
	case ',':
		tok = l.newTokenWithPos(token.COMMA, l.currentChar, startLine, startColumn)
		l.readChar()
	case '(':
		tok = l.newTokenWithPos(token.LEFT_PARENTHESIS, l.currentChar, startLine, startColumn)
		l.readChar()
	case ')':
		tok = l.newTokenWithPos(token.RIGHT_PARENTHESIS, l.currentChar, startLine, startColumn)
		l.readChar()
	case '[':
		tok = l.newTokenWithPos(token.LEFT_BRACKET, l.currentChar, startLine, startColumn)
		l.readChar()
	case ']':
		tok = l.newTokenWithPos(token.RIGHT_BRACKET, l.currentChar, startLine, startColumn)
		l.readChar()
	case '{':
		tok = l.newTokenWithPos(token.LEFT_BRACE, l.currentChar, startLine, startColumn)
		l.readChar()
	case '}':
		tok = l.newTokenWithPos(token.RIGHT_BRACE, l.currentChar, startLine, startColumn)
		l.readChar()

	case '\'', '"':
		tok = l.readString(l.currentChar, startLine, startColumn)

	case 0:
		tok = token.Token{Type: token.EOF, Literal: "EOF", Line: startLine, Column: startColumn}
		l.readChar()
	default:
		if isLetter(l.currentChar) {
			literal := l.readIdentifier()
			tok = token.Token{
				Type:    token.LookupIdentifier(literal),
				Literal: literal,
				Line:    startLine,
				Column:  startColumn,
			}
		} else if isNumber(l.currentChar) {
			literal := l.readNumber()
			tok = token.Token{
				Type:    token.NUMBER,
				Literal: literal,
				Line:    startLine,
				Column:  startColumn,
			}
		} else {
			tok = l.newTokenWithPos(token.ILLEGAL, l.currentChar, startLine, startColumn)
			l.readChar()
		}
	}

	return tok
}

func (l *Lexer) readString(terminator byte, line, column uint) token.Token {
	// consume opening quote
	l.readChar()
	start := l.lastPosition

	for {
		if l.currentChar == 0 {
			return l.errorToken("unterminated string", line, column)
		}
		if l.currentChar == '\n' {
			return l.errorToken("string literal cannot span multiple lines", line, column)
		}
		if l.currentChar == terminator {
			lit := l.source[start:l.lastPosition]
			l.readChar() // consume closing quote
			return token.Token{
				Type:    token.STRING,
				Literal: lit,
				Line:    line,
				Column:  column,
			}
		}
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	start := l.lastPosition
	for isNumber(l.currentChar) {
		l.readChar()
	}
	return l.source[start:l.lastPosition]
}

func (l *Lexer) skipWhitespace() {
	for {
		switch l.currentChar {
		case ' ', '\r', '\t':
			l.readChar()
		case '\n':
			l.readChar()
		default:
			return
		}
	}
}

func (l *Lexer) skipComments() {
	if l.currentChar == '/' {
		if l.peekChar() == '/' {
			l.readChar()
			for !l.isAtEnd() && l.currentChar != '\n' {
				l.readChar()
			}
			l.readChar()
		} else if l.peekChar() == '*' {
			l.readChar()
			for !l.isAtEnd() && !(l.currentChar == '*' && l.peekChar() == '/') {
				l.readChar()
			}
			l.readChar()
			l.readChar()
		}
	}
	l.skipWhitespace()
}

func (l *Lexer) readIdentifier() string {
	start := l.lastPosition
	for isAlphanumeric(l.currentChar) {
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

func (l *Lexer) readChar() {
	if l.isAtEnd() {
		l.currentChar = 0
	} else {
		l.currentChar = l.source[l.currentPosition]
	}
	l.lastPosition = l.currentPosition
	l.currentPosition++

	// Track line and column
	if l.currentChar == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) newCompound(regular, compound token.TokenType, line, column uint) token.Token {
	if l.peekChar() == '=' {
		start := l.currentChar
		l.readChar()
		return token.Token{
			Type:    compound,
			Literal: string(start) + string(l.currentChar),
			Line:    line,
			Column:  column,
		}
	} else {
		return l.newTokenWithPos(regular, l.currentChar, line, column)
	}
}

func (l *Lexer) newTokenWithPos(token_type token.TokenType, c byte, line, column uint) token.Token {
	return token.Token{
		Type:    token_type,
		Literal: string(c),
		Line:    line,
		Column:  column,
	}
}

func (l *Lexer) errorToken(msg string, line, column uint) token.Token {
	return token.Token{
		Type:    token.ERROR,
		Literal: msg,
		Line:    line,
		Column:  column,
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.currentPosition >= len(l.source)
}

func isAlphanumeric(ch byte) bool {
	return isLetter(ch) || isNumber(ch)
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
