package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/caelondev/monkey/src/ast"
)

type ObjectType string

const (
	NUMBER_OBJECT       = "NUMBER"
	STRING_OBJECT       = "STRING"
	ARRAY_OBJECT        = "ARRAY"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NIL_OBJECT          = "NIL"
	NAN_OBJECT          = "NAN"
	INFINITY_OBJECT     = "INFINITY"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION"
	HASH_OBJECT         = "HASH"
)

var (
	NIL          = &Nil{}
	INFINITY     = &Infinity{Sign: 1}
	NEG_INFINITY = &Infinity{Sign: -1}
	NAN          = &NaN{}
	TRUE         = &Boolean{Value: true}
	FALSE        = &Boolean{Value: false}
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (o *Hash) Type() ObjectType {
	return HASH_OBJECT
}

func (o *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range o.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
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

type String struct {
	Value string
}

func (o *String) Type() ObjectType {
	return STRING_OBJECT
}

func (o *String) Inspect() string {
	return fmt.Sprintf("\"%s\"", o.Value)
}

func (o *String) HashKey() HashKey {
	hash := fnv.New64a()
	hash.Write([]byte(o.Value))

	return HashKey{Type: o.Type(), Value: hash.Sum64()}
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

func (o *Number) HashKey() HashKey {
	return HashKey{Type: o.Type(), Value: uint64(o.Value)}
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

func (o *Boolean) HashKey() HashKey {
	var value uint64

	if o.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: o.Type(), Value: value}
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

type Function struct {
	Parameters []*ast.Identifier
	Name       *ast.Identifier
	Body       *ast.BlockStatement
	Scope      *Environment
}

func (o *Function) Type() ObjectType {
	return FUNCTION_OBJECT
}

func (o *Function) Inspect() string {
	if o.Name == nil {
		return "[ Anonymous Function ]"
	}

	return fmt.Sprintf("[ Function '%s' ]", o.Name)
}

type NativeFunctionFn func(
	callNode *ast.CallExpression,
	args []Object,
) Object

type NativeFunction struct {
	Fn NativeFunctionFn
}

func (o *NativeFunction) Type() ObjectType {
	return FUNCTION_OBJECT
}

func (o *NativeFunction) Inspect() string {
	return "[ Native Function ]"
}

type Array struct {
	Elements []Object
}

func (o *Array) Type() ObjectType {
	return ARRAY_OBJECT
}

func (o *Array) Inspect() string {
	var out bytes.Buffer

	out.WriteString("[")

	for i, elem := range o.Elements {
		out.WriteString(elem.Inspect())
		if i != len(o.Elements)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString("]")

	return out.String()
}
