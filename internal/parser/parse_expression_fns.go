package parser

import (
	"github.com/w-h-a/interpreter/internal/parser/ast"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
)

type (
	parsePrefixExpression func(*Parser) ast.Expression
	parseInfixExpression  func(*Parser, ast.Expression) ast.Expression
)

func parseIdentifier(p *Parser) ast.Expression {
	return &expression.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}
