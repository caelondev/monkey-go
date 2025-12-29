package object

import "fmt"

type ObjectType string

const (
	NUMBER_OBJECT   = "NUMBER"
	BOOLEAN_OBJECT  = "BOOLEAN"
	NIL_OBJECT      = "NIL"
	NAN_OBJECT      = "NAN"
	INFINITY_OBJECT = "INFINITY"
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
