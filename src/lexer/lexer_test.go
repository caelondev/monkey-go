package lexer

import (
	"testing"

	"github.com/caelondev/monkey/src/token"
)

func TestNextToken(t *testing.T) {
	input := `
var five = 5;
var ten = 10;
var add = fn(x, y) {
    x + y;
};
var result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true
} else {
	return false
}

10 == 10;
10 != 5;
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// var five = 5;
		{token.VAR, "var"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGNMENT, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},

		// var ten = 10;
		{token.VAR, "var"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGNMENT, "="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},

		// var add = fn(x, y) { x + y; };
		{token.VAR, "var"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGNMENT, "="},
		{token.FUNCTION, "fn"},
		{token.LEFT_PARENTHESIS, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.LEFT_BRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.SEMICOLON, ";"},

		// var result = add(five, ten);
		{token.VAR, "var"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGNMENT, "="},
		{token.IDENTIFIER, "add"},
		{token.LEFT_PARENTHESIS, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.SEMICOLON, ";"},

		// !-/*5;
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.STAR, "*"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},

		// 5 < 10 > 5;
		{token.NUMBER, "5"},
		{token.LESS, "<"},
		{token.NUMBER, "10"},
		{token.GREATER, ">"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},

		// if (5 < 10) { return true }
		{token.IF, "if"},
		{token.LEFT_PARENTHESIS, "("},
		{token.NUMBER, "5"},
		{token.LESS, "<"},
		{token.NUMBER, "10"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.RIGHT_BRACE, "}"},

		// else { return false }
		{token.ELSE, "else"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.RIGHT_BRACE, "}"},

		{token.NUMBER, "10"},
		{token.EQUAL, "=="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "10"},
		{token.NOT_EQUAL, "!="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf(
				"tests[%d] - tokentype wrong. expected=%q, got=%q",
				i,
				tt.expectedType,
				tok.Type,
			)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"tests[%d] - literal wrong. expected=%q, got=%q",
				i,
				tt.expectedLiteral,
				tok.Literal,
			)
		}
	}
}
