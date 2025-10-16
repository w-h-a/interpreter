package parser

import (
	"github.com/w-h-a/interpreter/internal/ast/expression"
	"github.com/w-h-a/interpreter/internal/ast/expression/call"
	infixoperator "github.com/w-h-a/interpreter/internal/ast/expression/infix_operator"
	prefixoperator "github.com/w-h-a/interpreter/internal/ast/expression/prefix_operator"
)

type (
	parsePrefixFn func(*Parser) (expression.Expression, error)
	parseInfixFn  func(*Parser, expression.Expression) (expression.Expression, error)
)

func parsePrefixOperatorExpression(p *Parser) (expression.Expression, error) {
	expression := &prefixoperator.PrefixOperator{
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

func parseInfixOperatorExpression(p *Parser, left expression.Expression) (expression.Expression, error) {
	expression := &infixoperator.InfixOperator{
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

func parseCallExpression(p *Parser, function expression.Expression) (expression.Expression, error) {
	exp := &call.Call{Token: p.curToken, Function: function}

	var err error

	exp.Arguments, err = p.parseCallArguments()
	if err != nil {
		return nil, err
	}

	return exp, nil
}
