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

func (e *Evaluator) evaluateToObjectBoolean(node bool) *object.Boolean {
	if node {
		return TRUE
	}
	return FALSE
}

func (e *Evaluator) evaluateUnaryExpression(operator token.TokenType, right object.Object) object.Object {
	switch operator {
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
		return &object.Infinity{Sign: -obj.Sign}
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

func (e *Evaluator) evaluateBinaryExpression(operator token.TokenType, left, right object.Object) object.Object {
	if left.Type() == object.NAN_OBJECT || right.Type() == object.NAN_OBJECT {
		return NAN
	}

	if left.Type() == object.INFINITY_OBJECT && right.Type() == object.INFINITY_OBJECT {
		return evaluateInfinityBinary(operator, left.(*object.Infinity), right.(*object.Infinity))
	}

	if left.Type() == object.INFINITY_OBJECT && right.Type() == object.NUMBER_OBJECT {
		return evaluateInfinityNumber(operator, left.(*object.Infinity), right.(*object.Number))
	}

	if left.Type() == object.NUMBER_OBJECT && right.Type() == object.INFINITY_OBJECT {
		return evaluateNumberInfinity(operator, left.(*object.Number), right.(*object.Infinity))
	}

	if left.Type() == object.NUMBER_OBJECT && right.Type() == object.NUMBER_OBJECT {
		return e.evaluateNumericBinaryExpression(operator, left, right)
	}

	return NIL
}

func evaluateInfinityBinary(operator token.TokenType, left, right *object.Infinity) object.Object {
	switch operator {
	case token.PLUS:
		if left.Sign == right.Sign {
			return &object.Infinity{Sign: left.Sign}
		}
		return NAN
	case token.MINUS:
		if left.Sign != right.Sign {
			return &object.Infinity{Sign: left.Sign}
		}
		return NAN
	case token.STAR:
		return &object.Infinity{Sign: left.Sign * right.Sign}
	case token.SLASH:
		return &object.Infinity{Sign: left.Sign * right.Sign}
	case token.CARET:
		if right.Sign < 0 {
			return NAN
		}
		if left.Sign < 0 && int(right.Sign)%2 != 0 {
			return &object.Infinity{Sign: -1}
		}
		return &object.Infinity{Sign: 1}
	}
	return NIL
}

func evaluateInfinityNumber(operator token.TokenType, left *object.Infinity, right *object.Number) object.Object {
	if math.IsNaN(right.Value) {
		return NAN
	}
	switch operator {
	case token.PLUS, token.MINUS:
		return &object.Infinity{Sign: left.Sign}
	case token.STAR:
		if right.Value == 0 {
			return NAN
		}
		sign := left.Sign
		if right.Value < 0 {
			sign = -sign
		}
		return &object.Infinity{Sign: sign}
	case token.SLASH:
		if right.Value == 0 {
			if left.Sign > 0 {
				return INFINITY
			}
			return NEG_INFINITY
		}
		sign := left.Sign
		if right.Value < 0 {
			sign = -sign
		}
		return &object.Infinity{Sign: sign}
	case token.CARET:
		if right.Value < 0 {
			return &object.Number{Value: 0}
		}
		if left.Sign < 0 && math.Mod(right.Value, 2) != 0 {
			return &object.Infinity{Sign: -1}
		}
		return &object.Infinity{Sign: 1}
	}
	return NAN
}

func evaluateNumberInfinity(operator token.TokenType, left *object.Number, right *object.Infinity) object.Object {
	if math.IsNaN(left.Value) {
		return NAN
	}
	switch operator {
	case token.PLUS:
		return &object.Infinity{Sign: right.Sign}
	case token.MINUS:
		return &object.Infinity{Sign: -right.Sign}
	case token.STAR:
		if left.Value == 0 {
			return NAN
		}

		sign := right.Sign
		if left.Value < 0 {
			sign = -sign
		}

		return &object.Infinity{Sign: sign}
	case token.SLASH:
		if left.Value == 0 {
			return &object.Number{Value: 0}
		}

		sign := right.Sign
		if left.Value < 0 {
			sign = -sign
		}

		return &object.Infinity{Sign: sign}
	case token.CARET:
		if left.Value == 0 {
			return &object.Number{Value: 0}
		}

		if right.Sign < 0 {
			return &object.Number{Value: 0}
		}

		if left.Value < 0 && int(right.Sign)%2 != 0 {
			return &object.Infinity{Sign: -1}
		}

		return &object.Infinity{Sign: 1}
	}

	return NAN
}

func (e *Evaluator) evaluateNumericBinaryExpression(operator token.TokenType, left, right object.Object) object.Object {
	leftVal := left.(*object.Number).Value
	rightVal := right.(*object.Number).Value
	var result float64

	switch operator {
	case token.PLUS:
		result = leftVal + rightVal
	case token.MINUS:
		result = leftVal - rightVal
	case token.STAR:
		result = leftVal * rightVal
	case token.SLASH:
		result = leftVal / rightVal
	case token.CARET:
		result = math.Pow(leftVal, rightVal)

	case token.LESS:
		return e.evaluateToObjectBoolean(leftVal < rightVal)
	case token.GREATER:
		return e.evaluateToObjectBoolean(leftVal > rightVal)
	case token.LESS_EQUAL:
		return e.evaluateToObjectBoolean(leftVal <= rightVal)
	case token.GREATER_EQUAL:
		return e.evaluateToObjectBoolean(leftVal >= rightVal)
	case token.EQUAL:
		return e.evaluateToObjectBoolean(leftVal == rightVal)
	case token.NOT_EQUAL:
		return e.evaluateToObjectBoolean(leftVal != rightVal)
	}

	if math.IsNaN(result) {
		return NAN
	}
	if math.IsInf(result, 64) {
		if result > 0 {
			return INFINITY
		}
		return NEG_INFINITY
	}

	return &object.Number{Value: result}
}
