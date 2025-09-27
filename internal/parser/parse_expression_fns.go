package parser

import "github.com/w-h-a/interpreter/internal/parser/ast"

type (
	parsePrefixExpression func() ast.Expression
	parseInfixExpression  func(ast.Expression) ast.Expression
)
