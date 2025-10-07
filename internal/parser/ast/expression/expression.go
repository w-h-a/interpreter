package expression

import "github.com/w-h-a/interpreter/internal/parser/ast"

type Expression interface {
	ast.Node
	ExpressionNode()
}
