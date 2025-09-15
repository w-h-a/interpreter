package statemachine

import (
	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/token"
)

type stateFn func(*statemachineLexer) stateFn

func lexToken(l *statemachineLexer) stateFn {
	l.skip()

	switch char := l.peek(); {
	case lexer.IsLetter(char):
		return lexIdentifier
	case lexer.IsDigit(char):
		return lexNumber
	default:
		return lexSymbol
	}
}

func lexIdentifier(l *statemachineLexer) stateFn {
	l.next()

	for l.pos < len(l.input) && lexer.IsLetter(l.input[l.pos]) {
		l.pos += 1
	}

	literal := l.input[l.start:l.pos]

	l.emit(token.LookupIdent(literal))

	return lexToken
}

func lexNumber(l *statemachineLexer) stateFn {
	l.next()

	for l.pos < len(l.input) && lexer.IsDigit(l.input[l.pos]) {
		l.pos += 1
	}

	l.emit(token.Int)

	return lexToken
}

func lexSymbol(l *statemachineLexer) stateFn {
	switch char := l.next(); char {
	case '=':
		return lexEqual
	case '+':
		l.emit(token.Plus)
	case '-':
		l.emit(token.Minus)
	case '!':
		return lexBang
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
	case 0:
		l.emit(token.EOF)
		return nil
	default:
		l.emit(token.Illegal)
	}

	return lexToken
}

func lexEqual(l *statemachineLexer) stateFn {
	if l.peek() == '=' {
		l.next()
		l.emit(token.Identical)
	} else {
		l.emit(token.Assign)
	}

	return lexToken
}

func lexBang(l *statemachineLexer) stateFn {
	if l.peek() == '=' {
		l.next()
		l.emit(token.NotIdentical)
	} else {
		l.emit(token.Bang)
	}

	return lexToken
}
