// evaluator.go
package evaluation

import (
	"strings"

	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
	"github.com/sanity-io/litter"
)

type Evaluator struct {
	line   uint
	column uint
	lines  []string
}

func New(source string) Evaluator {
	return Evaluator{
		lines: strings.Split(source, "\n"),
	}
}

func (e *Evaluator) Evaluate(node ast.Node, env *object.Environment) object.Object {
	e.line = node.GetLine()
	e.column = node.GetColumn()

	switch node := node.(type) {
	case *ast.Program:
		return e.evaluateProgram(node.Statements, env)
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.NilLiteral:
		return NIL
	case *ast.NaNLiteral:
		return NAN
	case *ast.InfinityLiteral:
		return INFINITY
	case *ast.BooleanExpression:
		return e.evaluateToObjectBoolean(node.Value)
	case *ast.UnaryExpression:
		right := e.Evaluate(node.Right, env)
		if isError(right) {
			return right
		}
		return e.evaluateUnaryExpression(node, right)

	case *ast.BinaryExpression:
		left := e.Evaluate(node.Left, env)
		if isError(left) {
			return left
		}
		right := e.Evaluate(node.Right, env)
		if isError(right) {
			return right
		}
		return e.evaluateBinaryExpression(node, left, right)
	case *ast.TernaryExpression:
		condition := e.Evaluate(node.Condition, env)
		if isError(condition) {
			return condition
		}
		return e.evaluateTernaryExpression(node, condition, env)
	case *ast.ExpressionStatement:
		return e.Evaluate(node.Expression, env)
	case *ast.BlockStatement:
		return e.evaluateBlockStatement(node, env)
	case *ast.IfStatement:
		condition := e.Evaluate(node.Condition, env)
		return e.evaluateIfStatement(node, condition, env)
	case *ast.ReturnStatement:
		if node.ReturnValue == nil {
			return &object.ReturnValue{Value: NIL}
		}
		value := e.Evaluate(node.ReturnValue, env)
		return &object.ReturnValue{Value: value}
	case *ast.VarStatement:
		return e.evaluateVariableDeclaration(node, env)
	case *ast.Identifier:
		return e.evaluateIdentifier(node, env)
	case *ast.AssignmentExpression:
		return e.evaluateAssignmentExpression(node, env)
	case *ast.BatchAssignmentStatement:
		return e.evaluateBatchAssignmentStatement(node, env)

	default:
		println("Unrecognized AST node:\n")
		litter.Dump(node)
		return NIL
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
