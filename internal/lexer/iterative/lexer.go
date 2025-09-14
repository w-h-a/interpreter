package iterative

import (
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

	l.skip()

	switch char := l.next(); char {
	case '=':
		tk, _ = l.generate(token.Assign)
	case '+':
		tk, _ = l.generate(token.Plus)
	case '(':
		tk, _ = l.generate(token.ParenLeft)
	case ')':
		tk, _ = l.generate(token.ParenRight)
	case '{':
		tk, _ = l.generate(token.BraceLeft)
	case '}':
		tk, _ = l.generate(token.BraceRight)
	case ',':
		tk, _ = l.generate(token.Comma)
	case ';':
		tk, _ = l.generate(token.Semicolon)
	case 0:
		tk.Literal = ""
		tk.Type = token.EOF
	default:
		if lexer.IsLetter(char) {
			tk, _ = l.lexIdentifier()
		} else if lexer.IsDigit(char) {
			tk, _ = l.lexNumber()
		} else {
			tk, _ = l.generate(token.Illegal)
		}
	}

	return tk
}

func (l *iterativeLexer) lexIdentifier() (token.Token, error) {
	for l.pos < len(l.input) && lexer.IsLetter(l.input[l.pos]) {
		l.pos += 1
	}

	literal := l.input[l.start:l.pos]

	return l.generate(token.LookupIdent(literal))
}

func (l *iterativeLexer) lexNumber() (token.Token, error) {
	for l.pos < len(l.input) && lexer.IsDigit(l.input[l.pos]) {
		l.pos += 1
	}

	return l.generate(token.Int)
}

// generate returns a token and resets the start
func (l *iterativeLexer) generate(t token.TokenType) (token.Token, error) {
	tk, err := token.Factory(t, l.input[l.start:l.pos])
	l.start = l.pos
	return tk, err
}

// next consumes the next byte
func (l *iterativeLexer) next() byte {
	if l.pos >= len(l.input) {
		return 0
	}

	b := l.input[l.pos]

	l.pos += 1

	return b
}

// skip skips whitespaces and resets the start
func (l *iterativeLexer) skip() {
	for l.pos < len(l.input) && lexer.IsSpace(l.input[l.pos]) {
		l.pos += 1
	}

	l.start = l.pos
}

func New(input string) lexer.Lexer {
	l := &iterativeLexer{
		input: input,
	}

	return l
}
