package evaluation

import (
	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
)

func (e *Evaluator) evaluateBlockStatement(node *ast.BlockStatement) object.Object {
	var lastEvaluated object.Object

	for _, stmt := range node.Statements {
		lastEvaluated = e.Evaluate(stmt)

		if lastEvaluated != nil {
			if lastEvaluated.Type() == object.RETURN_VALUE_OBJECT || lastEvaluated.Type() == object.ERROR_OBJECT {
				return lastEvaluated
			}
		}
	}

	return lastEvaluated
}

func (e *Evaluator) evaluateIfStatement(node *ast.IfStatement, condition object.Object) object.Object {
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Evaluate(node.Consequence)
	} else {
		if node.Alternative == nil {
			return NIL
		}

		return e.Evaluate(node.Alternative)
	}
}

func (e *Evaluator) evaluateTernaryExpression(node *ast.TernaryExpression, condition object.Object) object.Object {
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Evaluate(node.Consequence)
	} else {
		return e.Evaluate(node.Alternative)
	}
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Nil:
		return false

	case *object.Boolean:
		return obj.Value

	case *object.Number:
		return obj.Value != 0

	case *object.NaN:
		return false

	default:
		return true
	}
}
