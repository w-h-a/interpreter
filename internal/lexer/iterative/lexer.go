package iterative

import (
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/token"
)

type iterativeLexer struct {
	input  string
	start  int
	pos    int
	tokens chan token.Token
}

func (l *iterativeLexer) NextToken() token.Token {
	tk, ok := <-l.tokens
	if !ok {
		return token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	}
	return tk
}

func (l *iterativeLexer) lex() bool {
	l.skip()

	switch char := l.peek(); {
	case char == 0:
		l.lexStop()
		return false
	case lexer.IsLetter(char):
		l.lexIdentifier()
	case lexer.IsDigit(char):
		l.lexNumber()
	default:
		l.lexSymbol()
	}

	return true
}

func (l *iterativeLexer) lexIdentifier() {
	l.next()

	for l.pos < len(l.input) && lexer.IsLetter(l.input[l.pos]) {
		l.pos += 1
	}

	literal := l.input[l.start:l.pos]

	l.emit(token.LookupIdent(literal))
}

func (l *iterativeLexer) lexNumber() {
	l.next()

	for l.pos < len(l.input) && lexer.IsDigit(l.input[l.pos]) {
		l.pos += 1
	}

	l.emit(token.Int)
}

func (l *iterativeLexer) lexSymbol() {
	switch char := l.next(); char {
	case '=':
		l.lexEqual()
	case '+':
		l.emit(token.Plus)
	case '-':
		l.emit(token.Minus)
	case '!':
		l.lexBang()
	case '*':
		l.emit(token.Asterisk)
	case '/':
		l.emit(token.Slash)
	case '<':
		l.emit(token.LessThan)
	case '>':
		l.emit(token.GreaterThan)
	case '(':
		l.emit(token.ParenLeft)
	case ')':
		l.emit(token.ParenRight)
	case '{':
		l.emit(token.BraceLeft)
	case '}':
		l.emit(token.BraceRight)
	case ',':
		l.emit(token.Comma)
	case ';':
		l.emit(token.Semicolon)
	default:
		l.emit(token.Illegal)
	}
}

func (l *iterativeLexer) lexEqual() {
	if l.peek() == '=' {
		l.next()
		l.emit(token.Identical)
	} else {
		l.emit(token.Assign)
	}
}

func (l *iterativeLexer) lexBang() {
	if l.peek() == '=' {
		l.next()
		l.emit(token.NotIdentical)
	} else {
		l.emit(token.Bang)
	}
}

func (l *iterativeLexer) lexStop() {
	l.emit(token.EOF)
}

func (l *iterativeLexer) emit(t token.TokenType) {
	tk, _ := token.Factory(t, l.input[l.start:l.pos])
	l.tokens <- tk
	l.start = l.pos
}

func (l *iterativeLexer) run() {
	for l.lex() {
	}
	close(l.tokens)
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
		input:  input,
		tokens: make(chan token.Token, 2),
	}

	go l.run()

	return l
}
