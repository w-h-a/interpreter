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

	switch char := l.peek(); {
	case lexer.IsLetter(char):
		tk, _ = l.lexIdentifier()
	case lexer.IsDigit(char):
		tk, _ = l.lexNumber()
	default:
		tk, _ = l.lexSymbol()
	}

	return tk
}

func (l *iterativeLexer) lexIdentifier() (token.Token, error) {
	l.next()

	for l.pos < len(l.input) && lexer.IsLetter(l.input[l.pos]) {
		l.pos += 1
	}

	literal := l.input[l.start:l.pos]

	return l.generate(token.LookupIdent(literal))
}

func (l *iterativeLexer) lexNumber() (token.Token, error) {
	l.next()

	for l.pos < len(l.input) && lexer.IsDigit(l.input[l.pos]) {
		l.pos += 1
	}

	return l.generate(token.Int)
}

func (l *iterativeLexer) lexSymbol() (token.Token, error) {
	var tk token.Token

	switch char := l.next(); char {
	case '=':
		tk, _ = l.lexEqual()
	case '+':
		tk, _ = l.generate(token.Plus)
	case '-':
		tk, _ = l.generate(token.Minus)
	case '!':
		tk, _ = l.lexBang()
	case '*':
		tk, _ = l.generate(token.Asterisk)
	case '/':
		tk, _ = l.generate(token.Slash)
	case '<':
		tk, _ = l.generate(token.LessThan)
	case '>':
		tk, _ = l.generate(token.GreaterThan)
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
		tk, _ = l.generate(token.Illegal)
	}

	return tk, nil
}

func (l *iterativeLexer) lexEqual() (token.Token, error) {
	var tk token.Token

	if l.peek() == '=' {
		l.next()
		tk, _ = l.generate(token.Identical)
	} else {
		tk, _ = l.generate(token.Assign)
	}

	return tk, nil
}

func (l *iterativeLexer) lexBang() (token.Token, error) {
	var tk token.Token

	if l.peek() == '=' {
		l.next()
		tk, _ = l.generate(token.NotIdentical)
	} else {
		tk, _ = l.generate(token.Bang)
	}

	return tk, nil
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

// peek checks the next byte but does not consume it
func (l *iterativeLexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}

	return l.input[l.pos]
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
