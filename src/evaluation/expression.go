package evaluation

import (
	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
	"github.com/caelondev/monkey/src/token"
	"github.com/sanity-io/litter"
)

// Cache to avoid making new object instance every call
var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func (e *Evaluator) evaluateExpression(node ast.Expression) object.Object {
	switch node := node.(type) {
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.NilLiteral:
		return NIL
	case *ast.BooleanExpression:
		return e.evaluateToObjectBoolean(node.Value)
	case *ast.UnaryExpression:
		right := e.evaluateExpression(node.Right)
		return e.evaluateUnaryExpression(node.Operator.Type, right)

	default:
		println("Unrecognized AST node:\n")
		litter.Dump(node)
		return NIL
	}
}

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
	if right.Type() != object.NUMBER_OBJECT {
		return NIL
	}

	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}

func (e *Evaluator) evaluateNotExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NIL:
		return TRUE

	default:
		return FALSE
	}
}
