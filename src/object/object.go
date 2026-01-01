package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/caelondev/monkey/src/ast"
)

type ObjectType string

const (
	NUMBER_OBJECT       = "NUMBER"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NIL_OBJECT          = "NIL"
	NAN_OBJECT          = "NAN"
	INFINITY_OBJECT     = "INFINITY"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION_OBJECT"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Number struct {
	Value float64
}

func (o *Number) Type() ObjectType {
	return NUMBER_OBJECT
}

func (o *Number) Inspect() string {
	return fmt.Sprintf("%g", o.Value)
}

type Boolean struct {
	Value bool
}

func (o *Boolean) Type() ObjectType {
	return BOOLEAN_OBJECT
}

func (o *Boolean) Inspect() string {
	return fmt.Sprintf("%t", o.Value)
}

type Nil struct{}

func (o *Nil) Type() ObjectType {
	return NIL_OBJECT
}

func (o *Nil) Inspect() string {
	return "nil"
}

type NaN struct{}

func (o *NaN) Type() ObjectType {
	return NAN_OBJECT
}

func (o *NaN) Inspect() string {
	return "NotANumber"
}

type Infinity struct {
	Sign int
}

func (o *Infinity) Type() ObjectType {
	return INFINITY_OBJECT
}

func (o *Infinity) Inspect() string {
	if o.Sign > 0 {
		return "Infinity++"
	}

	return "Infinity--"
}

type ReturnValue struct {
	Value Object
}

func (o *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJECT
}

func (o *ReturnValue) Inspect() string {
	return fmt.Sprintf("return { %s }", o.Value.Inspect())
}

type Error struct {
	Line    uint
	Column  uint
	Message string
	Hint    string
	NodeStr string
}

func (o *Error) Type() ObjectType {
	return ERROR_OBJECT
}

func (o *Error) Inspect() string {
	return fmt.Sprintf("Error at Ln %d:%d - %s", o.Line, o.Column, o.Message)
}

type FunctionLiteral struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Scope      *Environment
}

func (o *FunctionLiteral) Type() ObjectType {
	return FUNCTION_OBJECT
}

func (o *FunctionLiteral) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range o.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(o.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type FunctionStatement struct {
	Parameters []*ast.Identifier
	Name       *ast.Identifier
	Body       *ast.BlockStatement
	Scope      *Environment
}

func (o *FunctionStatement) Type() ObjectType {
	return FUNCTION_OBJECT
}

func (o *FunctionStatement) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range o.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn ")
	out.WriteString(o.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(o.Body.String())
	out.WriteString("\n}")
	return out.String()
}
