package evaluation

import (
	"math"

	"github.com/caelondev/monkey/src/ast"
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

func (e *Evaluator) evaluateUnaryExpression(node *ast.UnaryExpression, right object.Object) object.Object {
	switch node.Operator.Type {
	case token.BANG:
		return e.evaluateNotExpression(right)
	case token.MINUS:
		return e.evaluateNegationExpression(node, right)
	default:
		return e.throwErr(
			node,
			"This error occurs when an unregistered unary operator was used.\nThis should only appear during language development.",
			"Unknown unary operator: '%v'",
			node.Operator.Type,
		)
	}
}

func (e *Evaluator) evaluateNegationExpression(node *ast.UnaryExpression, right object.Object) object.Object {
	switch obj := right.(type) {
	case *object.Infinity:
		return infinityWithSign(-obj.Sign)
	case *object.Number:
		return &object.Number{Value: -obj.Value}
	case *object.NaN:
		return NAN
	default:
		return e.throwErr(
			node,
			"This error occurs when you try to negate a non-number value.",
			"Cannot negate operand of type '%v'", right.Type(),
		)
	}
}

func (e *Evaluator) evaluateNotExpression(right object.Object) object.Object {
	if !isTruthy(right) {
		return TRUE
	}
	return FALSE
}

func (e *Evaluator) evaluateBinaryExpression(
	node *ast.BinaryExpression,
	left, right object.Object,
) object.Object {
	// Handle NaN operands ---
	if left.Type() == object.NAN_OBJECT || right.Type() == object.NAN_OBJECT {
		switch node.Operator.Type {
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
		return evalInfInf(node.Operator.Type, left.(*object.Infinity), right.(*object.Infinity))

	case left.Type() == object.INFINITY_OBJECT && right.Type() == object.NUMBER_OBJECT:
		return evalInfNum(node.Operator.Type, left.(*object.Infinity), right.(*object.Number))

	case left.Type() == object.NUMBER_OBJECT && right.Type() == object.INFINITY_OBJECT:
		return evalNumInf(node.Operator.Type, left.(*object.Number), right.(*object.Infinity))

	case left.Type() == object.NUMBER_OBJECT && right.Type() == object.NUMBER_OBJECT:
		return e.evaluateNumericBinaryExpression(node, left, right)
	}

	return e.throwErr(
		node,
		"This error occurs when the operands don't share the same type or cannot be used together to perform arithmetic.",
		"Cannot perform `%v %v %v` as they are an invalid operand combination",
		left.Type(),
		node.Operator.Type,
		right.Type(),
	)
}

func (e *Evaluator) evaluateNumericBinaryExpression(
	node *ast.BinaryExpression,
	left, right object.Object,
) object.Object {
	l := left.(*object.Number).Value
	r := right.(*object.Number).Value
	var result float64

	switch node.Operator.Type {
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

	default:
		return e.throwErr(
			node,
			"This error occurs when an unregistered binary operator is used.\nThis should only appear during language development.",
			"Unknown binary operator: '%v'",
			node.Operator.Type,
		)
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
