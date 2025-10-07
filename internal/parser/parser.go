package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression/boolean"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression/identifier"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression/integer"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement"
	expressionstatement "github.com/w-h-a/interpreter/internal/parser/ast/statement/expression"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement/let"
	returnstatement "github.com/w-h-a/interpreter/internal/parser/ast/statement/return"
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
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
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

	p.nextToken()

	// TODO: RHS expression

	for p.curToken.Type != token.Semicolon {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseReturnStatement() (*returnstatement.Return, error) {
	stmt := &returnstatement.Return{Token: p.curToken}

	p.nextToken()

	// TODO: expression

	for p.curToken.Type != token.Semicolon {
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
		exp, err = p.parseIdentifier()
	case token.ParenLeft:
		exp, err = p.parseGrouped()
	case token.Int:
		exp, err = p.parseInteger()
	case token.True, token.False:
		exp, err = p.parseBoolean()
	default:
		parsePrefixExpression := p.parsePrefixFns[p.curToken.Type]

		if parsePrefixExpression == nil {
			errDetail := fmt.Sprintf("no parse prefix function for %s found", p.curToken.Type)
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

func (p *Parser) parseIdentifier() (expression.Expression, error) {
	return &identifier.Identifier{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func (p *Parser) parseGrouped() (expression.Expression, error) {
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

func (p *Parser) parseInteger() (expression.Expression, error) {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		return nil, err
	}

	return &integer.Integer{Token: p.curToken, Value: value}, nil
}

func (p *Parser) parseBoolean() (expression.Expression, error) {
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

	p.registerParsePrefixFn(token.Bang, parsePrefixOperator)
	p.registerParsePrefixFn(token.Minus, parsePrefixOperator)

	p.registerParseInfixFn(token.Identical, parseInfixOperator)
	p.registerParseInfixFn(token.NotIdentical, parseInfixOperator)
	p.registerParseInfixFn(token.LessThan, parseInfixOperator)
	p.registerParseInfixFn(token.GreaterThan, parseInfixOperator)
	p.registerParseInfixFn(token.Plus, parseInfixOperator)
	p.registerParseInfixFn(token.Minus, parseInfixOperator)
	p.registerParseInfixFn(token.Asterisk, parseInfixOperator)
	p.registerParseInfixFn(token.Slash, parseInfixOperator)

	p.nextToken()
	p.nextToken()

	return p
}
