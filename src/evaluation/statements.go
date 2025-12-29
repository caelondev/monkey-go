package evaluation

import (
	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
)

func (e *Evaluator) evaluateIfStatement(condition object.Object, node *ast.IfStatement) object.Object {
	if isTruthy(condition) {
		return e.Evaluate(node.Consequence)
	} else {
		return e.Evaluate(node.Alternative)
	}
}

// TODO: Use if statement for ternary???
func (e *Evaluator) evaluateTernaryExpression(condition object.Object, node *ast.TernaryExpression) object.Object {
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
