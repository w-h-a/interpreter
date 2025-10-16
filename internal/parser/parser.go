package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/w-h-a/interpreter/internal/ast/expression"
	"github.com/w-h-a/interpreter/internal/ast/expression/boolean"
	"github.com/w-h-a/interpreter/internal/ast/expression/function"
	"github.com/w-h-a/interpreter/internal/ast/expression/identifier"
	ifexpression "github.com/w-h-a/interpreter/internal/ast/expression/if"
	"github.com/w-h-a/interpreter/internal/ast/expression/integer"
	"github.com/w-h-a/interpreter/internal/ast/statement"
	"github.com/w-h-a/interpreter/internal/ast/statement/block"
	expressionstatement "github.com/w-h-a/interpreter/internal/ast/statement/expression"
	"github.com/w-h-a/interpreter/internal/ast/statement/let"
	returnstatement "github.com/w-h-a/interpreter/internal/ast/statement/return"
	"github.com/w-h-a/interpreter/internal/token"
)

type Parser struct {
	tokens         chan token.Token
	curToken       token.Token
	peekToken      token.Token
	parsePrefixFns map[token.TokenType]parsePrefixFn
	parseInfixFns  map[token.TokenType]parseInfixFn
	errors         []string
}

func (p *Parser) ParseProgram() *statement.Program {
	program := &statement.Program{}
	program.Statements = []statement.Statement{}

	for p.curToken.Type != token.EOF {
		stmt, err := p.parseStatement()
		if err == nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() (statement.Statement, error) {
	switch p.curToken.Type {
	case token.Return:
		return p.parseReturnStatement()
	case token.Let:
		return p.parseLetStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() (*returnstatement.Return, error) {
	stmt := &returnstatement.Return{Token: p.curToken}

	p.nextToken()

	var err error

	stmt.Value, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekToken.Type == token.Semicolon {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseLetStatement() (*let.Let, error) {
	stmt := &let.Let{Token: p.curToken}

	if p.peekToken.Type != token.Ident {
		errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.Ident, p.peekToken.Type)
		p.appendError(errDetail)
		return nil, errors.New(errDetail)
	}

	p.nextToken()

	stmt.Name = &identifier.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type != token.Assign {
		errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.Assign, p.peekToken.Type)
		p.appendError(errDetail)
		return nil, errors.New(errDetail)
	}

	p.nextToken() // consume identifier
	p.nextToken() // consume assignment

	var err error

	stmt.Value, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekToken.Type == token.Semicolon {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseBlockStatement() (*block.Block, error) {
	stmt := &block.Block{Token: p.curToken}

	stmt.Statements = []statement.Statement{}

	p.nextToken() // consume '{'

	for p.curToken.Type != token.BraceRight && p.curToken.Type != token.EOF {
		s, err := p.parseStatement()
		if err == nil {
			stmt.Statements = append(stmt.Statements, s)
		}
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseExpressionStatement() (*expressionstatement.Expression, error) {
	stmt := &expressionstatement.Expression{Token: p.curToken}

	var err error

	stmt.Expression, err = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.Semicolon {
		p.nextToken()
	}

	return stmt, err
}

func (p *Parser) parseExpression(precedence int) (expression.Expression, error) {
	var exp expression.Expression
	var err error

	switch p.curToken.Type {
	case token.Ident:
		exp, err = p.parseIdentifierExpression()
	case token.ParenLeft:
		exp, err = p.parseGroupedExpression()
	case token.If:
		exp, err = p.parseIfExpression()
	case token.Function:
		exp, err = p.parseFunctionExpression()
	case token.Int:
		exp, err = p.parseIntegerExpression()
	case token.True, token.False:
		exp, err = p.parseBooleanExpression()
	default:
		parsePrefixExpression := p.parsePrefixFns[p.curToken.Type]

		if parsePrefixExpression == nil {
			errDetail := fmt.Sprintf("no parse function for %s found", p.curToken.Type)
			p.appendError(errDetail)
			return nil, errors.New(errDetail)
		}

		exp, err = parsePrefixExpression(p)
	}

	if err != nil {
		errDetail := fmt.Sprintf("failed to parse expression literal %q", p.curToken.Literal)
		p.appendError(fmt.Sprintf("%s: %v", errDetail, err))
		return nil, fmt.Errorf("%s: %w", errDetail, err)
	}

	for p.peekToken.Type != token.Semicolon && precedence < p.peekPrecedence() {
		parseInfixExpression := p.parseInfixFns[p.peekToken.Type]

		if parseInfixExpression == nil {
			return exp, nil
		}

		p.nextToken()

		exp, err = parseInfixExpression(p, exp)
		if err != nil {
			errDetail := fmt.Sprintf("failed to parse infix expression literal %q", p.curToken.Literal)
			p.appendError(fmt.Sprintf("%s: %v", errDetail, err))
			return nil, fmt.Errorf("%s: %w", errDetail, err)
		}
	}

	return exp, nil
}

func (p *Parser) parseIdentifierExpression() (expression.Expression, error) {
	return &identifier.Identifier{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func (p *Parser) parseGroupedExpression() (expression.Expression, error) {
	p.nextToken()

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekToken.Type != token.ParenRight {
		errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.ParenRight, p.peekToken.Type)
		return nil, errors.New(errDetail)
	}

	p.nextToken()

	return exp, nil
}

func (p *Parser) parseIfExpression() (expression.Expression, error) {
	exp := &ifexpression.If{Token: p.curToken}

	if p.peekToken.Type != token.ParenLeft {
		errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.ParenLeft, p.peekToken.Type)
		return nil, errors.New(errDetail)
	}

	p.nextToken() // move to '('

	p.nextToken() // consume '(' to get ready to parse condition

	var err error

	exp.Condition, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekToken.Type != token.ParenRight {
		errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.ParenRight, p.peekToken.Type)
		return nil, errors.New(errDetail)
	}

	p.nextToken() // move to ')'

	if p.peekToken.Type != token.BraceLeft {
		errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.BraceLeft, p.peekToken.Type)
		return nil, errors.New(errDetail)
	}

	p.nextToken() // move to '{'

	exp.Consequence, _ = p.parseBlockStatement()

	if p.peekToken.Type == token.Else {
		p.nextToken() // move to 'else'

		if p.peekToken.Type != token.BraceLeft {
			errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.BraceLeft, p.peekToken.Type)
			return nil, errors.New(errDetail)
		}

		p.nextToken() // move to '{'

		exp.Alternative, _ = p.parseBlockStatement()
	}

	return exp, nil
}

func (p *Parser) parseCallArguments() ([]expression.Expression, error) {
	args := []expression.Expression{}

	if p.peekToken.Type == token.ParenRight {
		p.nextToken()
		return args, nil
	}

	p.nextToken() // consume '('

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	args = append(args, exp)

	for p.peekToken.Type == token.Comma {
		p.nextToken() // consume previous arg
		p.nextToken() // consume ','
		exp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, exp)
	}

	if p.peekToken.Type != token.ParenRight {
		errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.ParenRight, p.peekToken.Type)
		return nil, errors.New(errDetail)
	}

	p.nextToken() // consume last arg

	return args, nil
}

func (p *Parser) parseFunctionExpression() (expression.Expression, error) {
	exp := &function.Function{Token: p.curToken}

	if p.peekToken.Type != token.ParenLeft {
		errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.ParenLeft, p.peekToken.Type)
		return nil, errors.New(errDetail)
	}

	p.nextToken() // consume 'fn'

	var err error

	exp.Parameters, err = p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	if p.peekToken.Type != token.BraceLeft {
		errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.BraceLeft, p.peekToken.Type)
		return nil, errors.New(errDetail)
	}

	p.nextToken() // consume ')'

	exp.Body, _ = p.parseBlockStatement()

	return exp, nil
}

func (p *Parser) parseFunctionParameters() ([]*identifier.Identifier, error) {
	identifiers := []*identifier.Identifier{}

	if p.peekToken.Type == token.ParenRight {
		p.nextToken()
		return identifiers, nil
	}

	p.nextToken() // consume '('

	if p.curToken.Type != token.Ident {
		errDetail := fmt.Sprintf("expected identifier as function parameter, got %s", p.curToken.Type)
		return nil, errors.New(errDetail)
	}

	ident := &identifier.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekToken.Type == token.Comma {
		p.nextToken() // consume previous param
		p.nextToken() // consume ','
		if p.curToken.Type != token.Ident {
			errDetail := fmt.Sprintf("expected identifier as function parameter, got %s", p.curToken.Type)
			return nil, errors.New(errDetail)
		}
		ident := &identifier.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if p.peekToken.Type != token.ParenRight {
		errDetail := fmt.Sprintf("expected next token to be %s, got %s", token.ParenRight, p.peekToken.Type)
		return nil, errors.New(errDetail)
	}

	p.nextToken() // consume last param

	return identifiers, nil
}

func (p *Parser) parseIntegerExpression() (expression.Expression, error) {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		return nil, err
	}

	return &integer.Integer{Token: p.curToken, Value: value}, nil
}

func (p *Parser) parseBooleanExpression() (expression.Expression, error) {
	return &boolean.Boolean{Token: p.curToken, Value: p.curToken.Type == token.True}, nil
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) currentPrecendence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) appendError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = <-p.tokens
}

