package lexer

import (
	"github.com/w-h-a/interpreter/internal/token"
)

type lexer struct {
	input  string
	start  int
	pos    int
	tokens chan token.Token
}

func (l *lexer) emit(t token.TokenType) {
	tk := token.Factory(t, l.input[l.start:l.pos])
	l.tokens <- tk
	l.start = l.pos
}

func (l *lexer) run() {
	for state := lex; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) next() byte {
	if l.pos >= len(l.input) {
		return 0
	}

	b := l.input[l.pos]

	l.pos += 1

	return b
}

func (l *lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}

	return l.input[l.pos]
}

func (l *lexer) skip() {
	for l.pos < len(l.input) && IsSpace(l.input[l.pos]) {
		l.pos += 1
	}

	l.start = l.pos
}

func Lex(input string) chan token.Token {
	tks := make(chan token.Token, 2)

	l := &lexer{
		input:  input,
		tokens: tks,
	}

	go l.run()

	return tks
}
