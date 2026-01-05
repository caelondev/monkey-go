package ast

import (
	"bytes"
	"strings"

	"github.com/caelondev/monkey/src/token"
)

// ---------------- NumberLiteral ----------------
type StringLiteral struct {
	Token token.Token
	Value string
}

func (n *StringLiteral) GetLine() uint {
	return n.Token.Line
}
func (n *StringLiteral) GetColumn() uint {
	return n.Token.Column
}

func (n *StringLiteral) expressionNode() {}
func (n *StringLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("\"")
	out.WriteString(n.Token.Literal)
	out.WriteString("\"")
	return out.String()
}
func (n *StringLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// ---------------- NumberLiteral ----------------
type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (n *NumberLiteral) GetLine() uint {
	return n.Token.Line
}

func (n *NumberLiteral) GetColumn() uint {
	return n.Token.Column
}

func (n *NumberLiteral) expressionNode() {}
func (n *NumberLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(n.Token.Literal)
	return out.String()
}
func (n *NumberLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// ---------------- NilLiteral ----------------
type NilLiteral struct {
	Token token.Token
}

func (n *NilLiteral) GetLine() uint {
	return n.Token.Line
}
func (n *NilLiteral) GetColumn() uint {
	return n.Token.Column
}

func (n *NilLiteral) expressionNode() {}
func (n *NilLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("nil")
	return out.String()
}
func (n *NilLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// ---------------- Identifier ----------------
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) GetLine() uint {
	return i.Token.Line
}
func (i *Identifier) GetColumn() uint {
	return i.Token.Column
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	var out bytes.Buffer
	out.WriteString(i.Value)
	return out.String()
}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// ---------------- UnaryExpression ----------------
type UnaryExpression struct {
	Token    token.Token
	Operator token.Token
	Right    Expression
}

func (ue *UnaryExpression) GetLine() uint {
	return ue.Token.Line
}
func (ue *UnaryExpression) GetColumn() uint {
	return ue.Token.Column
}

func (ue *UnaryExpression) expressionNode() {}
func (ue *UnaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ue.Operator.Literal)
	out.WriteString(ue.Right.String())
	return out.String()
}
func (ue *UnaryExpression) TokenLiteral() string {
	return ue.Token.Literal
}

// ---------------- BinaryExpression ----------------
type BinaryExpression struct {
	Token    token.Token
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (be *BinaryExpression) GetLine() uint {
	return be.Token.Line
}
func (be *BinaryExpression) GetColumn() uint {
	return be.Token.Column
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString(be.Left.String())
	out.WriteString(be.Operator.Literal)
	out.WriteString(be.Right.String())
	return out.String()
}
func (be *BinaryExpression) TokenLiteral() string {
	return be.Token.Literal
}

// ---------------- BooleanExpression ----------------
type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (be *BooleanExpression) GetLine() uint {
	return be.Token.Line
}
func (be *BooleanExpression) GetColumn() uint {
	return be.Token.Column
}

func (be *BooleanExpression) expressionNode() {}
func (be *BooleanExpression) String() string {
	var out bytes.Buffer
	out.WriteString(be.Token.Literal)
	return out.String()
}
func (be *BooleanExpression) TokenLiteral() string {
	return be.Token.Literal
}

// ---------------- TernaryExpression ----------------
type TernaryExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (te *TernaryExpression) GetLine() uint {
	return te.Token.Line
}
func (te *TernaryExpression) GetColumn() uint {
	return te.Token.Column
}

func (te *TernaryExpression) expressionNode() {}
func (te *TernaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString(te.Consequence.String())
	out.WriteString(" if ")
	out.WriteString(te.Condition.String())
	out.WriteString(" else ")
	out.WriteString(te.Alternative.String())
	return out.String()
}
func (te *TernaryExpression) TokenLiteral() string {
	return te.Token.Literal
}

// ---------------- FunctionLiteral ----------------
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) GetLine() uint {
	return fl.Token.Line
}
func (fl *FunctionLiteral) GetColumn() uint {
	return fl.Token.Column
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
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

	return out.String()
}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// ---------------- CallExpression ----------------
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) GetLine() uint {
	return ce.Token.Line
}
func (ce *CallExpression) GetColumn() uint {
	return ce.Token.Column
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	for i, arg := range ce.Arguments {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(arg.String())
	}
	out.WriteString(")")
	return out.String()
}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// ---------------- NaNLiteral ----------------
type NaNLiteral struct {
	Token token.Token
}

