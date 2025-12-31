package evaluation

import (
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

func (e *Evaluator) evaluateAssignmentExpression(node *ast.AssignmentExpression, env *object.Environment) object.Object {
	newValue := e.Evaluate(node.NewValue, env)
	assignee := node.Assignee.TokenLiteral()

	if _, ok := env.Get(assignee); ok {
		value, _ := env.Set(assignee, newValue)
		return value
	}

	return e.throwErr(
		node,
		"This error occurs when trying to assign a variable that doesn't exist",
		"Assignment to an undefined variable",
	)
}
