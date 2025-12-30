// evaluator.go
package evaluation

import (
	"fmt"
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

func (e *Evaluator) Evaluate(node ast.Node) object.Object {
	e.line = node.GetLine()
	e.column = node.GetColumn()

	switch node := node.(type) {
	case *ast.Program:
		return e.evaluateProgram(node.Statements)
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
		if isError(right) {
			return right
		}
		return e.evaluateUnaryExpression(node, right)
	case *ast.BinaryExpression:
		left := e.Evaluate(node.Left)
		if isError(left) {
			return left
		}
		right := e.Evaluate(node.Right)
		if isError(right) {
			return right
		}
		return e.evaluateBinaryExpression(node, left, right)
	case *ast.TernaryExpression:
		condition := e.Evaluate(node.Condition)
		if isError(condition) {
			return condition
		}
		return e.evaluateTernaryExpression(node, condition)
	case *ast.ExpressionStatement:
		return e.Evaluate(node.Expression)
	case *ast.BlockStatement:
		return e.evaluateBlockStatement(node)
	case *ast.IfStatement:
		condition := e.Evaluate(node.Condition)
		return e.evaluateIfStatement(node, condition)
	case *ast.ReturnStatement:
		if node.ReturnValue == nil {
			return &object.ReturnValue{Value: NIL}
		}
		value := e.Evaluate(node.ReturnValue)
		return &object.ReturnValue{Value: value}

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

		switch result := lastEval.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return lastEval
}

func (e *Evaluator) throwErr(node ast.Node, hint string, format string, a ...interface{}) *object.Error {
	lineColumn := fmt.Sprintf("[Ln %d:%d] Runtime::Error -> ", e.line, e.column)
	message := fmt.Sprintf(format, a...)

	snippet := "\n\n"

	// Show the actual source line if available
	if int(e.line) > 0 && int(e.line) <= len(e.lines) {
		sourceLine := e.lines[e.line-1]
		lineNumStr := fmt.Sprintf("Ln %d:%d", e.line, e.column)

		snippet += " Error caused by:\n"
		snippet += fmt.Sprintf("    %s | %s\n", lineNumStr, sourceLine)

		// Create pointer to error location
		padding := strings.Repeat(" ", len(lineNumStr))
		pointer := strings.Repeat(" ", int(e.column)-1) + "^"
		snippet += fmt.Sprintf("    %s | %s\n", padding, pointer)
	} else {
		// Fallback if source line isn't available
		snippet += " Error caused by:\n"
		snippet += fmt.Sprintf("\t%d:%d | %s\n", e.line, e.column, node.String())
	}

	if hint != "" {
		snippet += fmt.Sprintf("\n Hint: %s\n", hint)
	}

	return &object.Error{Message: lineColumn + message + snippet}
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJECT
}
