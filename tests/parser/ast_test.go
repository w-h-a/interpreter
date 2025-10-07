package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression/identifier"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement/let"
	"github.com/w-h-a/interpreter/internal/token"
)

func TestAST(t *testing.T) {
	testCases := []struct {
		name     string
		program  *statement.Program
		expected string
	}{
		{
			name: "let statement with identifier",
			program: &statement.Program{
				Statements: []statement.Statement{
					&let.Let{
						Token: token.Token{Type: token.Let, Literal: "let"},
						Name: &identifier.Identifier{
							Token: token.Token{Type: token.Ident, Literal: "myVar"},
							Value: "myVar",
						},
						Value: &identifier.Identifier{
							Token: token.Token{Type: token.Ident, Literal: "anotherVar"},
							Value: "anotherVar",
						},
					},
				},
			},
			expected: "let myVar = anotherVar;",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.program.String())
		})
	}
}
