package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/parser"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression/boolean"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression/call"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression/function"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression/identifier"
	ifexpression "github.com/w-h-a/interpreter/internal/parser/ast/expression/if"
	infixoperator "github.com/w-h-a/interpreter/internal/parser/ast/expression/infix_operator"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression/integer"
	prefixoperator "github.com/w-h-a/interpreter/internal/parser/ast/expression/prefix_operator"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement"
	expressionstatement "github.com/w-h-a/interpreter/internal/parser/ast/statement/expression"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement/let"
	returnstatement "github.com/w-h-a/interpreter/internal/parser/ast/statement/return"
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

type expectedIfExpression struct {
	condition      expectedInfixOperatorExpression
	consequenceLen int
	alternativeLen int
}

type expectedCallExpression struct {
	function string
	args     []any
}

type expectedFunctionExpression struct {
	params  []string
	bodyLen int
}

func TestParseProgram(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		testFn     func(t *testing.T, program *statement.Program)
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
			testFn: func(t *testing.T, program *statement.Program) {
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
			testFn: func(t *testing.T, program *statement.Program) {
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
				require.Equal(t, 4, len(errors))
				testParseErrors(t, "expected next token to be =, got INT", errors[0])
				testParseErrors(t, "expected next token to be IDENT, got =", errors[1])
				testParseErrors(t, "no parse function for = found", errors[2])
				testParseErrors(t, "expected next token to be IDENT, got INT", errors[3])
			},
		},
		{
			name:  "identifier expression",
			input: `foobar;`,
			testFn: func(t *testing.T, program *statement.Program) {
				require.Equal(t, 1, len(program.Statements))
				testExpressionStatement(t, program.Statements[0], "foobar")
			},
			expectErr: false,
		},
		{
			name:  "integer expression",
			input: `5;`,
			testFn: func(t *testing.T, program *statement.Program) {
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
			testFn: func(t *testing.T, program *statement.Program) {
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
			testFn: func(t *testing.T, program *statement.Program) {
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
			testFn: func(t *testing.T, program *statement.Program) {
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
			name: "if expressions",
			input: `
if (x < y) { x };
if (x > y) { x } else { y };
`,
			testFn: func(t *testing.T, program *statement.Program) {
				require.Equal(t, 2, len(program.Statements))
				testExpressionStatement(t, program.Statements[0], expectedIfExpression{
					condition:      expectedInfixOperatorExpression{operator: "<", left: "x", right: "y"},
					consequenceLen: 1,
					alternativeLen: 0,
				})
				testExpressionStatement(t, program.Statements[1], expectedIfExpression{
					condition:      expectedInfixOperatorExpression{operator: ">", left: "x", right: "y"},
					consequenceLen: 1,
					alternativeLen: 1,
				})
			},
			expectErr: false,
		},
		{
			name: "call expressions",
			input: `
add(1, 2 * 3);
`,
			testFn: func(t *testing.T, program *statement.Program) {
				require.Equal(t, 1, len(program.Statements))
				testExpressionStatement(t, program.Statements[0], expectedCallExpression{
					function: "add",
					args: []any{
						1,
						expectedInfixOperatorExpression{
							operator: "*",
							left:     2,
							right:    3,
						},
					},
				})
			},
			expectErr: false,
		},
		{
			name: "function expressions",
			input: `
fn(x, y) { x + y; };
fn() {};
fn(x, y, z) {};
`,
			testFn: func(t *testing.T, program *statement.Program) {
				require.Equal(t, 3, len(program.Statements))
				testExpressionStatement(t, program.Statements[0], expectedFunctionExpression{
					params:  []string{"x", "y"},
					bodyLen: 1,
				})
				testExpressionStatement(t, program.Statements[1], expectedFunctionExpression{
					params:  []string{},
					bodyLen: 0,
				})
				testExpressionStatement(t, program.Statements[2], expectedFunctionExpression{
					params:  []string{"x", "y", "z"},
					bodyLen: 0,
				})
			},
			expectErr: false,
		},
		{
			name: "malformed function expressions",
			input: `
fn(5) {};
`,
			expectErr: true,
			testErrsFn: func(t *testing.T, errors []string) {
				require.Equal(t, 4, len(errors))
				testParseErrors(t, `failed to parse expression literal "5": expected identifier as function parameter, got INT`, errors[0])
				testParseErrors(t, "no parse function for ) found", errors[1])
				testParseErrors(t, "no parse function for { found", errors[2])
				testParseErrors(t, "no parse function for } found", errors[3])
			},
		},
		{
			name:  "program string 1",
			input: `-a * b`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((-a) * b)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 2",
			input: `!-a`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "(!(-a))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 3",
			input: `a + b + c`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((a + b) + c)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 4",
			input: `a + b - c`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((a + b) - c)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 5",
			input: `a * b * c`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((a * b) * c)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 6",
			input: `a * b / c`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((a * b) / c)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 7",
			input: `a + b / c`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "(a + (b / c))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 8",
			input: `a + b * c + d / e - f`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "(((a + (b * c)) + (d / e)) - f)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 9",
			input: `3 + 4; -5 * 5`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "(3 + 4)((-5) * 5)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 10",
			input: `5 > 4 == 3 < 4`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((5 > 4) == (3 < 4))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 11",
			input: `5 < 4 != 3 > 4`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((5 < 4) != (3 > 4))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 12",
			input: `3 + 4 * 5 == 3 * 1 + 4 * 5`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 13",
			input: `3 > 5 == false`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((3 > 5) == false)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 14",
			input: `3 < 5 == true`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((3 < 5) == true)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 15",
			input: `1 + (2 + 3) + 4`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((1 + (2 + 3)) + 4)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 16",
			input: `(5 + 5) * 2`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((5 + 5) * 2)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 17",
			input: `2 / (5 + 5)`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "(2 / (5 + 5))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 18",
			input: `-(5 + 5)`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "(-(5 + 5))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 19",
			input: `!(true == true)`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "(!(true == true))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 20",
			input: `a + add(b * c) + d`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "((a + add((b * c))) + d)", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 21",
			input: `add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))", got)
			},
			expectErr: false,
		},
		{
			name:  "program string 22",
			input: `add(a + b + c * d / f + g)`,
			testFn: func(t *testing.T, program *statement.Program) {
				got := program.String()
				require.Equal(t, "add((((a + b) + ((c * d) / f)) + g))", got)
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

			t.Logf("ERRORS: %v", errors)

			if tc.expectErr {
				tc.testErrsFn(t, errors)
				return
			}

			require.Equal(t, 0, len(errors))
			tc.testFn(t, program)
		})
	}
}

func testLetStatement(t *testing.T, s statement.Statement, name string) {
	require.Equal(t, "let", s.TokenLiteral())
	letStmt, ok := s.(*let.Let)
	require.True(t, ok)
	require.Equal(t, name, letStmt.Name.TokenLiteral())
	require.Equal(t, name, letStmt.Name.Value)
}

func testReturnStatement(t *testing.T, s statement.Statement) {
	require.Equal(t, "return", s.TokenLiteral())
	_, ok := s.(*returnstatement.Return)
	require.True(t, ok)
}

func testExpressionStatement(t *testing.T, s statement.Statement, expected any) {
	expStmt, ok := s.(*expressionstatement.Expression)
	require.True(t, ok)

	testExpression(t, expStmt.Expression, expected)
}

func testExpression(t *testing.T, e expression.Expression, expected any) {
	switch v := expected.(type) {
	case int:
		integer, ok := e.(*integer.Integer)
		require.True(t, ok)
		require.Equal(t, int64(v), integer.Value)
		require.Equal(t, fmt.Sprintf("%d", v), integer.TokenLiteral())
	case bool:
		boolean, ok := e.(*boolean.Boolean)
		require.True(t, ok)
		require.Equal(t, v, boolean.Value)
		require.Equal(t, fmt.Sprintf("%t", v), boolean.TokenLiteral())
	case string:
		identifier, ok := e.(*identifier.Identifier)
		require.True(t, ok)
		require.Equal(t, v, identifier.Value)
		require.Equal(t, v, identifier.TokenLiteral())
	case expectedPrefixOperatorExpression:
		prefixOperator, ok := e.(*prefixoperator.PrefixOperator)
		require.True(t, ok)
		require.Equal(t, v.operator, prefixOperator.Operator)
		testExpression(t, prefixOperator.Right, v.right)
	case expectedInfixOperatorExpression:
		infixOperator, ok := e.(*infixoperator.InfixOperator)
		require.True(t, ok)
		require.Equal(t, v.operator, infixOperator.Operator)
		testExpression(t, infixOperator.Left, v.left)
		testExpression(t, infixOperator.Right, v.right)
	case expectedIfExpression:
		ifExpression, ok := e.(*ifexpression.If)
		require.True(t, ok)
		testExpression(t, ifExpression.Condition, v.condition)
		require.NotNil(t, ifExpression.Consequence)
		require.Equal(t, v.consequenceLen, len(ifExpression.Consequence.Statements))
		if v.alternativeLen > 0 {
			require.NotNil(t, ifExpression.Alternative)
			require.Equal(t, v.alternativeLen, len(ifExpression.Alternative.Statements))
		} else {
			require.Nil(t, ifExpression.Alternative)
		}
	case expectedCallExpression:
		callExpression, ok := e.(*call.Call)
		require.True(t, ok)
		testExpression(t, callExpression.Function, v.function)
		for i, arg := range v.args {
			testExpression(t, callExpression.Arguments[i], arg)
		}
	case expectedFunctionExpression:
		functionExpression, ok := e.(*function.Function)
		require.True(t, ok)
		require.Equal(t, len(v.params), len(functionExpression.Parameters))
		for i, ident := range v.params {
			testExpression(t, functionExpression.Parameters[i], ident)
		}
		require.Equal(t, v.bodyLen, len(functionExpression.Body.Statements))
	default:
		t.Errorf("Expression assertion type not handled. got=%T", expected)
	}
}

func testParseErrors(t *testing.T, want string, got string) {
	require.Equal(t, want, got)
}
