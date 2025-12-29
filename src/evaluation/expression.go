package evaluation

import (
	"math"

	"github.com/caelondev/monkey/src/object"
	"github.com/caelondev/monkey/src/token"
)

var (
	NIL          = &object.Nil{}
	INFINITY     = &object.Infinity{Sign: 1}
	NEG_INFINITY = &object.Infinity{Sign: -1}
	NAN          = &object.NaN{}
	TRUE         = &object.Boolean{Value: true}
	FALSE        = &object.Boolean{Value: false}
)

func (e *Evaluator) evaluateToObjectBoolean(v bool) *object.Boolean {
	if v {
		return TRUE
	}
	return FALSE
}

func (e *Evaluator) evaluateUnaryExpression(op token.TokenType, right object.Object) object.Object {
	switch op {
	case token.BANG:
		return e.evaluateNotExpression(right)
	case token.MINUS:
		return e.evaluateNegationExpression(right)
	default:
		return NIL
	}
}

func (e *Evaluator) evaluateNegationExpression(right object.Object) object.Object {
	switch obj := right.(type) {
	case *object.Infinity:
		return infinityWithSign(-obj.Sign)
	case *object.Number:
		return &object.Number{Value: -obj.Value}
	case *object.NaN:
		return NAN
	default:
		return NAN
	}
}

func (e *Evaluator) evaluateNotExpression(right object.Object) object.Object {
	if !isTruthy(right) {
		return TRUE
	}
	return FALSE
}

func (e *Evaluator) evaluateBinaryExpression(
	op token.TokenType,
	left, right object.Object,
) object.Object {
	if left.Type() == object.NAN_OBJECT || right.Type() == object.NAN_OBJECT {
		switch op {
		case token.EQUAL, token.LESS, token.GREATER, token.LESS_EQUAL, token.GREATER_EQUAL:
			return FALSE
		case token.NOT_EQUAL:
			return TRUE
		default:
			return NAN
		}
	}

	switch {
	case left.Type() == object.INFINITY_OBJECT && right.Type() == object.INFINITY_OBJECT:
		return evalInfInf(op, left.(*object.Infinity), right.(*object.Infinity))

	case left.Type() == object.INFINITY_OBJECT && right.Type() == object.NUMBER_OBJECT:
		return evalInfNum(op, left.(*object.Infinity), right.(*object.Number))

	case left.Type() == object.NUMBER_OBJECT && right.Type() == object.INFINITY_OBJECT:
		return evalNumInf(op, left.(*object.Number), right.(*object.Infinity))

	case left.Type() == object.NUMBER_OBJECT && right.Type() == object.NUMBER_OBJECT:
		return e.evaluateNumericBinaryExpression(op, left, right)
	}

	return NIL
}

func (e *Evaluator) evaluateNumericBinaryExpression(
	op token.TokenType,
	left, right object.Object,
) object.Object {

	l := left.(*object.Number).Value
	r := right.(*object.Number).Value
	var result float64

	switch op {
	case token.PLUS:
		result = l + r
	case token.MINUS:
		result = l - r
	case token.STAR:
		result = l * r
	case token.SLASH:
		result = l / r
	case token.CARET:
		result = math.Pow(l, r)

	case token.LESS:
		return e.evaluateToObjectBoolean(l < r)
	case token.GREATER:
		return e.evaluateToObjectBoolean(l > r)
	case token.LESS_EQUAL:
		return e.evaluateToObjectBoolean(l <= r)
	case token.GREATER_EQUAL:
		return e.evaluateToObjectBoolean(l >= r)
	case token.EQUAL:
		return e.evaluateToObjectBoolean(l == r)
	case token.NOT_EQUAL:
		return e.evaluateToObjectBoolean(l != r)
	}

	if math.IsNaN(result) {
		return NAN
	}

	if math.IsInf(result, 1) {
		return INFINITY
	}
	if math.IsInf(result, -1) {
		return NEG_INFINITY
	}

	return &object.Number{Value: result}
}
