package parser

import (
	"strconv"

	"github.com/w-h-a/interpreter/internal/parser/ast"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
)

type (
	parsePrefixExpression func(*Parser) (ast.Expression, error)
	parseInfixExpression  func(*Parser, ast.Expression) (ast.Expression, error)
)

func parseIdentifier(p *Parser) (ast.Expression, error) {
	return &expression.Identifier{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func parseInteger(p *Parser) (ast.Expression, error) {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		return nil, err
	}

	return &expression.Integer{Token: p.curToken, Value: value}, nil
}
