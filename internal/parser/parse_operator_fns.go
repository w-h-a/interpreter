package parser

import (
	"github.com/w-h-a/interpreter/internal/parser/ast"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
)

type (
	parsePrefixFn func(*Parser) (ast.Expression, error)
	parseInfixFn  func(*Parser, ast.Expression) (ast.Expression, error)
)

func parsePrefixOperator(p *Parser) (ast.Expression, error) {
	expression := &expression.PrefixOperator{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	var err error

	expression.Right, err = p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func parseInfixOperator(p *Parser, left ast.Expression) (ast.Expression, error) {
	expression := &expression.InfixOperator{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecendence()

	p.nextToken()

	var err error

	expression.Right, err = p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	return expression, nil
}
