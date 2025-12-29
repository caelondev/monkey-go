package object

import "fmt"

type ObjectType string

const (
	NUMBER_OBJECT  = "NUMBER"
	BOOLEAN_OBJECT = "BOOLEAN"
	NIL_OBJECT     = "NIL"
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
	return NUMBER_OBJECT
}

func (o *Nil) Inspect() string {
	return "nil"
}