func (p *Parser) registerParsePrefixFn(tokenType token.TokenType, fn parsePrefixFn) {
	p.parsePrefixFns[tokenType] = fn
}

func (p *Parser) registerParseInfixFn(tokenType token.TokenType, fn parseInfixFn) {
	p.parseInfixFns[tokenType] = fn
}

func New(tks chan token.Token) *Parser {
	p := &Parser{
		tokens:         tks,
		errors:         []string{},
		parsePrefixFns: map[token.TokenType]parsePrefixFn{},
		parseInfixFns:  map[token.TokenType]parseInfixFn{},
	}

	p.registerParsePrefixFn(token.Bang, parsePrefixOperatorExpression)
	p.registerParsePrefixFn(token.Minus, parsePrefixOperatorExpression)

	p.registerParseInfixFn(token.Identical, parseInfixOperatorExpression)
	p.registerParseInfixFn(token.NotIdentical, parseInfixOperatorExpression)
	p.registerParseInfixFn(token.LessThan, parseInfixOperatorExpression)
	p.registerParseInfixFn(token.GreaterThan, parseInfixOperatorExpression)
	p.registerParseInfixFn(token.Plus, parseInfixOperatorExpression)
	p.registerParseInfixFn(token.Minus, parseInfixOperatorExpression)
	p.registerParseInfixFn(token.Asterisk, parseInfixOperatorExpression)
	p.registerParseInfixFn(token.Slash, parseInfixOperatorExpression)
	p.registerParseInfixFn(token.ParenLeft, parseCallExpression)

	p.nextToken()
	p.nextToken()

	return p
}
