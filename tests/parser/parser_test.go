package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/parser"
	"github.com/w-h-a/interpreter/internal/parser/ast"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement"
)

func TestParseProgram(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		testFn     func(t *testing.T, program *ast.Program)
		expectErr  bool
		testErrsFn func(t *testing.T, errors []string)
	}{
		{
			name: "let statements happy path",
			input: `
let x = 5;
let y = 10;
let foobar = 838383;
`,
			testFn: func(t *testing.T, program *ast.Program) {
				require.Equal(t, 3, len(program.Statements))
				testLetStatement(t, program.Statements[0], "x")
				testLetStatement(t, program.Statements[1], "y")
				testLetStatement(t, program.Statements[2], "foobar")
			},
			expectErr: false,
		},
		{
			name: "return statements happy path",
			input: `
return 5;
return 10;
return 993322;
`,
			testFn: func(t *testing.T, program *ast.Program) {
				require.Equal(t, 3, len(program.Statements))
				testReturnStatement(t, program.Statements[0])
				testReturnStatement(t, program.Statements[1])
				testReturnStatement(t, program.Statements[2])
			},
			expectErr: false,
		},
		{
			name: "malformed let statements error path",
			input: `
let x 5;
let = 10;
let 838383;
`,
			expectErr: true,
			testErrsFn: func(t *testing.T, errors []string) {
				t.Logf("errors %+v", errors)
				require.Equal(t, 7, len(errors))
				require.Equal(t, "expected next token to be =, got INT", errors[0])
				require.Equal(t, "no parse prefix function for INT found", errors[1])
				require.Equal(t, "expected next token to be IDENT, got =", errors[2])
				require.Equal(t, "no parse prefix function for = found", errors[3])
				require.Equal(t, "no parse prefix function for INT found", errors[4])
				require.Equal(t, "expected next token to be IDENT, got INT", errors[5])
				require.Equal(t, "no parse prefix function for INT found", errors[6])
			},
		},
		{
			name:  "identifier expression",
			input: `foobar;`,
			testFn: func(t *testing.T, program *ast.Program) {
				require.Equal(t, 1, len(program.Statements))
				stmt, ok := program.Statements[0].(*statement.Expression)
				require.True(t, ok)
				ident, ok := stmt.Expression.(*expression.Identifier)
				require.True(t, ok)
				require.Equal(t, "foobar", ident.Value)
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tks := lexer.Lex(tc.input)
			p := parser.New(tks)

			program := p.ParseProgram()
			errors := p.Errors()

			if tc.expectErr {
				tc.testErrsFn(t, errors)
				return
			}

			require.Equal(t, 0, len(errors))
			tc.testFn(t, program)
		})
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	require.Equal(t, "let", s.TokenLiteral())
	letStmt, ok := s.(*statement.Let)
	require.True(t, ok)
	require.Equal(t, name, letStmt.Name.TokenLiteral())
	require.Equal(t, name, letStmt.Name.Value)
}

func testReturnStatement(t *testing.T, s ast.Statement) {
	require.Equal(t, "return", s.TokenLiteral())
	_, ok := s.(*statement.Return)
	require.True(t, ok)
}
