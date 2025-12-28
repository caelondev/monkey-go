package ast

import "github.com/caelondev/monkey/src/token"

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (n *NumberLiteral) expressionNode() {}
func (n *NumberLiteral) TokenLiteral() string {
	return n.Token.Literal
}

type Identifier struct {
	Token token.Token // IDENTIFIER Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type UnaryExpression struct {
	Token    token.Token
	Operator token.Token
	Right    Expression
}

func (ue *UnaryExpression) expressionNode() {}
func (ue *UnaryExpression) TokenLiteral() string {
	return ue.Token.Literal
}

type BinaryExpression struct {
	Token    token.Token
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) TokenLiteral() string {
	return be.Token.Literal
}

type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (be *BooleanExpression) expressionNode() {}
func (be *BooleanExpression) TokenLiteral() string {
	return be.Token.Literal
}

type TernaryExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (te *TernaryExpression) expressionNode() {}
func (te *TernaryExpression) TokenLiteral() string {
	return te.Token.Literal
}
