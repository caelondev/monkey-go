package evaluation

import (
	"fmt"

	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
)

func (e *Evaluator) NATIVE_LEN_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	arg := args[0]

	switch arg.Type() {
	case object.STRING_OBJECT:
		s, _ := arg.(*object.String)
		return &object.Number{Value: float64(len(s.Value))}

	default:
		return e.throwErr(
			callNode.Arguments[0],
			"This error occurs when trying to get the length of an unsupported value",
			"Cannot get length of type '%s'",
			arg.Type(),
		)
	}
}

func (e *Evaluator) NATIVE_PRINT_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	for i, arg := range args {
		fmt.Printf("%s", arg.Inspect())

		if i != len(args)-1 {
			fmt.Printf(", ")
		}
	}

	fmt.Println()
	return object.NIL
}
