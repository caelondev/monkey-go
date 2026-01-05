package evaluation

import (
	"math"

	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
	"github.com/caelondev/monkey/src/token"
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
	case token.NOT:
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
		return object.NAN
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
		return object.TRUE
	}
	return object.FALSE
}

func (e *Evaluator) evaluateBinaryExpression(node *ast.BinaryExpression, env *object.Environment) object.Object {
	switch node.Operator.Type {
	case token.OR, token.AND:
		return e.evaluateComparisonExpression(node, env)
	}

	left := e.unwrapReturnValue(e.Evaluate(node.Left, env))

	if isError(left) {
		return left
	}

	right := e.unwrapReturnValue(e.Evaluate(node.Right, env))

	if isError(right) {
		return right
	}

	// Handle NaN operands ---
	if left.Type() == object.NAN_OBJECT || right.Type() == object.NAN_OBJECT {
		switch node.Operator.Type {
		case token.EQUAL, token.LESS, token.GREATER, token.LESS_EQUAL, token.GREATER_EQUAL:
			return object.FALSE
		case token.NOT_EQUAL:
			return object.TRUE
		default:
			return object.NAN
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
	case left.Type() == object.STRING_OBJECT && right.Type() == object.STRING_OBJECT:
		return e.evaluateStringBinaryExpression(node, left, right)
	case left.Type() == object.BOOLEAN_OBJECT && right.Type() == object.BOOLEAN_OBJECT:
		return e.evaluateBooleanBinaryExpression(node, left, right)
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

func (e *Evaluator) evaluateBooleanBinaryExpression(node *ast.BinaryExpression, left, right object.Object) object.Object {
	switch node.Operator.Type {
	case token.EQUAL:
		return e.evaluateToObjectBoolean(left.(*object.Boolean).Value == right.(*object.Boolean).Value)
	case token.NOT_EQUAL:
		return e.evaluateToObjectBoolean(left.(*object.Boolean).Value != right.(*object.Boolean).Value)
	default:
		return object.NIL // Unreachable
	}
}

func (e *Evaluator) evaluateComparisonExpression(node *ast.BinaryExpression, env *object.Environment) object.Object {
	left := e.Evaluate(node.Left, env)
	right := e.Evaluate(node.Right, env)

	var result bool

	switch node.Operator.Type {
	case token.AND:
		result = isTruthy(left) && isTruthy(right)
	case token.OR:
		result = isTruthy(left) || isTruthy(right)
	}

	// Prevents multiple boolean allocation ---
	if result {
		return object.TRUE
	} else {
		return object.FALSE
	}
}

func (e *Evaluator) evaluateStringBinaryExpression(node *ast.BinaryExpression, left, right object.Object) object.Object {
	l := left.(*object.String).Value
	r := right.(*object.String).Value

	var result string

	switch node.Operator.Type {
	case token.PLUS:
		result = l + r

	default:
		return e.throwErr(
			node,
			"This error occurs when an invalid string operator was used",
			"Invalid string operator '%s'",
			node.Operator.Literal,
		)
	}

	return &object.String{Value: result}
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
		return object.NAN
	}
	if math.IsInf(result, 1) {
		return object.INFINITY
	}
	if math.IsInf(result, -1) {
		return object.NEG_INFINITY
	}

	return &object.Number{Value: result}
}

func (e *Evaluator) evaluateCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	var fn object.Object

	if ident, ok := node.Function.(*ast.Identifier); ok {
		// If we got here, this means that we're calling a function in a variable
		fnName := ident.Value
		foundFn, exists := env.Get(fnName)

		if !exists {
			return e.throwErr(
				node.Function,
				"This error occurs when an undeclared function was called",
				"Cannot call function '%s', as it is undefined",
				fnName,
			)
		}

		fn = foundFn
	} else if fnLit, ok := node.Function.(*ast.FunctionLiteral); ok {
		args := e.evaluateExpressions(node.Arguments, env)

		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		fnLitObj := e.Evaluate(fnLit, env)
		fn = fnLitObj
	} else {
		return e.throwErr(
			node.Function,
			"This error occurs when trying to call an invalid expression as a function",
			"Unexpected call to an invalid expression",
		)
	}

	args := e.evaluateExpressions(node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	return e.applyFunction(node.Function, node, fn, args)
}

func (e *Evaluator) applyFunction(
	fnNode ast.Node,
	callNode *ast.CallExpression,
	function object.Object,
	args []object.Object,
) object.Object {
	switch fn := function.(type) {
	case *object.Function:
		if e.callDepth >= e.MaxCallDepth {
			return e.throwErr(
				fnNode,
				"This error occurs when a function calls itself (or other functions) too deeply",
				"Maximum call stack depth exceeded (%d calls)",
				e.callDepth,
			)
		}

		if len(fn.Parameters) != len(args) {
			return e.throwErr(
				callNode,
				"Argument count mismatch",
				"Expected %d arguments, got %d",
				len(fn.Parameters),
				len(args),
			)
		}

		e.callDepth++
		defer func() { e.callDepth-- }()

		extendedEnv := e.extendFunctionEnv(fn, args)
		evaluated := e.Evaluate(fn.Body, extendedEnv)
		return e.unwrapReturnValue(evaluated)

	case *object.NativeFunction:
		// Native functions still contributes to call stack depth
		// but doesnt have a limiter
		e.callDepth++
		defer func() { e.callDepth-- }()

		return fn.Fn(callNode, args)

	default:
		return e.throwErr(
			fnNode,
			"Cannot call non-function expression",
			"Attempted to call %s",
			function.Type(),
		)
	}
}

func (e *Evaluator) extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	// fn env is the outer env (for closure) ---
	env := object.NewEnvironment(fn.Scope)

	for idx, param := range fn.Parameters {
		// Assign args to params
		env.Set(param.Value, args[idx])
	}

	return env
}

func (e *Evaluator) evaluateArrayLiteral(node *ast.ArrayLiteral, env *object.Environment) object.Object {
	exprs := e.evaluateExpressions(node.Elements, env)
	return &object.Array{Elements: exprs}
}

func (e *Evaluator) evaluateIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	target := e.Evaluate(node.Target, env)
	if isError(target) {
		return target
	}

	index := e.Evaluate(node.Index, env)
	if isError(index) {
		return index
	}

	switch {
	case target.Type() == object.ARRAY_OBJECT && index.Type() == object.NUMBER_OBJECT:
		return e.evaluateArrayIndexExpression(node, target, index)
	case target.Type() == object.HASH_OBJECT:
		return e.evaluateHashIndexExpreesion(node, target, index)

	default:
		return e.throwErr(
			node,
			"This error occurs when trying to index an invalid expression",
			"Cannot index expression type '%s' with index type of '%s'",
			target.Type(),
			index.Type(),
		)
	}
}

