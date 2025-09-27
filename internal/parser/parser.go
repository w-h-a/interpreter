package parser

import (
	"fmt"

	"github.com/w-h-a/interpreter/internal/parser/ast"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement"
	"github.com/w-h-a/interpreter/internal/token"
)

type Parser struct {
	tokens         chan token.Token
	curToken       token.Token
	peekToken      token.Token
	parsePrefixFns map[token.TokenType]parsePrefixExpression
	parseInfixFns  map[token.TokenType]parseInfixExpression
	errors         []string
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		// this is required
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *statement.Let {
	stmt := &statement.Let{Token: p.curToken}

	if p.peekToken.Type != token.Ident {
		p.appendError(fmt.Sprintf("expected next token to be %s, got %s", token.Ident, p.peekToken.Type))
		return nil
	}

	p.nextToken()

	stmt.Name = &expression.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type != token.Assign {
		p.appendError(fmt.Sprintf("expected next token to be %s, got %s", token.Assign, p.peekToken.Type))
		return nil
	}

	p.nextToken()

	// TODO: RHS expression

	for p.curToken.Type != token.Semicolon {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *statement.Return {
	stmt := &statement.Return{Token: p.curToken}

	p.nextToken()

	// TODO: expression

	for p.curToken.Type != token.Semicolon {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *statement.Expression {
	stmt := &statement.Expression{Token: p.curToken}

	stmt.Expression = p.parseExpression()

	if p.peekToken.Type == token.Semicolon {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	parsePrefixExpression := p.parsePrefixFns[p.curToken.Type]

	if parsePrefixExpression == nil {
		p.appendError(fmt.Sprintf("no parse prefix function for %s found", p.curToken.Type))
		return nil
	}

	leftExp := parsePrefixExpression(p)

	return leftExp
}

func (p *Parser) appendError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = <-p.tokens
}

func (p *Parser) registerParsePrefixFn(tokenType token.TokenType, fn parsePrefixExpression) {
	p.parsePrefixFns[tokenType] = fn
}

func (p *Parser) registerParseInfixFn(tokenType token.TokenType, fn parseInfixExpression) {
	p.parseInfixFns[tokenType] = fn
}

func New(tks chan token.Token) *Parser {
	p := &Parser{
		tokens:         tks,
		errors:         []string{},
		parsePrefixFns: map[token.TokenType]parsePrefixExpression{},
		parseInfixFns:  map[token.TokenType]parseInfixExpression{},
	}

	p.registerParsePrefixFn(token.Ident, parseIdentifier)

	p.nextToken()
	p.nextToken()

	return p
}
