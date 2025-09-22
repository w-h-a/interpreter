package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/parser"
	"github.com/w-h-a/interpreter/internal/parser/ast"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement"
)

type want struct {
	expectedIdentifier string
}

type wantError struct {
	expectedError string
}

func TestLetStatements(t *testing.T) {
	happyTestCases := []struct {
		input string
		wants []want
	}{
		{

			input: `
let x = 5;
let y = 10;
let foobar = 838383;
`,
			wants: []want{
				{expectedIdentifier: "x"},
				{expectedIdentifier: "y"},
				{expectedIdentifier: "foobar"},
			},
		},
	}

	for _, tc := range happyTestCases {
		tks := lexer.Lex(tc.input)
		p := parser.New(tks)

		program := p.ParseProgram()
		errors := p.Errors()

		require.Equal(t, len(tc.wants), len(program.Statements))
		require.Equal(t, 0, len(errors))

		testLetStatement(t, tc.wants, program.Statements)
	}
}

func testLetStatement(t *testing.T, wants []want, stmts []ast.Statement) {
	for i, want := range wants {
		s := stmts[i]
		require.Equal(t, "let", s.TokenLiteral())

		letStmt, ok := s.(*statement.Let)
		require.True(t, ok)

		identExp := letStmt.Name

		require.Equal(t, want.expectedIdentifier, identExp.TokenLiteral())
		require.Equal(t, want.expectedIdentifier, identExp.Value)
	}
}

func TestLetStatements_Errors(t *testing.T) {
	errorTestCases := []struct {
		input string
		wants []wantError
	}{
		{
			input: `
let x 5;
let = 10;
let 838383;
`,
			wants: []wantError{
				{"expected next token to be =, got INT"},
				{"expected next token to be IDENT, got ="},
				{"expected next token to be IDENT, got INT"},
			},
		},
	}

	for _, tc := range errorTestCases {
		tks := lexer.Lex(tc.input)
		p := parser.New(tks)

		p.ParseProgram()
		errors := p.Errors()

		require.Equal(t, len(tc.wants), len(errors))

		testLetStatements_errors(t, tc.wants, errors)
	}
}

func testLetStatements_errors(t *testing.T, wants []wantError, errors []string) {
	for i, want := range wants {
		require.Equal(t, want.expectedError, errors[i])
	}
}

func TestReturnStatements(t *testing.T) {
	happyTestCases := []struct {
		input string
		wants []want
	}{
		{
			input: `
return 5;
return 10;
return 993322;
`,
			wants: []want{
				{},
				{},
				{},
			},
		},
	}

	for _, tc := range happyTestCases {
		tks := lexer.Lex(tc.input)
		p := parser.New(tks)

		program := p.ParseProgram()
		errors := p.Errors()

		require.Equal(t, len(tc.wants), len(program.Statements))
		require.Equal(t, 0, len(errors))

		testReturnStatement(t, tc.wants, program.Statements)
	}
}

func testReturnStatement(t *testing.T, wants []want, stmts []ast.Statement) {
	for i := range wants {
		s := stmts[i]
		require.Equal(t, "return", s.TokenLiteral())

		_, ok := s.(*statement.Return)
		require.True(t, ok)
	}
}

// TODO
// func TestReturnStatements_Errors(t *testing.T) {}
