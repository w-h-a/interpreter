package lexer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/lexer/iterative"
	statemachine "github.com/w-h-a/interpreter/internal/lexer/state_machine"
	"github.com/w-h-a/interpreter/internal/token"
)

var (
	testCases = []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Assign, "="},
		{token.Plus, "+"},
		{token.ParenLeft, "("},
		{token.ParenRight, ")"},
		{token.BraceLeft, "{"},
		{token.BraceRight, "}"},
		{token.Comma, ","},
		{token.Semicolon, ";"},
		{token.EOF, ""},
	}
)

func TestIterativeLexer(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip("SKIPPING UNIT TEST")
		return
	}

	input := `=+(){},;`
	l := iterative.New(input)
	runLexerTest(t, l)
}

func TestStatemachineLexer(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip("SKIPPING UNIT TEST")
		return
	}

	input := `=+(){},;`
	l := statemachine.New(input)
	runLexerTest(t, l)
}

func runLexerTest(t *testing.T, l lexer.Lexer) {
	for _, tc := range testCases {
		tk := l.NextToken()
		require.Equal(t, tc.expectedType, tk.Type)
		require.Equal(t, tc.expectedLiteral, tk.Literal)
	}
}
