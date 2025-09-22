package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/token"
)

type want struct {
	expectedType    token.TokenType
	expectedLiteral string
}

var (
	testCases = []struct {
		input string
		wants []want
	}{
		{
			input: `=+(){},;!-/*<>`,
			wants: []want{
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
			input: `10 == 10;
10 != 9;`,
			wants: []want{
				{token.Int, "10"},
				{token.Identical, "=="},
				{token.Int, "10"},
				{token.Semicolon, ";"},
				{token.Int, "10"},
				{token.NotIdentical, "!="},
				{token.Int, "9"},
				{token.Semicolon, ";"},
				{token.EOF, ""},
			},
		},
		{
			input: `let five = 5;
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
			wants: []want{
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

func TestLexer(t *testing.T) {
	for _, tc := range testCases {
		tks := lexer.Lex(tc.input)
		runLexerTest(t, tc.wants, tks)
	}
}

func runLexerTest(t *testing.T, wants []want, tks chan token.Token) {
	for _, want := range wants {
		tk := <-tks
		require.Equal(t, want.expectedType, tk.Type)
		require.Equal(t, want.expectedLiteral, tk.Literal)
	}
	_, ok := <-tks
	require.False(t, ok, "channel was not closed")
}
