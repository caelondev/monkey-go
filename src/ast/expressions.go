package ast

import (
	"bytes"

	"github.com/caelondev/monkey/src/token"
)

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (n *NumberLiteral) expressionNode() {}
func (n *NumberLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(n.Token.Literal)
	out.WriteString(")")

	return out.String()
}
func (n *NumberLiteral) TokenLiteral() string {
	return n.Token.Literal
}

type NilLiteral struct {
	Token token.Token
}

func (n *NilLiteral) expressionNode() {}
func (n *NilLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString("nil")
	out.WriteString(")")

	return out.String()
}
func (n *NilLiteral) TokenLiteral() string {
	return n.Token.Literal
}

type Identifier struct {
	Token token.Token // IDENTIFIER Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Value)
	out.WriteString(")")

	return out.String()
}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type UnaryExpression struct {
	Token    token.Token
	Operator token.Token
	Right    Expression
}

func (ue *UnaryExpression) expressionNode() {}
func (ue *UnaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ue.Operator.Literal)
	out.WriteString(ue.Right.String())
	out.WriteString(")")

	return out.String()
}
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
func (be *BinaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(be.Left.String())
	out.WriteString(be.Operator.Literal)
	out.WriteString(be.Right.String())
	out.WriteString(")")

	return out.String()
}
func (be *BinaryExpression) TokenLiteral() string {
	return be.Token.Literal
}

type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (be *BooleanExpression) expressionNode() {}
func (be *BooleanExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(be.Token.Literal)
	out.WriteString(")")

	return out.String()
}
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
func (te *TernaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(te.Consequence.String())
	out.WriteString(" if ")
	out.WriteString(te.Condition.String())
	out.WriteString(" else ")
	out.WriteString(te.Alternative.String())
	out.WriteString(")")

	return out.String()
}
func (te *TernaryExpression) TokenLiteral() string {
	return te.Token.Literal
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(fl.Token.Literal)
	out.WriteString("(")

	for i, param := range fl.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(param.String())
	}

	out.WriteString(") ")
	out.WriteString("{\n")
	out.WriteString(fl.Body.String())
	out.WriteString("}")
	out.WriteString(")")

	return out.String()
}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ce.Function.String())
	out.WriteString("(")

	for i, arg := range ce.Arguments {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(arg.String())
	}

	out.WriteString("))")

	return out.String()
}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
