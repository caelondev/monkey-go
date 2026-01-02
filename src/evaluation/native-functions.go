package evaluation

import (
	"bufio"
	"fmt"
	"os"

	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
)

func (e *Evaluator) NATIVE_LEN_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	if len(args) != 1 {
		return e.throwErr(
			callNode,
			"This error occurs when an argument passed was less than or greater than expected amount",
			"Expected 1 argument, got %d",
			len(args),
		)
	}

	arg := args[0]

	switch arg.Type() {
	case object.STRING_OBJECT:
		s, _ := arg.(*object.String)
		return &object.Number{Value: float64(len(s.Value))}
	case object.ARRAY_OBJECT:
		a, _ := arg.(*object.Array)
		return &object.Number{Value: float64(len(a.Elements))}

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
		if arg.Type() == object.STRING_OBJECT {
			msg := arg.Inspect()
			trimmed := msg[1 : len(msg)-1] // Trim ""
			fmt.Printf("%s", trimmed)
		} else {
			fmt.Printf("%s", arg.Inspect())
		}

		if i != len(args)-1 {
			fmt.Printf(", ")
		}
	}

	fmt.Println()
	return object.NIL
}

func (e *Evaluator) NATIVE_PROMPT_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	if len(args) != 1 {
		return e.throwErr(
			callNode,
			"This error occurs when trying to pass more than one argument",
			"Expected 1 argument, got %d",
			len(args),
		)
	}

	arg := args[0]
	message, ok := arg.(*object.String)

	if !ok {
		return e.throwErr(
			callNode.Arguments[0],
			"This error occurs when trying to pass a non-string value as a prompt",
			"Cannot send value type of '%s' as message prompt",
			arg.Type(),
		)
	}

	fmt.Print(message.Value)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return &object.String{Value: scanner.Text()}
	}

	return e.throwErr(
		callNode,
		"This rare error occurs when the Stdin failed to return the input of the prompt",
		"Failed I/O error",
	)
}
