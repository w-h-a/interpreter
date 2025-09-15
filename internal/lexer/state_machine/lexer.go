package statemachine

import (
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/token"
)

type statemachineLexer struct {
	input  string
	start  int
	pos    int
	tokens chan token.Token
}

func (l *statemachineLexer) NextToken() token.Token {
	tk, ok := <-l.tokens
	if !ok {
		return token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	}
	return tk
}

func (l *statemachineLexer) emit(t token.TokenType) {
	tk, _ := token.Factory(t, l.input[l.start:l.pos])
	l.tokens <- tk
	l.start = l.pos
}

func (l *statemachineLexer) run() {
	for state := lex; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// next consumes the next byte
func (l *statemachineLexer) next() byte {
	if l.pos >= len(l.input) {
		return 0
	}

	b := l.input[l.pos]

	l.pos += 1

	return b
}

// peek checks the next byte but does not consume it
func (l *statemachineLexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}

	return l.input[l.pos]
}

// skip skips whitespace and resets the start
func (l *statemachineLexer) skip() {
	for l.pos < len(l.input) && lexer.IsSpace(l.input[l.pos]) {
		l.pos += 1
	}

	l.start = l.pos
}

func New(input string) lexer.Lexer {
	l := &statemachineLexer{
		input:  input,
		tokens: make(chan token.Token, 2),
	}

	go l.run()

	return l
}
