package ast

import "bytes"

type Node interface {
	String() string
	GetLine() uint
	GetColumn() uint
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
		out.WriteString(";\n")
	}

	return out.String()
}

func (p *Program) GetLine() uint {
	return 0
}

func (p *Program) GetColumn() uint {
	return 0
}
