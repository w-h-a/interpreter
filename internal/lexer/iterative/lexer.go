package iterative

import (
	"unicode/utf8"

	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/token"
)

type iterativeLexer struct {
	input string
	start int
	pos   int
}

func (l *iterativeLexer) NextToken() token.Token {
	var tk token.Token

	switch char := l.next(); char {
	case '=':
		tk, _ = l.result(token.Assign)
	case '+':
		tk, _ = l.result(token.Plus)
	case '(':
		tk, _ = l.result(token.ParenLeft)
	case ')':
		tk, _ = l.result(token.ParenRight)
	case '{':
		tk, _ = l.result(token.BraceLeft)
	case '}':
		tk, _ = l.result(token.BraceRight)
	case ',':
		tk, _ = l.result(token.Comma)
	case ';':
		tk, _ = l.result(token.Semicolon)
	case 0:
		tk.Literal = ""
		tk.Type = token.EOF
	default:
		tk, _ = l.result(token.Illegal)
	}

	return tk
}

// next consumes the next rune
func (l *iterativeLexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}

	r, size := utf8.DecodeRuneInString(l.input[l.pos:])

	l.pos += size

	return r
}

// result returns a token and resets the start
func (l *iterativeLexer) result(t token.TokenType) (token.Token, error) {
	tk, err := token.Factory(t, l.input[l.start:l.pos])
	l.start = l.pos
	return tk, err
}

func New(input string) lexer.Lexer {
	l := &iterativeLexer{
		input: input,
	}

	return l
}
