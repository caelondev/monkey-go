// evaluator.go
package evaluation

import (
	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
)

type Evaluator struct {
	line         uint
	column       uint
	callDepth    int
	MaxCallDepth int
}

func New() Evaluator {
	return Evaluator{
		MaxCallDepth: 10_000,
	}
}

func (e *Evaluator) Evaluate(node ast.Node, env *object.Environment) object.Object {
	if env.GetOuter() == nil { // Global env
		e.InitializeNativeFunctions(env)
	}

	e.line = node.GetLine()
	e.column = node.GetColumn()

	switch node := node.(type) {
	case *ast.Program:
		return e.evaluateProgram(node.Statements, env)
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.NilLiteral:
		return object.NIL
	case *ast.NaNLiteral:
		return object.NAN
	case *ast.InfinityLiteral:
		return object.INFINITY
	case *ast.BooleanExpression:
		return e.evaluateToObjectBoolean(node.Value)
	case *ast.UnaryExpression:
		return e.evaluateUnaryExpression(node, env)
	case *ast.BinaryExpression:
		return e.evaluateBinaryExpression(node, env)
	case *ast.TernaryExpression:
		return e.evaluateTernaryExpression(node, env)
	case *ast.ExpressionStatement:
		return e.Evaluate(node.Expression, env)
	case *ast.BlockStatement:
		return e.evaluateBlockStatement(node, env)
	case *ast.IfStatement:
		return e.evaluateIfStatement(node, env)
	case *ast.ReturnStatement:
		return e.evaluateReturnStatement(node, env)
	case *ast.VarStatement:
		return e.evaluateVariableDeclaration(node, env)
	case *ast.Identifier:
		return e.evaluateIdentifier(node, env)
	case *ast.AssignmentExpression:
		return e.evaluateAssignmentExpression(node, env)
	case *ast.BatchAssignmentStatement:
		return e.evaluateBatchAssignmentStatement(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Scope: env}
	case *ast.FunctionDeclarationStatement:
		return e.evaluateFunctionDeclaration(node, env)
	case *ast.CallExpression:
		return e.evaluateCallExpression(node, env)
	case *ast.ArrayLiteral:
		return e.evaluateArrayLiteral(node, env)
	case *ast.IndexExpression:
		return e.evaluateIndexExpression(node, env)
	case *ast.IndexAssignmentExpression:
		return e.evaluateIndexAssignmentExpression(node, env)
	case *ast.HashLiteral:
		return e.evaluateHashLiteral(node, env)

	default:
		return e.throwErr(
			node,
			"This error occurs when an unhandled AST was passed.\nThis error should only happen in language development",
			"Unrecognized Abstract Syntax Tree node:\n%v",
			node,
		)
	}
}

func (e *Evaluator) evaluateProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var lastEval object.Object

	for _, stmt := range statements {
		lastEval = e.Evaluate(stmt, env)

		switch result := lastEval.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return lastEval
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJECT
}

func (e *Evaluator) InitializeNativeFunctions(env *object.Environment) {
	e.registerNativeFn(env, "len", e.NATIVE_LEN_FUNCTION)
	e.registerNativeFn(env, "print", e.NATIVE_PRINT_FUNCTION)
	e.registerNativeFn(env, "prompt", e.NATIVE_PROMPT_FUNCTION)
	e.registerNativeFn(env, "time", e.NATIVE_TIME_FUNCTION)
	e.registerNativeFn(env, "to_string", e.NATIVE_TO_STRING_FUNCTION)
	e.registerNativeFn(env, "to_number", e.NATIVE_TO_NUMBER_FUNCTION)
	e.registerNativeFn(env, "type", e.NATIVE_TYPE_FUNCTION)
	e.registerNativeFn(env, "is_NaN", e.NATIVE_IS_NAN_FUNCTION)
	e.registerNativeFn(env, "is_Inf", e.NATIVE_IS_INF_FUNCTION)
	e.registerNativeFn(env, "is_nil", e.NATIVE_IS_NIL_FUNCTION)
	e.registerNativeFn(env, "random", e.NATIVE_RANDOM_FUNCTION)
}

func (e *Evaluator) registerNativeFn(env *object.Environment, name string, fn object.NativeFunctionFn) {
	fnObject := &object.NativeFunction{Fn: fn}
	env.Declare(name, fnObject)
}