func (e *Evaluator) evaluateHashIndexExpreesion(node *ast.IndexExpression, target, index object.Object) object.Object {
	hash := target.(*object.Hash).Pairs
	key, ok := index.(object.Hashable)

	if !ok {
		return e.throwErr(
			node.Index,
			"This error occurs when trying to use an invalid key for accessing a value in a hash",
			"Cannot use key type '%s' for accessing a hash",
			index.Type(),
		)
	}

	pair, ok := hash[key.HashKey()]
	if !ok {
		return object.NIL
	}

	return pair.Value
}

func (e *Evaluator) evaluateArrayIndexExpression(node *ast.IndexExpression, target object.Object, index object.Object) object.Object {
	t := target.(*object.Array).Elements
	i := int(index.(*object.Number).Value)
	maxLen := len(t) - 1

	if i < 0 || i > maxLen {
		return e.throwErr(
			node.Index,
			"This error occurs when trying to index an array smaller or bigger than its current length",
			"Array index '%d' out-of-bounds",
			i,
		)
	}

	return t[i]
}

func (e *Evaluator) evaluateIndexAssignmentExpression(node *ast.IndexAssignmentExpression, env *object.Environment) object.Object {
	target := e.Evaluate(node.Target, env)
	if isError(target) {
		return target
	}

	switch target.Type() {
	case object.ARRAY_OBJECT:
		return e.evaluateArrayIndexAssignmentExpression(node, target, env)
	case object.HASH_OBJECT:
		return e.evaluateHashIndexAssignmentExpression(node, target, env)

	default:
		return e.throwErr(
			node.Target,
			"This error occurs when trying to assign a non-indexable expression",
			"Cannot re-assign non-indexable expression type '%s'",
			target.Type(),
		)
	}
}

func (e *Evaluator) evaluateHashIndexAssignmentExpression(node *ast.IndexAssignmentExpression, target object.Object, env *object.Environment) object.Object {
	hash := target.(*object.Hash).Pairs
	index := e.Evaluate(node.Index, env)
	if isError(index) {
		return index
	}

	key, ok := index.(object.Hashable)
	if !ok {
		return e.throwErr(
			node.Index,
			"This error occurs when trying to use an invalid key for accessing a value in a hash",
			"Cannot use key type '%s' for accessing a hash",
			index.Type(),
		)
	}

	newValue := e.Evaluate(node.NewValue, env)

	hash[key.HashKey()] = object.HashPair{Key: index, Value: newValue}

	return newValue
}

func (e *Evaluator) evaluateArrayIndexAssignmentExpression(node *ast.IndexAssignmentExpression, target object.Object, env *object.Environment) object.Object {
	array := target.(*object.Array)
	i := e.Evaluate(node.Index, env)
	if isError(i) {
		return i
	}

	newValue := e.Evaluate(node.NewValue, env)
	if isError(newValue) {
		return newValue
	}

	index, ok := i.(*object.Number)
	if !ok {
		return e.throwErr(
			node.Index,
			"This error occurs when trying to index an array with a non-number index",
			"Cannot index an array with index type '%s'",
			i.Type(),
		)
	}

	array.Elements[int(index.Value)] = newValue
	return newValue
}

func (e *Evaluator) evaluateHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := e.Evaluate(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return e.throwErr(
				keyNode,
				"This error occurs when trying to use an unsupported value as a key",
				"Cannot access hash with key type '%s'",
				key.Type(),
			)
		}

		value := e.Evaluate(valueNode, env)
		if isError(value) {
			return value
		}

		hash := hashKey.HashKey()
		pairs[hash] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}
