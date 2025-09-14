package statemachine

import (
	"unicode/utf8"

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

func (l *statemachineLexer) run() {
	for state := lexToken; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// next consumes the next rune
func (l *statemachineLexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}

	r, size := utf8.DecodeRuneInString(l.input[l.pos:])

	l.pos += size

	return r
}

// emit sends a token to the channel and resets the start
func (l *statemachineLexer) emit(t token.TokenType) {
	tk, _ := token.Factory(t, l.input[l.start:l.pos])
	l.tokens <- tk
	l.start = l.pos
}

func New(input string) lexer.Lexer {
	l := &statemachineLexer{
		input:  input,
		tokens: make(chan token.Token),
	}

	go l.run()

	return l
}
