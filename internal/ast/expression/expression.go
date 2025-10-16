package expression

import "github.com/w-h-a/interpreter/internal/ast"

type Expression interface {
	ast.Node
	ExpressionNode()
}
