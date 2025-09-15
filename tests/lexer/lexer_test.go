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

type want struct {
	expectedType    token.TokenType
	expectedLiteral string
}

var (
	testCases = []map[string]any{
		{
			"input": `=+(){},;!-/*<>`,
			"want": []want{
				{token.Assign, "="},
				{token.Plus, "+"},
				{token.ParenLeft, "("},
				{token.ParenRight, ")"},
				{token.BraceLeft, "{"},
				{token.BraceRight, "}"},
				{token.Comma, ","},
				{token.Semicolon, ";"},
				{token.Bang, "!"},
				{token.Minus, "-"},
				{token.Slash, "/"},
				{token.Asterisk, "*"},
				{token.LessThan, "<"},
				{token.GreaterThan, ">"},
				{token.EOF, ""},
			},
		},
		{
			"input": `10 == 10;
10 != 9;`,
			"want": []want{
				{token.Int, "10"},
				{token.Identical, "=="},
				{token.Int, "10"},
				{token.Semicolon, ";"},
				{token.Int, "10"},
				{token.NotIdentical, "!="},
				{token.Int, "9"},
				{token.Semicolon, ";"},
			},
		},
		{
			"input": `let five = 5;
let ten = 10;
let add = fn(x, y) {
x + y;
};
let result = add(five, ten);
if (5 < 10) {
  return true;
} else {
  return false;
}
`,
			"want": []want{
				{token.Let, "let"},
				{token.Ident, "five"},
				{token.Assign, "="},
				{token.Int, "5"},
				{token.Semicolon, ";"},
				{token.Let, "let"},
				{token.Ident, "ten"},
				{token.Assign, "="},
				{token.Int, "10"},
				{token.Semicolon, ";"},
				{token.Let, "let"},
				{token.Ident, "add"},
				{token.Assign, "="},
				{token.Function, "fn"},
				{token.ParenLeft, "("},
				{token.Ident, "x"},
				{token.Comma, ","},
				{token.Ident, "y"},
				{token.ParenRight, ")"},
				{token.BraceLeft, "{"},
				{token.Ident, "x"},
				{token.Plus, "+"},
				{token.Ident, "y"},
				{token.Semicolon, ";"},
				{token.BraceRight, "}"},
				{token.Semicolon, ";"},
				{token.Let, "let"},
				{token.Ident, "result"},
				{token.Assign, "="},
				{token.Ident, "add"},
				{token.ParenLeft, "("},
				{token.Ident, "five"},
				{token.Comma, ","},
				{token.Ident, "ten"},
				{token.ParenRight, ")"},
				{token.Semicolon, ";"},
				{token.If, "if"},
				{token.ParenLeft, "("},
				{token.Int, "5"},
				{token.LessThan, "<"},
				{token.Int, "10"},
				{token.ParenRight, ")"},
				{token.BraceLeft, "{"},
				{token.Return, "return"},
				{token.True, "true"},
				{token.Semicolon, ";"},
				{token.BraceRight, "}"},
				{token.Else, "else"},
				{token.BraceLeft, "{"},
				{token.Return, "return"},
				{token.False, "false"},
				{token.Semicolon, ";"},
				{token.BraceRight, "}"},
				{token.EOF, ""},
			},
		},
	}
)

func TestIterativeLexer(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip("SKIPPING UNIT TEST")
		return
	}

	for _, tc := range testCases {
		input := tc["input"]
		l := iterative.New(input.(string))
		runLexerTest(t, tc["want"].([]want), l)
	}
}

func TestStatemachineLexer(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip("SKIPPING UNIT TEST")
		return
	}

	for _, tc := range testCases {
		input := tc["input"]
		l := statemachine.New(input.(string))
		runLexerTest(t, tc["want"].([]want), l)
	}
}

func runLexerTest(t *testing.T, wants []want, l lexer.Lexer) {
	for _, want := range wants {
		tk := l.NextToken()
		require.Equal(t, want.expectedType, tk.Type)
		require.Equal(t, want.expectedLiteral, tk.Literal)
	}
}
