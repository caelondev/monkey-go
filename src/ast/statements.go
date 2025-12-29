package ast

import (
	"bytes"

	"github.com/caelondev/monkey/src/token"
)

type VarStatement struct {
	Token token.Token   // LET Token
	Names []*Identifier // All names will receive same value
	Value Expression
}

func (vs *VarStatement) statementNode() {}
func (vs *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.Token.Literal)
	out.WriteString(" ")

	if len(vs.Names) > 1 {
		for i := range len(vs.Names) - 1 {
			out.WriteString(vs.Names[i].String())
			out.WriteString(", ")
		}
	}

	out.WriteString(vs.Names[len(vs.Names)-1].String())

	out.WriteString(" = ")
	out.WriteString(vs.Value.String())

	return out.String()
}
func (vs *VarStatement) TokenLiteral() string {
	return vs.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token // RETURN Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.Token.Literal)

	if rs.ReturnValue != nil {
		out.WriteString(" ")
		out.WriteString(rs.ReturnValue.String())
	}

	return out.String()
}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(es.Expression.String())
	out.WriteString(")")

	return out.String()
}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString("\t")
		out.WriteString(stmt.String())
		out.WriteString(";\n")
	}

	return out.String()
}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence Statement
	Alternative Statement
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString(is.Token.Literal)
	out.WriteString(" (")
	out.WriteString(is.Condition.String())
	out.WriteString(") {\n")
	out.WriteString(is.Consequence.String())
	out.WriteString("} ")

	if is.Alternative != nil {
		out.WriteString("else {\n")
		out.WriteString(is.Alternative.String())
		out.WriteString("}")
	}

	return out.String()
}
func (is *IfStatement) TokenLiteral() string {
	return is.Token.Literal
}
