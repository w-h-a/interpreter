package evaluator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/interpreter/internal/evaluator"
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/object"
	"github.com/w-h-a/interpreter/internal/object/integer"
	"github.com/w-h-a/interpreter/internal/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output int64
	}{
		{"should evaluate '5' as 5", "5", 5},
		{"should evaluate '10' as 10", "10", 10},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluated := testEval(t, test.input)
			testIntegerObject(t, test.output, evaluated)
		})
	}
}

func testEval(t *testing.T, input string) object.Object {
	tks := lexer.Lex(input)
	p := parser.New(tks)
	program := p.ParseProgram()
	errors := p.Errors()

	require.True(t, len(errors) == 0)

	return evaluator.Eval(program)
}

func testIntegerObject(t *testing.T, expected int64, obj object.Object) {
	result, ok := obj.(*integer.Integer)
	require.True(t, ok)
	require.Equal(t, expected, result.Value)
}
