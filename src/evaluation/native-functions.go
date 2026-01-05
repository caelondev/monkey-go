package evaluation

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

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

func (e *Evaluator) NATIVE_TYPE_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	if len(args) != 1 {
		return e.throwErr(
			callNode,
			"This error occurs when trying to pass more than 1 argument value to the function",
			"Expected 1 argument, got %d",
			len(args),
		)
	}

	arg := args[0]
	typeStr := arg.Type()
	return &object.String{Value: string(typeStr)}
}

func (e *Evaluator) NATIVE_TO_NUMBER_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	if len(args) != 1 {
		return e.throwErr(
			callNode,
			"This error occurs when trying to pass more than 1 argument value to the function",
			"Expected 1 argument, got %d",
			len(args),
		)
	}
	switch obj := args[0].(type) {
	case *object.Number, *object.NaN, *object.Infinity:
		return obj
	case *object.String:
		v, err := strconv.ParseFloat(obj.Value, 64)
		if err != nil {
			return object.NAN
		}
		return &object.Number{Value: v}
	case *object.Boolean:
		if obj.Value {
			return &object.Number{Value: 1}
		}
		return &object.Number{Value: 0}
	default:
		return object.NAN
	}
}

func (e *Evaluator) NATIVE_TO_STRING_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	if len(args) != 1 {
		return e.throwErr(
			callNode,
			"This error occurs when trying to pass more than 1 argument value to the function",
			"Expected 1 argument, got %d",
			len(args),
		)
	}
	return &object.String{Value: args[0].Inspect()}
}

func (e *Evaluator) NATIVE_TIME_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	if len(args) != 0 {
		return e.throwErr(
			callNode,
			"This error occurs when trying to pass more than 0 argument value to the function",
			"Expected 0 argument, got %d",
			len(args),
		)
	}
	t := float64(time.Now().UnixNano()) / 1e6 // milliseconds
	return &object.Number{Value: t}
}

func (e *Evaluator) NATIVE_PRINT_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	for i, arg := range args {
		if arg.Type() == object.STRING_OBJECT {
			msg := arg.Inspect()
			fmt.Printf("%s", msg[1:len(msg)-1]) // Trim quotes
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
			"This error occurs when trying to pass more than 1 argument value to the function",
			"Expected 1 argument, got %d",
			len(args),
		)
	}
	message, ok := args[0].(*object.String)
	if !ok {
		return e.throwErr(
			callNode,
			"This error occurs when trying to pass more than 1 argument value to the function",
			"Cannot use prompt message type '%s' as prompt message",
			args[0].Type(),
		)
	}
	fmt.Print(message.Value)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return &object.String{Value: scanner.Text()}
	}
	return e.throwErr(
		callNode,
		"This rare error happens when an I/O error happened during interpretation.\nRestart your program and try again",
		"I/O error",
	)
}

func (e *Evaluator) NATIVE_IS_NAN_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	if len(args) != 1 {
		return e.throwErr(
			callNode,
			"This error occurs when trying to pass more than 1 argument value to the function",
			"Expected 1 argument, got %d",
			len(args),
		)
	}

	isNan := args[0].Type() == object.NAN_OBJECT

	if isNan {
		return object.TRUE
	}

	return object.FALSE
}

func (e *Evaluator) NATIVE_IS_INF_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	if len(args) != 1 {
		return e.throwErr(
			callNode,
			"This error occurs when trying to pass more than 1 argument value to the function",
			"Expected 1 argument, got %d",
			len(args),
		)
	}

	isInf := args[0].Type() == object.INFINITY_OBJECT

	if isInf {
		return object.TRUE
	}

	return object.FALSE
}

func (e *Evaluator) NATIVE_IS_NIL_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	if len(args) != 1 {
		return e.throwErr(
			callNode,
			"This error occurs when trying to pass more than 1 argument value to the function",
			"Expected 1 argument, got %d",
			len(args),
		)
	}

	isInf := args[0].Type() == object.NIL_OBJECT

	if isInf {
		return object.TRUE
	}

	return object.FALSE
}
func (e *Evaluator) NATIVE_RANDOM_FUNCTION(callNode *ast.CallExpression, args []object.Object) object.Object {
	if len(args) != 0 {
		return e.throwErr(
			callNode,
			"This error occurs when trying to pass more than 0 argument value to the function",
			"Expected o argument, got %d",
			len(args),
		)
	}

	return &object.Number{Value: rand.Float64()}
}
