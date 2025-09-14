package statemachine

import (
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/token"
)

type stateFn func(*statemachineLexer) stateFn

func lexToken(l *statemachineLexer) stateFn {
	l.skip()

	switch char := l.next(); char {
	case '=':
		l.emit(token.Assign)
	case '+':
		l.emit(token.Plus)
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
	case 0:
		l.emit(token.EOF)
		return nil
	default:
		if lexer.IsLetter(char) {
			return lexIdentifier
		} else if lexer.IsDigit(char) {
			return lexNumber
		} else {
			l.emit(token.Illegal)
		}
	}

	return lexToken
}

func lexIdentifier(l *statemachineLexer) stateFn {
	for l.pos < len(l.input) && lexer.IsLetter(l.input[l.pos]) {
		l.pos += 1
	}

	literal := l.input[l.start:l.pos]

	l.emit(token.LookupIdent(literal))

	return lexToken
}

func lexNumber(l *statemachineLexer) stateFn {
	for l.pos < len(l.input) && lexer.IsDigit(l.input[l.pos]) {
		l.pos += 1
	}

	l.emit(token.Int)

	return lexToken
}
