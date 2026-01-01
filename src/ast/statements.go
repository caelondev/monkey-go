package ast

import (
	"bytes"

	"github.com/caelondev/monkey/src/token"
)

// ---------------- VarStatement ----------------
type VarStatement struct {
	Token token.Token   // LET Token
	Names []*Identifier // All names will receive same value
	Value Expression
}

func (vs *VarStatement) GetLine() uint {
	return vs.Token.Line
}
func (vs *VarStatement) GetColumn() uint {
	return vs.Token.Column
}

func (vs *VarStatement) statementNode() {}
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.Token.Literal)
	out.WriteString(" ")

	if len(vs.Names) > 1 {
		for i := 0; i < len(vs.Names)-1; i++ {
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

// ---------------- ReturnStatement ----------------
type ReturnStatement struct {
	Token       token.Token // RETURN Token
	ReturnValue Expression
}

func (rs *ReturnStatement) GetLine() uint {
	return rs.Token.Line
}
func (rs *ReturnStatement) GetColumn() uint {
	return rs.Token.Column
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

// ---------------- ExpressionStatement ----------------
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) GetLine() uint {
	return es.Token.Line
}
func (es *ExpressionStatement) GetColumn() uint {
	return es.Token.Column
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

// ---------------- BlockStatement ----------------
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) GetLine() uint {
	return bs.Token.Line
}
func (bs *BlockStatement) GetColumn() uint {
	return bs.Token.Column
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

// ---------------- IfStatement ----------------
type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence Statement
	Alternative Statement
}

func (is *IfStatement) GetLine() uint {
	return is.Token.Line
}
func (is *IfStatement) GetColumn() uint {
	return is.Token.Column
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

// ---------------- BatchAssignmentStatement ----------------
type BatchAssignmentStatement struct {
	Token     token.Token
	Assignees []*Identifier
	NewValue  Expression
}

func (bs *BatchAssignmentStatement) GetLine() uint {
	return bs.Token.Line
}
func (bs *BatchAssignmentStatement) GetColumn() uint {
	return bs.Token.Column
}

func (ba *BatchAssignmentStatement) statementNode() {}
func (ba *BatchAssignmentStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ba.Token.Literal)

	for i := range len(ba.Assignees) - 1 {
		assignee := ba.Assignees[i]
		out.WriteString(assignee.String())
		out.WriteString(", ")
	}

	assignee := ba.Assignees[len(ba.Assignees)-1]
	out.WriteString(assignee.String())
	out.WriteString(" = ")
	out.WriteString(ba.NewValue.String())

	return out.String()
}
func (ba *BatchAssignmentStatement) TokenLiteral() string {
	return ba.Token.Literal
}

// ---------------- FunctionDeclarationStatement ----------------
type FunctionDeclarationStatement struct {
	Token      token.Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (bs *FunctionDeclarationStatement) GetLine() uint {
	return bs.Token.Line
}
func (bs *FunctionDeclarationStatement) GetColumn() uint {
	return bs.Token.Column
}

func (ba *FunctionDeclarationStatement) statementNode() {}
func (ba *FunctionDeclarationStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ba.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ba.Name.TokenLiteral())
	out.WriteString("(")

	if len(ba.Parameters) != 0 {
		for i := range len(ba.Parameters) - 1 {

			param := ba.Parameters[i]
			out.WriteString(param.String())
			out.WriteString(", ")
		}

	param := ba.Parameters[len(ba.Parameters)-1]
	out.WriteString(param.String())
	}

	out.WriteString(") {\n")

	for _, stmt := range ba.Body.Statements {
		out.WriteString("\t")
		out.WriteString(stmt.String())
		out.WriteString("\n")
	}

	out.WriteString("}\n")

	return out.String()
}
func (ba *FunctionDeclarationStatement) TokenLiteral() string {
	return ba.Token.Literal
}