func (n *NaNLiteral) GetLine() uint {
	return n.Token.Line
}
func (n *NaNLiteral) GetColumn() uint {
	return n.Token.Column
}

func (n *NaNLiteral) expressionNode() {}
func (n *NaNLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("NaN")
	return out.String()
}
func (n *NaNLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// ---------------- InfinityLiteral ----------------
type InfinityLiteral struct {
	Token token.Token
	Sign  int // -1 +1 ---
}

func (n *InfinityLiteral) GetLine() uint {
	return n.Token.Line
}
func (n *InfinityLiteral) GetColumn() uint {
	return n.Token.Column
}

func (n *InfinityLiteral) expressionNode() {}
func (n *InfinityLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("Inf")
	return out.String()
}
func (n *InfinityLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// ---------------- AssignmentExpression ----------------
type AssignmentExpression struct {
	Token    token.Token
	Assignee Expression
	NewValue Expression
}

func (n *AssignmentExpression) GetLine() uint {
	return n.Token.Line
}
func (n *AssignmentExpression) GetColumn() uint {
	return n.Token.Column
}

func (n *AssignmentExpression) expressionNode() {}
func (n *AssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString(n.Assignee.String())
	out.WriteString(" = ")
	out.WriteString(n.NewValue.String())
	return out.String()
}
func (n *AssignmentExpression) TokenLiteral() string {
	return n.Token.Literal
}

// ---------------- ArrayLiteral ----------------
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (n *ArrayLiteral) GetLine() uint {
	return n.Token.Line
}
func (n *ArrayLiteral) GetColumn() uint {
	return n.Token.Column
}

func (n *ArrayLiteral) expressionNode() {}
func (n *ArrayLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("[")

	for i, elem := range n.Elements {
		out.WriteString(elem.String())
		if i != len(n.Elements)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString("]")

	return out.String()
}

func (n *ArrayLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// ---------------- IndexExpression ----------------
type IndexExpression struct {
	Token  token.Token
	Index  Expression
	Target Expression
}

func (n *IndexExpression) GetLine() uint {
	return n.Token.Line
}
func (n *IndexExpression) GetColumn() uint {
	return n.Token.Column
}

func (n *IndexExpression) expressionNode() {}
func (n *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString(n.Target.String())
	out.WriteString("[")
	out.WriteString(n.Index.String())
	out.WriteString("]")

	return out.String()
}
func (n *IndexExpression) TokenLiteral() string {
	return n.Token.Literal
}

// ---------------- IndexAssignmentExpression ----------------
type IndexAssignmentExpression struct {
	Token    token.Token
	Index    Expression
	Target   Expression
	NewValue Expression
}

func (n *IndexAssignmentExpression) GetLine() uint {
	return n.Token.Line
}
func (n *IndexAssignmentExpression) GetColumn() uint {
	return n.Token.Column
}

func (n *IndexAssignmentExpression) expressionNode() {}
func (n *IndexAssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(n.Target.String())
	out.WriteString("[")
	out.WriteString(n.Index.String())
	out.WriteString("] = ")
	out.WriteString(n.NewValue.String())

	return out.String()
}
func (n *IndexAssignmentExpression) TokenLiteral() string {
	return n.Token.Literal
}

// ---------------- HashLiteral ----------------
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (n *HashLiteral) GetLine() uint {
	return n.Token.Line
}
func (n *HashLiteral) GetColumn() uint {
	return n.Token.Column
}

func (n *HashLiteral) expressionNode() {}
func (n *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range n.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	if len(pairs) <= 0 {
		out.WriteString("{}")
	} else {
		out.WriteString("{")
		out.WriteString(strings.Join(pairs, ", "))
		out.WriteString("}")
	}
	return out.String()
}
func (n *HashLiteral) TokenLiteral() string {
	return n.Token.Literal
}
