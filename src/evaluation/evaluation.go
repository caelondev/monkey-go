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
		return e.evaluateProgram(node.Statements)
	case *ast.ExpressionStatement:
		return e.evaluateExpression(node.Expression)

	default:
		println("Unrecognized AST node:\n")
		litter.Dump(node)
		return NIL
	}
}

func (e *Evaluator) evaluateProgram(statements []ast.Statement) object.Object {
	var lastEval object.Object

	for _, stmt := range statements {
		lastEval = e.Evaluate(stmt)
	}

	return lastEval
}
