package parser

import (
	"fmt"

	"github.com/w-h-a/interpreter/internal/parser/ast"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement"
	"github.com/w-h-a/interpreter/internal/token"
)

type Parser struct {
	tokens    chan token.Token
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
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
		return nil
	}
}

func (p *Parser) parseLetStatement() *statement.Let {
	stmt := &statement.Let{Token: p.curToken}

	if p.peekToken.Type != token.Ident {
		p.appendError(token.Ident)
		return nil
	}

	p.nextToken()

	stmt.Name = &expression.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type != token.Assign {
		p.appendError(token.Assign)
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

func (p *Parser) appendError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = <-p.tokens
}

func New(tks chan token.Token) *Parser {
	p := &Parser{
		tokens: tks,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}
