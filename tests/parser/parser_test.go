package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/parser"
	"github.com/w-h-a/interpreter/internal/parser/ast"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement"
)

type expectedPrefixOperatorExpression struct {
	operator string
	right    any
}

type expectedInfixOperatorExpression struct {
	operator string
	left     any
	right    any
}

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
				require.Equal(t, 4, len(errors))
				testParseErrors(t, "expected next token to be =, got INT", errors[0])
				testParseErrors(t, "expected next token to be IDENT, got =", errors[1])
				testParseErrors(t, "no parse prefix function for = found", errors[2])
				testParseErrors(t, "expected next token to be IDENT, got INT", errors[3])
			},
		},
		{
			name:  "identifier expression",
			input: `foobar;`,
			testFn: func(t *testing.T, program *ast.Program) {
				require.Equal(t, 1, len(program.Statements))
				testExpressionStatement(t, program.Statements[0], "foobar")
			},
			expectErr: false,
		},
		{
			name:  "integer expression",
			input: `5;`,
			testFn: func(t *testing.T, program *ast.Program) {
				require.Equal(t, 1, len(program.Statements))
				testExpressionStatement(t, program.Statements[0], 5)
			},
			expectErr: false,
		},
		{
			name: "boolean expression",
			input: `
true;
false;
`,
			testFn: func(t *testing.T, program *ast.Program) {
				require.Equal(t, 2, len(program.Statements))
				testExpressionStatement(t, program.Statements[0], true)
				testExpressionStatement(t, program.Statements[1], false)
			},
			expectErr: false,
		},
		{
			name: "prefix operator expressions",
			input: `
!5;
-15;
!true;
!false;
`,
			testFn: func(t *testing.T, program *ast.Program) {
				require.Equal(t, 4, len(program.Statements))
				testExpressionStatement(t, program.Statements[0], expectedPrefixOperatorExpression{operator: "!", right: 5})
				testExpressionStatement(t, program.Statements[1], expectedPrefixOperatorExpression{operator: "-", right: 15})
				testExpressionStatement(t, program.Statements[2], expectedPrefixOperatorExpression{operator: "!", right: true})
				testExpressionStatement(t, program.Statements[3], expectedPrefixOperatorExpression{operator: "!", right: false})
			},
			expectErr: false,
		},
		{
			name: "infix operator expressions",
			input: `
5 + 5;
5 - 5;
5 * 5;
5 / 5;
5 > 5;
5 < 5;
5 == 5;
5 != 5;
true == true;
true != false;
false == false;
`,
			testFn: func(t *testing.T, program *ast.Program) {
				require.Equal(t, 11, len(program.Statements))
				testExpressionStatement(t, program.Statements[0], expectedInfixOperatorExpression{operator: "+", left: 5, right: 5})
				testExpressionStatement(t, program.Statements[1], expectedInfixOperatorExpression{operator: "-", left: 5, right: 5})
				testExpressionStatement(t, program.Statements[2], expectedInfixOperatorExpression{operator: "*", left: 5, right: 5})
				testExpressionStatement(t, program.Statements[3], expectedInfixOperatorExpression{operator: "/", left: 5, right: 5})
				testExpressionStatement(t, program.Statements[4], expectedInfixOperatorExpression{operator: ">", left: 5, right: 5})
				testExpressionStatement(t, program.Statements[5], expectedInfixOperatorExpression{operator: "<", left: 5, right: 5})
				testExpressionStatement(t, program.Statements[6], expectedInfixOperatorExpression{operator: "==", left: 5, right: 5})
				testExpressionStatement(t, program.Statements[7], expectedInfixOperatorExpression{operator: "!=", left: 5, right: 5})
				testExpressionStatement(t, program.Statements[8], expectedInfixOperatorExpression{operator: "==", left: true, right: true})
				testExpressionStatement(t, program.Statements[9], expectedInfixOperatorExpression{operator: "!=", left: true, right: false})
				testExpressionStatement(t, program.Statements[10], expectedInfixOperatorExpression{operator: "==", left: false, right: false})
			},
			expectErr: false,
		},
		{
			name:  "program string 1",
			input: `-a * b`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "((-a) * b)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 2",
			input: `!-a`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "(!(-a))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 3",
			input: `a + b + c`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "((a + b) + c)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 4",
			input: `a + b - c`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "((a + b) - c)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 5",
			input: `a * b * c`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "((a * b) * c)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 6",
			input: `a * b / c`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "((a * b) / c)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 7",
			input: `a + b / c`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "(a + (b / c))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 8",
			input: `a + b * c + d / e - f`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "(((a + (b * c)) + (d / e)) - f)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 9",
			input: `3 + 4; -5 * 5`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "(3 + 4)((-5) * 5)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 10",
			input: `5 > 4 == 3 < 4`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "((5 > 4) == (3 < 4))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 11",
			input: `5 < 4 != 3 > 4`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "((5 < 4) != (3 > 4))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 12",
			input: `3 + 4 * 5 == 3 * 1 + 4 * 5`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 13",
			input: `3 > 5 == false`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "((3 > 5) == false)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 14",
			input: `3 < 5 == true`,
			testFn: func(t *testing.T, program *ast.Program) {
				got := program.String()
				require.Equal(t, "((3 < 5) == true)", got)
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

func testExpressionStatement(t *testing.T, s ast.Statement, expected any) {
	expStmt, ok := s.(*statement.Expression)
	require.True(t, ok)

	testExpression(t, expStmt.Expression, expected)
}

func testExpression(t *testing.T, e ast.Expression, expected any) {
	switch v := expected.(type) {
	case int:
		integer, ok := e.(*expression.Integer)
		require.True(t, ok)
		require.Equal(t, int64(v), integer.Value)
		require.Equal(t, fmt.Sprintf("%d", v), integer.TokenLiteral())
	case bool:
		boolean, ok := e.(*expression.Boolean)
		require.True(t, ok)
		require.Equal(t, v, boolean.Value)
		require.Equal(t, fmt.Sprintf("%t", v), boolean.TokenLiteral())
	case string:
		identifier, ok := e.(*expression.Identifier)
		require.True(t, ok)
		require.Equal(t, v, identifier.Value)
		require.Equal(t, v, identifier.TokenLiteral())
	case expectedPrefixOperatorExpression:
		prefixOperator, ok := e.(*expression.PrefixOperator)
		require.True(t, ok)
		require.Equal(t, v.operator, prefixOperator.Operator)
		testExpression(t, prefixOperator.Right, v.right)
	case expectedInfixOperatorExpression:
		infixOperator, ok := e.(*expression.InfixOperator)
		require.True(t, ok)
		require.Equal(t, v.operator, infixOperator.Operator)
		testExpression(t, infixOperator.Left, v.left)
		testExpression(t, infixOperator.Right, v.right)
	default:
		t.Errorf("Expression assertion type not handled. got=%T", expected)
	}
}

func testParseErrors(t *testing.T, want string, got string) {
	require.Equal(t, want, got)
}
