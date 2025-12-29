package evaluation

import (
	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
	"github.com/sanity-io/litter"
)

type Evaluator struct {
	line   uint
	column uint
}

func New() Evaluator {
	return Evaluator{}
}

func (e *Evaluator) Evaluate(node ast.Node) object.Object {
	e.line = node.GetLine()
	e.column = node.GetColumn()

	switch node := node.(type) {
	case *ast.Program:
		return e.evaluateStatements(node.Statements)
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
		right := e.Evaluate(node.Right)
		return e.evaluateUnaryExpression(node.Operator.Type, right)
	case *ast.BinaryExpression:
		left := e.Evaluate(node.Left)
		right := e.Evaluate(node.Right)
		return e.evaluateBinaryExpression(node.Operator.Type, left, right)
	case *ast.TernaryExpression:
		condition := e.Evaluate(node.Condition)
		return e.evaluateTernaryExpression(condition, node)
	case *ast.ExpressionStatement:
		return e.Evaluate(node.Expression)
	case *ast.BlockStatement:
		return e.evaluateStatements(node.Statements)
	case *ast.IfStatement:
		condition := e.Evaluate(node.Condition)
		return e.evaluateIfStatement(condition, node)

	default:
		println("Unrecognized AST node:\n")
		litter.Dump(node)
		return NIL
	}
}

func (e *Evaluator) evaluateStatements(statements []ast.Statement) object.Object {
	var lastEval object.Object

	for _, stmt := range statements {
		lastEval = e.Evaluate(stmt)
	}

	return lastEval
}
