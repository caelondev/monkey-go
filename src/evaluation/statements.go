package evaluation

import (
	"github.com/caelondev/monkey/src/ast"
	"github.com/caelondev/monkey/src/object"
)

func (e *Evaluator) evaluateBlockStatement(node *ast.BlockStatement, env *object.Environment) object.Object {
	var lastEvaluated object.Object

	for _, stmt := range node.Statements {
		lastEvaluated = e.Evaluate(stmt, env)

		if lastEvaluated != nil {
			if lastEvaluated.Type() == object.RETURN_VALUE_OBJECT || lastEvaluated.Type() == object.ERROR_OBJECT {
				return lastEvaluated
			}
		}
	}

	return lastEvaluated
}

func (e *Evaluator) evaluateIfStatement(node *ast.IfStatement, env *object.Environment) object.Object {
	condition := e.Evaluate(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Evaluate(node.Consequence, env)
	} else {
		if node.Alternative == nil {
			return object.NIL
		}

		return e.Evaluate(node.Alternative, env)
	}
}

func (e *Evaluator) evaluateVariableDeclaration(node *ast.VarStatement, env *object.Environment) object.Object {
	// Check if every assignees are valid ---
	// Then discard everything if not ---
	for _, name := range node.Names {
		if !env.DoesExist(name.Value) {
			continue
		}

		return e.throwErr(
			node.Value,
			"This error occurs when a variable that is already declared was redeclared again in the same scope",
			"Cannot declare '%s' as it already exists",
			name.Value,
		)
	}

	value := e.Evaluate(node.Value, env)
	if isError(value) {
		return value
	}

	for _, name := range node.Names {
		env.Set(name.Value, value)
	}

	return value
}

func (e *Evaluator) evaluateBatchAssignmentStatement(node *ast.BatchAssignmentStatement, env *object.Environment) object.Object {
	// Check if every assignees are valid ---
	// Then discard everything if not ---
	for _, assignee := range node.Assignees {
		exists := env.DoesExist(assignee.Value)

		if !exists {
			return e.throwErr(
				assignee,
				"This error occurs when the assignee variable doesnt exist",
				"Cannot resolve variable '%s'",
				assignee.Value,
			)
		}
	}

	newValue := e.Evaluate(node.NewValue, env)

	for _, assignee := range node.Assignees {
		env.Set(assignee.Value, newValue)
	}

	return newValue
}

func (e *Evaluator) evaluateReturnStatement(node *ast.ReturnStatement, env *object.Environment) object.Object {
	if node.ReturnValue == nil {
		return &object.ReturnValue{Value: object.NIL}
	}

	value := e.Evaluate(node.ReturnValue, env)
	return &object.ReturnValue{Value: value}
}

func (e *Evaluator) evaluateFunctionDeclaration(node *ast.FunctionDeclarationStatement, env *object.Environment) object.Object {
	function := &object.Function{
		Parameters: node.Parameters,
		Name:       node.Name,
		Body:       node.Body,
		Scope:      env,
	}

	value := env.Declare(node.Name.Value, function)
	return value
}
