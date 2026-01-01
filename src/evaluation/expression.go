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

func (e *Evaluator) evaluateIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	e.line = node.GetLine()
	e.column = node.GetColumn()

	if value, ok := env.Get(node.Value); ok {
		return value
	}

	return e.throwErr(
		node,
		"This error happens when a variable with that given name doesn't exist",
		"Cannot resolve variable '%s'",
		node.Value,
	)
}

func (e *Evaluator) evaluateExpressions(
	exprs []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, expr := range exprs {
		evaluated := e.Evaluate(expr, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) evaluateUnaryExpression(node *ast.UnaryExpression, env *object.Environment) object.Object {
	right := e.Evaluate(node.Right, env)
	if isError(right) {
		return right
	}

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

func (e *Evaluator) evaluateTernaryExpression(node *ast.TernaryExpression, env *object.Environment) object.Object {
	condition := e.Evaluate(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Evaluate(node.Consequence, env)
	} else {
		return e.Evaluate(node.Alternative, env)
	}
}

func (e *Evaluator) evaluateNotExpression(right object.Object) object.Object {
	if !isTruthy(right) {
		return TRUE
	}
	return FALSE
}

func (e *Evaluator) evaluateBinaryExpression(node *ast.BinaryExpression, env *object.Environment) object.Object {
	left := e.Evaluate(node.Left, env)

	if isError(left) {
		return left
	}
	right := e.Evaluate(node.Right, env)

	if isError(right) {
		return right
	}

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

func (e *Evaluator) evaluateCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	fn := e.Evaluate(node.Function, env)
	if isError(fn) {
		return fn
	}

	args := e.evaluateExpressions(node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	return e.applyFunction(node.Function, node, fn, args)
}

func (e *Evaluator) applyFunction(fnNode, callNode ast.Node, function object.Object, args []object.Object) object.Object {
	fn, ok := function.(*object.FunctionLiteral)
	if !ok {
		return e.throwErr(
			fnNode,
			"This error occurs when you're trying to call a non-function expression",
			"Invalid call to a non-function expression",
		)
	}

	if len(fn.Parameters) != len(args) {
		return e.throwErr(
			callNode,
			"This error occurs when the amount of argument doesnt match the function parameter",
			"Argument count mismatch. Expected %d argument count, got %d",
			len(fn.Parameters),
			len(args),
		)
	}

	extendedEnv := e.extendFunctionEnv(fn, args)
	evaluated := e.Evaluate(fn.Body, extendedEnv)
	return e.unwrapFunctionValue(evaluated)
}

func (e *Evaluator) extendFunctionEnv(fn *object.FunctionLiteral, args []object.Object) *object.Environment {
	// fn env is the outer env (for closure) ---
	env := object.NewEnclosedEnvironment(fn.Scope)

	for idx, param := range fn.Parameters {
		env.Set(param.Value, args[idx])
	}

	return env
}

func (e *Evaluator) unwrapFunctionValue(evaluated object.Object) object.Object {
	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		return returnValue
	}

	return evaluated
}
