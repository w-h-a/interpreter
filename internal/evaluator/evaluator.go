package evaluator

import (
	"github.com/w-h-a/interpreter/internal/ast"
	intexp "github.com/w-h-a/interpreter/internal/ast/expression/integer"
	"github.com/w-h-a/interpreter/internal/ast/statement"
	expressionstatement "github.com/w-h-a/interpreter/internal/ast/statement/expression"
	"github.com/w-h-a/interpreter/internal/object"
	intobj "github.com/w-h-a/interpreter/internal/object/integer"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *statement.Program:
		return evalStatements(node.Statements)
	case *expressionstatement.Expression:
		return Eval(node.Expression)
	case *intexp.Integer:
		return &intobj.Integer{Value: node.Value}
	default:
		return nil
	}
}

func evalStatements(stmts []statement.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}
