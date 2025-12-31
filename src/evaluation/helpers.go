package evaluation

import (
	"fmt"
	"math"
	"strings"

	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
	"github.com/caelondev/monkey/src/token"
	"github.com/jwalton/gchalk"
)

func (e *Evaluator) evaluateToObjectBoolean(v bool) *object.Boolean {
	if v {
		return TRUE
	}
	return FALSE
}

func (e *Evaluator) evaluateNotExpression(right object.Object) object.Object {
	if !isTruthy(right) {
		return TRUE
	}
	return FALSE
}

func (e *Evaluator) evaluateBinaryExpression(
	node *ast.BinaryExpression,
	left, right object.Object,
) object.Object {
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

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Nil:
		return false

	case *object.Boolean:
		return obj.Value

	case *object.Number:
		return obj.Value != 0

	case *object.NaN:
		return false

	default:
		return true
	}
}

func (e *Evaluator) throwErr(node ast.Node, hint string, format string, a ...interface{}) *object.Error {
	var line, column int
	switch n := node.(type) {
	case *ast.Identifier:
		line = int(n.Token.Line)
		column = int(n.Token.Column)
	default:
		line = int(e.line)
		column = int(e.column)
	}

	// Styled header with bold red
	lineColumn := gchalk.WithBold().Red(fmt.Sprintf("[Ln %d:%d] Runtime::Error", line, column))
	message := gchalk.Red(" -> " + fmt.Sprintf(format, a...))

	snippet := "\n\n"

	// Show the actual source line if available
	if line > 0 && line <= len(e.lines) {
		sourceLine := e.lines[line-1]
		lineNumStr := fmt.Sprintf("Ln %d:%d", line, column)

		snippet += gchalk.WithBold().White(" Error caused by:\n")
		
		// Line number in light blue/cyan
		snippet += gchalk.Cyan(fmt.Sprintf("    %s | ", lineNumStr))
		snippet += gchalk.White(sourceLine + "\n")

		// Pointer in bright red
		padding := strings.Repeat(" ", len(lineNumStr))
		pointer := strings.Repeat(" ", column-1) + "^"
		snippet += gchalk.Cyan(fmt.Sprintf("    %s | ", padding))
		snippet += gchalk.BrightRed(pointer + "\n")
	} else {
		snippet += gchalk.WithBold().White(" Error caused by:\n")
		snippet += gchalk.Cyan(fmt.Sprintf("\t%d:%d | ", line, column))
		snippet += gchalk.White(node.String() + "\n")
	}

	if hint != "" {
		snippet += "\n"
		snippet += gchalk.Cyan(" Hint: ")
		snippet += gchalk.White(hint + "\n")
	}

	return &object.Error{Message: lineColumn + message + snippet}
}
