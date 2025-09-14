package statemachine

import (
	"github.com/w-h-a/interpreter/internal/token"
)

type stateFn func(*statemachineLexer) stateFn

func lexToken(l *statemachineLexer) stateFn {
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
		l.emit(token.Illegal)
	}

	return lexToken
}
