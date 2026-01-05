package evaluation

import (
	"fmt"

	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
)

func (e *Evaluator) evaluateToObjectBoolean(v bool) *object.Boolean {
	if v {
		return object.TRUE
	}
	return object.FALSE
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

func (e *Evaluator) throwErr(node ast.Node, hint string, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Line:    node.GetLine(),
		Column:  node.GetColumn(),
		Message: fmt.Sprintf(format, a...),
		Hint:    hint,
		NodeStr: node.String(),
	}
}

func (e *Evaluator) unwrapReturnValue(obj object.Object) object.Object {
	if returnVal, ok := obj.(*object.ReturnValue); ok {
		return returnVal.Value
	}

	return obj
}
